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

func CreateFoodStock(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	var input dto.FoodStockInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	// เช็คก่อนหรือหลัง create ก็ได้
	var count int64
	if err := config.DB.Model(&models.FoodStock{}).Where("name = ?", input.Name).Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Food name already exists"})
	}
	// fix ไว้
	// var FoodType string
	// switch input.Name {
	// case "อาหารลูกหมู", "อาหารหมูขุน", "อาหารพ่อเเม่พันธุ์", "อาหารหมูท้อง", "อาหารหมูให้นม":
	// 	FoodType = "อาหารหลัก"
	// case "โพรไบโอติก", "เเร่ธาตุรวม", "วิตามินรวม":
	// 	FoodType = "อาหารเสริม"
	// }
	// loc, _ := time.LoadLocation("Asia/Bangkok")
	loc, _ := time.LoadLocation("Asia/Bangkok")

	if input.DateTime == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "DateTime is required",
		})
	}
	parsedDate, err := time.ParseInLocation("2006-01-02 15:04", input.DateTime, loc)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD HH:MM"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
	}
	if input.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "amount must be greater than 0"})
	}

	foodStock := models.FoodStock{
		Name:      input.Name,
		Type:      input.Type,
		DateTime:  parsedDate,
		Amount:    input.Amount,
		Note:      input.Note,
		CreatedBy: user_id,
		UpdatedBy: user_id,
	}

	if err := config.DB.Create(&foodStock).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(400).JSON(fiber.Map{"error": "Food name already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "failed to create FoodStock"})
	}

	if err := config.DB.Preload("Creator").Preload("Updater").First(&foodStock, "id = ?", foodStock.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch preload FoodStock"})
	}

	resp := dto.FoodStockResponse{
		ID:          foodStock.ID,
		Name:        foodStock.Name,
		Type:        foodStock.Type,
		DateTime:    foodStock.DateTime,
		Amount:      foodStock.Amount,
		Note:        foodStock.Note,
		CreatedName: foodStock.Creator.FullName,
		UpdatedName: foodStock.Updater.FullName,
	}

	return c.Status(201).JSON(resp)
}

func SearchFoodStock(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	var foodStock []models.FoodStock

	db := config.DB.Model(&models.FoodStock{}).Preload("Creator").Preload("Updater").
		Joins("LEFT JOIN users creator on creator.id = food_stocks.created_by").
		Joins("LEFT JOIN users updater on updater.id = food_stocks.updated_by")
	kw := "%" + keyword + "%"

	if keyword != "" {
		if err := db.Where("food_stocks.name ILIKE ? OR food_stocks.type ILIKE ? OR creator.role ILIKE ?", kw, kw, kw).Find(&foodStock).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "search foodstock failed"})
		}
	} else {
		err := db.Find(&foodStock).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "cannot fetch FoodStock"})
		}
	}

	var resp []dto.FoodStockResponse
	for _, food := range foodStock {
		resp = append(resp, dto.FoodStockResponse{
			ID:          food.ID,
			Name:        food.Name,
			Type:        food.Type,
			DateTime:    food.DateTime,
			Amount:      food.Amount,
			Note:        food.Note,
			CreatedName: food.Creator.FullName,
			CreatedRole: food.Creator.Role,
			UpdatedName: food.Updater.FullName,
			UpdatedRole: food.Updater.Role,
		})
	}
	return c.JSON(resp)
}

func GetAllFoodStock(c *fiber.Ctx) error {
	var foodStock []models.FoodStock
	if err := config.DB.Preload("Creator").Preload("Updater").Find(&foodStock).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch FoodStock"})
	}

	var resp []dto.FoodStockResponse
	for _, food := range foodStock {
		resp = append(resp, dto.FoodStockResponse{
			ID:          food.ID,
			Name:        food.Name,
			Type:        food.Type,
			DateTime:    food.DateTime,
			Amount:      food.Amount,
			Note:        food.Note,
			CreatedName: food.Creator.FullName,
			CreatedRole: food.Creator.Role,
			UpdatedName: food.Updater.FullName,
			UpdatedRole: food.Updater.Role,
		})
	}
	return c.JSON(resp)

}

func GetFoodStockByID(c *fiber.Ctx) error {
	id := c.Params("id")
	stock_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ID"})
	}
	foodStock := &models.FoodStock{}
	if err := config.DB.Preload("Creator").Preload("Updater").First(&foodStock, "id = ?", stock_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not Found FoodStock ID"})
	}

	resp := dto.FoodStockResponse{
		ID:          foodStock.ID,
		Name:        foodStock.Name,
		Type:        foodStock.Type,
		DateTime:    foodStock.DateTime,
		Amount:      foodStock.Amount,
		Note:        foodStock.Note,
		CreatedName: foodStock.Creator.FullName,
		CreatedRole: foodStock.Creator.Role,
		UpdatedName: foodStock.Updater.FullName,
		UpdatedRole: foodStock.Updater.Role,
	}
	return c.JSON(resp)
}

func UpdateFoodStock(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	id := c.Params("id")
	stock_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ID"})
	}
	foodStock := &models.FoodStock{}
	if err := config.DB.First(&foodStock, "id = ?", stock_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found FoodStock ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	var input dto.FoodStockUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	update := make(map[string]interface{})
	// ไม่ให้เเก้ name type เพราะจะไปกระทบการให้อาหาร
	// if input.Name != nil {
	// 	var count int64
	// 	// <> ไม่เอา record ตัวเองมาคํานวณด้วยจึงจะ update ชื่อเดิมได้
	// 	if err := config.DB.Model(&models.FoodStock{}).Where("name = ? AND id <> ?", input.Name, stock_id).Count(&count).Error; err != nil {
	// 		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	// 	}
	// 	if count > 0 {
	// 		return c.Status(400).JSON(fiber.Map{"error": "Food name already exists"})
	// 	}
	// 	update["name"] = *input.Name
	// }
	// if input.Type != nil {
	// 	update["type"] = *input.Type
	// }
	if input.DateTime != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		parsedDate, err := time.ParseInLocation("2006-01-02 15:04", *input.DateTime, loc)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD HH:MM"})
		}
		if parsedDate.After(time.Now()) {
			return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
		}
		update["date_time"] = parsedDate
	}
	if input.Amount != nil {
		if *input.Amount <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "amount must be greater than 0"})
		}
		update["amount"] = *input.Amount
	}
	if input.Note != nil {
		update["note"] = *input.Note
	}
	update["updated_by"] = user_id
	if err := config.DB.Model(foodStock).Updates(update).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(400).JSON(fiber.Map{"error": "Food name already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Fail to update FoodStock"})
	}

	if err := config.DB.Preload("Creator").Preload("Updater").First(foodStock, "id = ?", stock_id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to preload food stock"})
	}

	resp := dto.FoodStockResponse{
		ID:          foodStock.ID,
		Name:        foodStock.Name,
		Type:        foodStock.Type,
		DateTime:    foodStock.DateTime,
		Amount:      foodStock.Amount,
		Note:        foodStock.Note,
		CreatedName: foodStock.Creator.FullName,
		CreatedRole: foodStock.Creator.Role,
		UpdatedName: foodStock.Updater.FullName,
		UpdatedRole: foodStock.Updater.Role,
	}
	return c.JSON(resp)

}

func DeleteFoodStock(c *fiber.Ctx) error {
	id := c.Params("id")
	stock_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ID"})
	}
	foodStock := &models.FoodStock{}
	if err := config.DB.First(&foodStock, "id = ?", stock_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found FoodStock ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	err = config.DB.Unscoped().Delete(&foodStock).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete FoodStock"})
	}
	return c.JSON(fiber.Map{"message": "FoodStock Deleted"})
}
