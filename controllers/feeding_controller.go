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

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Transaction commit failed"})
	}

	if err := config.DB.Preload("FoodStock").Preload("Creator").Preload("Updater").First(&feeding, "id = ?", feeding.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to preload feeding"})
	}

	resp := dto.FeedingResponse{
		ID:          feeding.ID,
		FoodID:      food.ID,
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
	err := config.DB.Preload("FoodStock").Preload("Creator").Preload("Updater").Find(&feedings).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch Feedings"})
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
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	id := c.Params("id")
	feeding_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	feeding := &models.Feeding{}
	if err := config.DB.First(feeding, "id = ?", feeding_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Feeding ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error Feeding"})
	}

	var input dto.FeedingUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// เริ่ม transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	updates := make(map[string]interface{})

	// ให้ให้เป็นค่าเดิมถ้าไม่มีการเปลี่ยนค่า เพราะจะได้เอาไปใช้ถูก
	// ค่า food ใหม่เเละเก่า จะใช้ตัวเเปรเดียวกัน
	newFoodID := feeding.FoodID
	newAmount := feeding.Amount

	// ถ้าเปลี่ยนค่า
	if input.FoodID != nil {
		newFoodID = *input.FoodID
	}
	if input.Amount != nil {
		if *input.Amount <= 0 {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Amount must be non-negative"})
		}
		newAmount = *input.Amount
	}

	// ถ้ามีการเปลี่ยน foodID or Amount
	stockUpdate := (input.FoodID != nil) || (input.Amount != nil)
	if stockUpdate {
		// คืนค่าของ foodID เดิม
		oldFood := &models.FoodStock{}
		if err := tx.First(oldFood, "id = ?", feeding.FoodID).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to find old food stock"})
		}
		oldFood.Amount += feeding.Amount
		if err := tx.Save(oldFood).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to restore old food stock"})
		}

		// ตัดสต็อกจะตัดตามค่า newFoodID ถ้าไม่มีค่าใหม่ก็ใช้ตัว newFoodID ซึ่งจะเป็นได้ทั้งค่าใหม่เเละค่าเก่า
		newFood := &models.FoodStock{}
		if err := tx.First(newFood, "id = ?", newFoodID).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(404).JSON(fiber.Map{"error": "FoodStock not found"})
			}
			return c.Status(500).JSON(fiber.Map{"error": "Database error (food stock)"})
		}

		if newFood.Amount < newAmount {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Not enough food in stock"})
		}

		newFood.Amount -= newAmount
		if err := tx.Save(newFood).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update new food stock"})
		}

		// อัพเดท feeding
		if input.FoodID != nil {
			updates["food_id"] = newFoodID
		}
		if input.Amount != nil {
			updates["amount"] = newAmount
		}

	}
	if input.DateTime != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		parsedDate, err := time.ParseInLocation("2006-01-02 15:04", *input.DateTime, loc)
		if err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD HH:MM"})
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
	updates["updated_by"] = user_id

	if err := tx.Model(feeding).Updates(updates).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Fail to update feeding"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	if err := config.DB.Preload("FoodStock").Preload("Creator").Preload("Updater").
		First(feeding, "id = ?", feeding.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reload feeding"})
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
