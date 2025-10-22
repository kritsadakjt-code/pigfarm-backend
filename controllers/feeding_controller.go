package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateFeeding(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	var input dto.FeedingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	parsedDate, err := time.ParseInLocation("2006-01-02 15:04", input.DateTime, loc)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD HH:MM"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tx := config.DB.Begin()
	// error ที่ไม่ได้มาจากเราเช็คเองเช่น array เกินขนาด x:=10/0
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// ดึงข้อมูลหมูที่เลือก
	var validPigsToFeed []models.Pig
	if len(input.PigIDs) > 0 {
		var pigs []models.Pig
		if err := tx.Where("id IN ?", input.PigIDs).Find(&pigs).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Pigs not found"})
		}

		// แยกหมูที่สามารถให้อาหารได้
		for _, pig := range pigs {
			if pig.Status != "ตายเเล้ว" && pig.Status != "ขายเเล้ว" {
				validPigsToFeed = append(validPigsToFeed, pig)
			}
		}
	}

	// check ชนิดอาหาร
	food := &models.FoodStock{}
	if err := tx.First(food, "id = ?", input.FoodID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "FoodStock not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error (food stock)"})
	}

	if input.Amount <= 0 {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "amount must be greater than 0"})
	}
	if food.Amount < input.Amount {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Not enough food in stock"})
	}
	food.Amount -= input.Amount

	if err := tx.Save(&food).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update food stock"})
	}

	feeding := models.Feeding{
		FoodID:    food.ID,
		DateTime:  parsedDate,
		Amount:    input.Amount,
		Note:      input.Note,
		CreatedBy: user_id,
		UpdatedBy: user_id,
	}

	if err := tx.Create(&feeding).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "failed to create feeding"})
	}
	var items []models.FeedingItem
	var validPigIDs []uint
	for _, pig := range validPigsToFeed {
		items = append(items, models.FeedingItem{
			FeedingID: feeding.ID,
			PigID:     pig.ID,
		})
		validPigIDs = append(validPigIDs, pig.ID)
	}
	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create feeding Items"})
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Transaction commit failed"})
	}

	if err := config.DB.Preload("FoodStock").Preload("Creator").Preload("Updater").Preload("Items.Pig").First(&feeding, "id = ?", feeding.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to preload feeding"})
	}

	resp := dto.FeedingResponse{
		ID:          feeding.ID,
		FoodID:      food.ID,
		PigIDs:      validPigIDs,
		DateTime:    parsedDate,
		Amount:      feeding.Amount,
		Note:        feeding.Note,
		FoodName:    feeding.FoodStock.Name,
		CreatedName: feeding.Creator.FullName,
		UpdatedName: feeding.Updater.FullName,
	}

	return c.Status(201).JSON(resp)

}

func SearchFeeding(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("keyword"))
	var feedings []models.Feeding
	db := config.DB.Model(&models.Feeding{}).Preload("FoodStock").Preload("Creator").Preload("Updater").
		Joins("LEFT JOIN food_stocks food on food.id = feedings.food_id").
		Joins("LEFT JOIN users creator on creator.id = feedings.created_by").
		Joins("LEFT JOIN users updater on updater.id = feedings.updated_by")
	kw := "%" + keyword + "%"
	if keyword != "" {
		if err := db.Where("food.name ILIKE ? OR creator.full_name ILIKE ? OR creator.role ILIKE ? OR updater.full_name ILIKE ?",
			kw, kw, kw, kw).Find(&feedings).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "search feedings failed"})
		}
	} else {
		if err := db.Find(&feedings).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "cannot fetch feedings"})
		}
	}

	var resp []dto.FeedingResponse
	for _, feeding := range feedings {
		resp = append(resp, dto.FeedingResponse{
			ID:          feeding.ID,
			FoodID:      feeding.FoodStock.ID,
			DateTime:    feeding.DateTime,
			Amount:      feeding.Amount,
			Note:        feeding.Note,
			FoodName:    feeding.FoodStock.Name,
			CreatedName: feeding.Creator.FullName,
			CreatedRole: feeding.Creator.Role,
			UpdatedName: feeding.Updater.FullName,
			UpdatedRole: feeding.Updater.Role,
		})
	}

	return c.JSON(resp)

}

func GetAllFeeding(c *fiber.Ctx) error {
	var feedings []models.Feeding
	err := config.DB.Order("date_time DESC").Preload("FoodStock").Preload("Items.Pig").Preload("Creator").Preload("Updater").Find(&feedings).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch Feedings"})
	}
	var resp []dto.FeedingResponse
	for _, feeding := range feedings {
		var pigIDs []uint
		var pigCodeNames []string

		for _, item := range feeding.Items {
			if item.Pig.ID != 0 {
				pigIDs = append(pigIDs, item.Pig.ID)
				pigCodeNames = append(pigCodeNames, item.Pig.CodeName)
			}
		}
		resp = append(resp, dto.FeedingResponse{
			ID:          feeding.ID,
			FoodID:      feeding.FoodStock.ID,
			PigIDs:      pigIDs,
			PigCodeName: pigCodeNames,
			DateTime:    feeding.DateTime,
			Amount:      feeding.Amount,
			Note:        feeding.Note,
			FoodName:    feeding.FoodStock.Name,
			CreatedName: feeding.Creator.FullName,
			CreatedRole: feeding.Creator.Role,
			UpdatedName: feeding.Updater.FullName,
			UpdatedRole: feeding.Updater.Role,
		})
	}
	return c.JSON(resp)
}

func GetFeedingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	feeding_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	feeding := &models.Feeding{}
	if err := config.DB.Preload("FoodStock").Preload("Creator").Preload("Updater").First(feeding, "id = ?", feeding_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Feeding ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error Feeding"})
	}

	resp := dto.FeedingResponse{
		ID:          feeding.ID,
		FoodID:      feeding.FoodStock.ID,
		DateTime:    feeding.DateTime,
		Amount:      feeding.Amount,
		Note:        feeding.Note,
		FoodName:    feeding.FoodStock.Name,
		CreatedName: feeding.Creator.FullName,
		CreatedRole: feeding.Creator.Role,
		UpdatedName: feeding.Updater.FullName,
		UpdatedRole: feeding.Updater.Role,
	}

	return c.JSON(resp)
}
func UpdateFeeding(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userID64, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	userID := uint(userID64)

	id := c.Params("id")
	feedingID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	feeding := &models.Feeding{}
	if err := config.DB.Preload("Items").First(feeding, "id = ?", feedingID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Feeding not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	var input dto.FeedingUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})

	// ----------------- อัปเดต FoodStock -----------------
	newFoodID := feeding.FoodID
	newAmount := feeding.Amount

	if input.FoodID != nil {
		newFoodID = *input.FoodID
	}
	if input.Amount != nil {
		if *input.Amount <= 0 {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Amount must be > 0"})
		}
		newAmount = *input.Amount
	}

	if input.FoodID != nil || input.Amount != nil {
		// คืน stock ของ food เดิม
		oldFood := &models.FoodStock{}
		if err := tx.First(oldFood, "id = ?", feeding.FoodID).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to find old food stock"})
		}
		oldFood.Amount += feeding.Amount
		tx.Save(oldFood)

		// ตัด stock ของ food ใหม่
		newFood := &models.FoodStock{}
		if err := tx.First(newFood, "id = ?", newFoodID).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "FoodStock not found"})
		}
		if newFood.Amount < newAmount {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Not enough food in stock"})
		}
		newFood.Amount -= newAmount
		tx.Save(newFood)

		updates["food_id"] = newFoodID
		updates["amount"] = newAmount
	}

	// ----------------- อัปเดตวันที่และ note -----------------
	if input.DateTime != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		parsedDate, err := time.ParseInLocation("2006-01-02 15:04", *input.DateTime, loc)
		if err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format"})
		}
		if parsedDate.After(time.Now()) {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
		}
		updates["date_time"] = parsedDate
	}
	if input.Note != nil {
		updates["note"] = *input.Note
	}

	updates["updated_by"] = userID

	if err := tx.Model(feeding).Updates(updates).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update feeding"})
	}

	// ----------------- อัปเดต FeedingItem -----------------
	// ลบ item เดิม
	if len(feeding.Items) > 0 {
		if err := tx.Where("feeding_id = ?", feeding.ID).Delete(&models.FeedingItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete old feeding items"})
		}
	}

	var validPigIDs []uint
	var items []models.FeedingItem
	var pigCodeNames []string
	if input.PigIDs != nil && len(*input.PigIDs) > 0 {
		var pigs []models.Pig
		if err := tx.Where("id IN ?", *input.PigIDs).Find(&pigs).Error; err != nil {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Pigs not found"})
		}
		for _, pig := range pigs {
			if pig.Status != "ตายเเล้ว" && pig.Status != "ขายเเล้ว" {
				items = append(items, models.FeedingItem{
					FeedingID: feeding.ID,
					PigID:     pig.ID,
				})
				validPigIDs = append(validPigIDs, pig.ID)
			}
		}

		for _, pig := range pigs {
			if pig.Status != "ตายเเล้ว" && pig.Status != "ขายเเล้ว" {
				pigCodeNames = append(pigCodeNames, pig.CodeName)
			}
		}
	}

	if len(items) > 0 {
		if err := tx.Create(&items).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create feeding items"})
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	// Reload
	if err := config.DB.Preload("FoodStock").Preload("Creator").Preload("Updater").First(feeding, "id = ?", feeding.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reload feeding"})
	}

	resp := dto.FeedingResponse{
		ID:          feeding.ID,
		FoodID:      feeding.FoodStock.ID,
		PigIDs:      validPigIDs,
		PigCodeName: pigCodeNames,
		DateTime:    feeding.DateTime,
		Amount:      feeding.Amount,
		Note:        feeding.Note,
		FoodName:    feeding.FoodStock.Name,
		CreatedName: feeding.Creator.FullName,
		CreatedRole: feeding.Creator.Role,
		UpdatedName: feeding.Updater.FullName,
		UpdatedRole: feeding.Updater.Role,
	}

	return c.JSON(resp)
}

func DeleteFeeding(c *fiber.Ctx) error {
	id := c.Params("id")
	feeding_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	feeding := &models.Feeding{}
	if err := config.DB.Preload("FoodStock").First(feeding, "id = ?", feeding_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Feeding ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error Feeding"})
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// คืน stock ก่อนลบ
	foodStock := &models.FoodStock{}
	if err := tx.First(foodStock, "id = ?", feeding.FoodID).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Food stock not found"})
	}

	foodStock.Amount += feeding.Amount

	if err := tx.Save(foodStock).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to restore food stock"})
	}

	err = tx.Delete(feeding).Error
	if err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete breeding"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.JSON(fiber.Map{"message": "Feeding Deleted",
		"restored_stock": fiber.Map{
			"foodName":        feeding.FoodStock.Name,
			"amount_restored": feeding.Amount,
			"newFoodAmount":   foodStock.Amount}})

}
