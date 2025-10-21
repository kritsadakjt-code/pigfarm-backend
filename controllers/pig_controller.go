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

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

func CreatePig(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	var input dto.PigInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// 2. ดึงข้อมูลเวลา (สัปดาห์ และ ปี)
	now := time.Now()
	year := now.Format("06") // "06" คือ format สำหรับปี 2 ตัวท้าย (เช่น 25)
	_, week := now.ISOWeek()
	weekStr := fmt.Sprintf("%02d", week) // "%02d" คือ format ให้มี 0 นำหน้าถ้าเป็นเลขหลักเดียว

	// 3. ค้นหาลำดับเลขถัดไป (001-999)
	var lastPig models.Pig
	// สร้าง Pattern สำหรับค้นหา เช่น "2001-4225" (หาหมูพันธุ์ดูร็อกที่เกิดในสัปดาห์ที่ 42 ปี 25)
	pattern := fmt.Sprintf("%s%%-%s%s", input.CodePrefix, weekStr, year)

	// ค้นหาหมูตัวล่าสุดที่ตรงกับ Pattern นี้ โดยเรียงจาก code_name มากไปน้อย
	err = config.DB.Where("code_name LIKE ?", pattern).Order("code_name DESC").First(&lastPig).Error

	nextNumber := 1
	if err == nil { // ถ้าเจอ record ล่าสุด
		var lastNum int
		// ดึงตัวเลขลำดับจากรหัสล่าสุด เช่น "D-005-42-25" จะได้เลข 5
		fmt.Sscanf(lastPig.CodeName, input.CodePrefix+"%d", &lastNum)
		nextNumber = lastNum + 1 // บวก 1 เพื่อเป็นเลขถัดไป
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// ถ้าเกิด Error อื่นๆ ที่ไม่ใช่ "หาไม่เจอ"
		return c.Status(500).JSON(fiber.Map{"error": "Failed to query for the last pig code"})
	}
	// ถ้าหาไม่เจอ (gorm.ErrRecordNotFound) nextNumber จะเป็น 1 ซึ่งถูกต้องแล้ว

	newCodeName := fmt.Sprintf("%s%03d-%s%s", input.CodePrefix, nextNumber, weekStr, year)

	if input.Gender == "ผู้" && input.Type == "เเม่พันธุ์" {
		return c.Status(400).JSON(fiber.Map{"error": "เพศผู้ไม่สามารถเป็นเเม่พันธุ์ได้"})
	}
	if input.Gender == "เมีย" && input.Type == "พ่อพันธุ์" {
		return c.Status(400).JSON(fiber.Map{"error": "เพศเมียไม่สามารถเป็นพ่อพันธุ์ได้"})
	}

	if input.Gender == "ผู้" {
		if input.Status == "อุ้มท้อง" || input.Status == "ให้นมลูก" {
			return c.Status(400).JSON(fiber.Map{"error": "เพศผู้ไม่สามารถมีสถานะอุ้มท้อง, ให้นมลูกได้"})
		}
	}

	if input.Type == "ลูกหมู" {
		if input.Status == "อุ้มท้อง" || input.Status == "ให้นมลูก" || input.Status == "พร้อมผสม" {
			return c.Status(400).JSON(fiber.Map{"error": "ลูกหมูไม่สามารถมีสถานะอุ้มท้อง, ให้นมลูก, หรือพร้อมผสมได้"})
		}
	}
	if input.Type == "หมูขุน" {
		if input.Status == "อุ้มท้อง" || input.Status == "ให้นมลูก" || input.Status == "พร้อมผสม" {
			return c.Status(400).JSON(fiber.Map{"error": "ลูกหมูไม่สามารถมีสถานะอุ้มท้อง, ให้นมลูก, หรือพร้อมผสมได้"})
		}
	}

	if input.Weight <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "น้ำหนักต้องมากกว่า 0"})
	}

	status := input.Status
	switch input.Type {
	case "ลูกหมู":
		status = "กําลังเลี้ยง"
	case "หมูขุน":
		status = "กําลังขุน"
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	parsedDate, err := time.ParseInLocation("2006-01-02", input.BirthDate, loc)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Birth date cannot be in the future"})
	}

	pig := models.Pig{
		CodeName:  newCodeName,
		Breed:     input.Breed,
		Gender:    input.Gender,
		Type:      input.Type,
		BirthDate: parsedDate,
		Weight:    input.Weight,
		Status:    status,
		CreatedBy: user_id,
		UpdatedBy: user_id,
	}

	if err := config.DB.Create(&pig).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return c.Status(400).JSON(fiber.Map{"error": "CodeName already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create pig"})
	}

	if err := config.DB.Preload("Creator").Preload("Updater").First(&pig, "id = ?", pig.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to preload"})
	}

	resp := dto.PigResponse{
		ID:          pig.ID,
		CodeName:    pig.CodeName,
		Breed:       pig.Breed,
		Gender:      pig.Gender,
		Type:        pig.Type,
		BirthDate:   pig.BirthDate,
		Weight:      pig.Weight,
		Status:      pig.Status,
		CreatedName: pig.Creator.FullName,
		UpdatedName: pig.Updater.FullName,
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Pig created successfully",
		"data":    resp,
	})
}

func SearchPigs(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("keyword"))
	pigType := strings.TrimSpace(c.Query("type"))  // "พ่อพันธุ์" หรือ "เเม่พันธุ์"
	status := strings.TrimSpace(c.Query("status")) // "พร้อมผสม"

	db := config.DB.Model(&models.Pig{}).
		Preload("Creator").Preload("Updater").
		Joins("LEFT JOIN users creator on creator.id = pigs.created_by").
		Joins("LEFT JOIN users updater on updater.id = pigs.updated_by")

	// filters
	if pigType != "" {
		db = db.Where("pigs.type = ?", pigType)
	}
	if status != "" {
		db = db.Where("pigs.status = ?", status)
	}
	if keyword != "" {
		kw := "%" + keyword + "%"
		db = db.Where("pigs.code_name ILIKE ? OR pigs.breed ILIKE ? OR creator.role ILIKE ? OR updater.role ILIKE ?", kw, kw, kw, kw)
	}

	var pigs []models.Pig
	if err := db.Find(&pigs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch pigs"})
	}

	var resp []dto.PigResponse
	for _, pig := range pigs {
		resp = append(resp, dto.PigResponse{
			ID:          pig.ID,
			CodeName:    pig.CodeName,
			Breed:       pig.Breed,
			Gender:      pig.Gender,
			Type:        pig.Type,
			BirthDate:   pig.BirthDate,
			Weight:      pig.Weight,
			Status:      pig.Status,
			CreatedName: pig.Creator.FullName,
			CreatedRole: pig.Creator.Role,
			UpdatedName: pig.Updater.FullName,
			UpdatedRole: pig.Updater.Role,
		})
	}

	return c.JSON(resp)
}

func GetAllPigs(c *fiber.Ctx) error {
	var pigs []models.Pig
	if err := config.DB.Preload("Creator").Preload("Updater").Find(&pigs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get pigs"})
	}

	var resp []dto.PigResponse
	for _, pig := range pigs {
		resp = append(resp, dto.PigResponse{
			ID:          pig.ID,
			CodeName:    pig.CodeName,
			Breed:       pig.Breed,
			Gender:      pig.Gender,
			Type:        pig.Type,
			BirthDate:   pig.BirthDate,
			Weight:      pig.Weight,
			Status:      pig.Status,
			CreatedName: pig.Creator.FullName,
			CreatedRole: pig.Creator.Role,
			UpdatedName: pig.Updater.FullName,
			UpdatedRole: pig.Updater.Role,
		})
	}
	return c.JSON(resp)
}

func GetPigByID(c *fiber.Ctx) error {
	id := c.Params("id")
	pig := &models.Pig{}
	if err := config.DB.Preload("Creator").Preload("Updater").First(pig, "id = ?", id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Pig not found"})
	}

	resp := dto.PigResponse{
		ID:          pig.ID,
		CodeName:    pig.CodeName,
		Breed:       pig.Breed,
		Gender:      pig.Gender,
		Type:        pig.Type,
		BirthDate:   pig.BirthDate,
		Weight:      pig.Weight,
		Status:      pig.Status,
		CreatedName: pig.Creator.FullName,
		CreatedRole: pig.Creator.Role,
		UpdatedName: pig.Updater.FullName,
		UpdatedRole: pig.Updater.Role,
	}
	return c.JSON(resp)
}

func UpdatePig(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, _ := strconv.ParseUint(user_id_str, 10, 64)
	user_id := uint(user_id64)

	id := c.Params("id")
	pig_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid id"})
	}

	pig := &models.Pig{}
	if err := config.DB.First(pig, "id = ?", pig_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Pig not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	var input dto.PigUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := make(map[string]interface{})

	if *input.Gender == "ผู้" && *input.Type == "เเม่พันธุ์" {
		return c.Status(400).JSON(fiber.Map{"error": "เพศผู้ไม่สามารถเป็นเเม่พันธุ์ได้"})
	}
	if *input.Gender == "เมีย" && *input.Type == "พ่อพันธุ์" {
		return c.Status(400).JSON(fiber.Map{"error": "เพศเมียไม่สามารถเป็นพ่อพันธุ์ได้"})
	}
	if *input.Gender == "ผู้" {
		if *input.Status == "อุ้มท้อง" || *input.Status == "ให้นมลูก" {
			return c.Status(400).JSON(fiber.Map{"error": "เพศผู้ไม่สามารถมีสถานะอุ้มท้อง, ให้นมลูกได้"})
		}
	}

	if *input.Type == "ลูกหมู" {
		if *input.Status == "อุ้มท้อง" || *input.Status == "ให้นมลูก" || *input.Status == "พร้อมผสม" {
			return c.Status(400).JSON(fiber.Map{"error": "ลูกหมูไม่สามารถมีสถานะอุ้มท้อง, ให้นมลูก, หรือพร้อมผสมได้"})
		}
	}
	if *input.Type == "หมูขุน" {
		if *input.Status == "อุ้มท้อง" || *input.Status == "ให้นมลูก" || *input.Status == "พร้อมผสม" {
			return c.Status(400).JSON(fiber.Map{"error": "ลูกหมูไม่สามารถมีสถานะอุ้มท้อง, ให้นมลูก, หรือพร้อมผสมได้"})
		}
	}

	// if input.CodeName != nil {
	// 	// ตรวจสอบก็ต่อเมื่อ code_name ใหม่ที่ส่งมาไม่ตรงกับของเดิมเท่านั้น

	// 	if *input.CodeName != pig.CodeName {
	// 		var existingPig models.Pig
	// 		// ค้นหาหมูที่มี code_name ตรงกับที่ส่งมา และ ID ไม่ใช่ตัวที่กำลังแก้ไข
	// 		err := config.DB.Where("code_name = ? AND id != ?", *input.CodeName, pig_id).First(&existingPig).Error

	// 		if err == nil {
	// 			// ถ้า err เป็น nil = เจอ
	// 			return c.Status(400).JSON(fiber.Map{"error": "รหัสหมูนี้ถูกใช้โดยหมูตัวอื่นแล้ว"})
	// 		}

	// 		// ถ้า err ไม่ใช่ gorm.ErrRecordNotFound แสดงว่าเป็น error อื่นๆ ของฐานข้อมูล
	// 		if !errors.Is(err, gorm.ErrRecordNotFound) {
	// 			return c.Status(500).JSON(fiber.Map{"error": "Database error while checking for duplicate code name"})
	// 		}

	// 	}
	// 	updates["code_name"] = *input.CodeName
	// }

	if input.Breed != nil {
		updates["breed"] = *input.Breed
	}

	if input.Gender != nil {
		updates["gender"] = *input.Gender
	}
	if input.Type != nil {

		updates["type"] = *input.Type
	}
	if input.BirthDate != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		parsedDate, err := time.ParseInLocation("2006-01-02", *input.BirthDate, loc)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}
		if parsedDate.After(time.Now()) {
			return c.Status(400).JSON(fiber.Map{"error": "Birth date cannot be in the future"})
		}
		updates["birth_date"] = parsedDate
	}
	if input.Weight != nil {
		if *input.Weight <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Weight must be greater than 0"})
		}
		updates["weight"] = *input.Weight
	}
	if input.Status != nil {
		updates["status"] = *input.Status
	}

	updates["updated_by"] = user_id // อัปเดต Foreign Key

	if err := config.DB.Model(pig).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update pig"})
	}

	if err := config.DB.Preload("Updater").First(&pig, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to preload"})
	}

	resp := dto.PigResponse{
		ID:          pig.ID,
		CodeName:    pig.CodeName,
		Breed:       pig.Breed,
		Gender:      pig.Gender,
		Type:        pig.Type,
		BirthDate:   pig.BirthDate,
		Weight:      pig.Weight,
		Status:      pig.Status,
		UpdatedName: pig.Updater.FullName,
		UpdatedRole: pig.Updater.Role,
	}

	return c.JSON(fiber.Map{
		"message": "Pig updated successfully",
		"data":    resp,
	})
}

func DeletePig(c *fiber.Ctx) error {
	id := c.Params("id")
	pig_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	pig := &models.Pig{}
	if err := config.DB.First(pig, "id = ?", pig_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Pig not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	err = config.DB.Unscoped().Delete(pig).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to delete pig"})
	}

	return c.JSON(fiber.Map{"message": "Pig deleted"})
}
