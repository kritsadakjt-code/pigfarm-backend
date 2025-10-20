package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateHealth(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	var input dto.HealthInput
	if err := c.BodyParser(&input); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	if err := validate.Struct(input); err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// ตรวจสอบว่าหมูมีอยู่จริง
	pig := &models.Pig{}
	if err := config.DB.First(pig, "id = ?", input.PigID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found pig id"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
	}

	health := &models.Health{
		PigID:     input.PigID,
		Date:      parsedDate,
		Type:      input.Type,
		Detail:    input.Detail,
		Note:      input.Note,
		CreatedBy: user_id,
		UpdatedBy: user_id,
	}

	if err := config.DB.Create(health).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create Health"})
	}

	if err := config.DB.Preload("Pig").Preload("Creator").Preload("Updater").
		First(health, "id = ?", health.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to load health"})
	}

	resp := dto.HealthResponse{
		ID:          health.ID,
		PigID:       health.PigID,
		PigCodeName: health.Pig.CodeName,
		Date:        health.Date,
		Type:        health.Type,
		Detail:      health.Detail,
		Note:        health.Note,
		CreatedName: health.Creator.FullName,
		CreatedRole: health.Creator.Role,
		UpdatedName: health.Updater.FullName,
		UpdatedRole: health.Updater.Role,
	}

	return c.Status(201).JSON(resp)
}

func SearchHealth(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	var healths []models.Health
	db := config.DB.Model(&models.Health{}).Preload("Pig").Preload("Creator").Preload("Updater").
		Joins("LEFT JOIN users creator on creator.id = healths.created_by").
		Joins("LEFT JOIN users updater on updater.id = healths.updated_by").
		Joins("LEFT JOIN pigs pig on pig.id = healths.pig_id")
	if keyword != "" {
		kw := "%" + keyword + "%"
		if err := db.Where("healths.type ILIKE ? OR pig.code_name ILIKE ? OR creator.full_name ILIKE ? OR updater.role ILIKE ?",
			kw, kw, kw, kw).Find(&healths).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "search healths failed"})
		}
	} else {
		if err := db.Find(&healths).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "failed to fetch healths"})
		}
	}

	var resp []dto.HealthResponse
	for _, health := range healths {
		resp = append(resp, dto.HealthResponse{
			ID:          health.ID,
			PigID:       health.PigID,
			Date:        health.Date,
			Type:        health.Type,
			Detail:      health.Detail,
			Note:        health.Note,
			CreatedName: health.Creator.FullName,
			CreatedRole: health.Creator.Role,
			UpdatedName: health.Updater.FullName,
			UpdatedRole: health.Updater.Role,
		})
	}

	return c.JSON(resp)

}

func GetAllHealth(c *fiber.Ctx) error {
	var healths []models.Health
	err := config.DB.Preload("Pig").Preload("Creator").Preload("Updater").Find(&healths).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch Health"})
	}

	var resp []dto.HealthResponse
	for _, health := range healths {
		resp = append(resp, dto.HealthResponse{
			ID:          health.ID,
			PigID:       health.PigID,
			PigCodeName: health.Pig.CodeName,
			Date:        health.Date,
			Type:        health.Type,
			Detail:      health.Detail,
			Note:        health.Note,
			CreatedName: health.Creator.FullName,
			CreatedRole: health.Creator.Role,
			UpdatedName: health.Updater.FullName,
			UpdatedRole: health.Updater.Role,
		})
	}

	return c.JSON(resp)
}

func GetHealthByID(c *fiber.Ctx) error {
	id := c.Params("id")
	healthID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	health := &models.Health{}
	if err := config.DB.Preload("Pig").Preload("Creator").Preload("Updater").First(health, "id = ?", healthID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Health ID"})
		}
		return c.Status(404).JSON(fiber.Map{"error": "Database error Health"})
	}
	resp := dto.HealthResponse{
		ID:          health.ID,
		PigID:       health.PigID,
		Date:        health.Date,
		Type:        health.Type,
		Detail:      health.Detail,
		Note:        health.Note,
		CreatedName: health.Creator.FullName,
		CreatedRole: health.Creator.Role,
		UpdatedName: health.Updater.FullName,
		UpdatedRole: health.Updater.Role,
	}

	return c.JSON(resp)
}

func UpdateHealth(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	id := c.Params("id")
	healthID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	health := &models.Health{}
	if err := config.DB.First(health, "id = ?", healthID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Health ID"})
		}
		return c.Status(404).JSON(fiber.Map{"error": "Database error Health"})
	}

	var input dto.HealthUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := make(map[string]interface{})

	if input.PigID != nil {
		pig := &models.Pig{}
		if err := config.DB.First(pig, "id = ?", input.PigID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(404).JSON(fiber.Map{"error": "Not Found Pig ID"})
			}
			return c.Status(404).JSON(fiber.Map{"error": "Database error Pig"})
		}
		updates["pig_id"] = *input.PigID
	}
	if input.Date != nil {
		parsedDate, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}
		if parsedDate.After(time.Now()) {
			return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
		}
		updates["date"] = parsedDate
	}

	if input.Type != nil {
		updates["type"] = *input.Type
	}
	if input.Detail != nil {
		updates["detail"] = *input.Detail
	}
	if input.Note != nil {
		updates["note"] = *input.Note
	}
	updates["updated_by"] = user_id
	if err := config.DB.Model(health).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Fail to update Health"})
	}
	if err := config.DB.Preload("Pig").Preload("Creator").Preload("Updater").First(health, "id = ?", healthID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Fail to update Health"})
	}
	resp := dto.HealthResponse{
		ID:          health.ID,
		PigID:       health.PigID,
		PigCodeName: health.Pig.CodeName,
		Date:        health.Date,
		Type:        health.Type,
		Detail:      health.Detail,
		Note:        health.Note,
		CreatedName: health.Creator.FullName,
		CreatedRole: health.Creator.Role,
		UpdatedName: health.Updater.FullName,
		UpdatedRole: health.Updater.Role,
	}

	return c.JSON(resp)
}

func DeleteHealth(c *fiber.Ctx) error {
	id := c.Params("id")
	healthID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	health := &models.Health{}
	if err := config.DB.First(health, "id = ?", healthID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Health ID"})
		}
		return c.Status(404).JSON(fiber.Map{"error": "Database error Health"})
	}

	err = config.DB.Unscoped().Delete(health).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete Health"})
	}
	return c.JSON(fiber.Map{"message": "Health Deleted"})
}
