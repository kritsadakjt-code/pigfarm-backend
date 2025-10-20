package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// CreateFeedingSchedule สร้างตารางเวลาให้อาหารใหม่พร้อมรายการอาหาร
func CreateFeedingSchedule(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userID64, _ := strconv.ParseUint(userIDStr, 10, 64)
	userID := uint(userID64)

	var input dto.FeedingScheduleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ข้อมูลที่ส่งมาไม่ถูกต้อง"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// เริ่ม Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	schedule := models.FeedingSchedule{
		Name:          input.Name,
		ScheduledTime: input.ScheduledTime,
		IsActive:      input.IsActive,
		Note:          input.Note,
		CreatedBy:     userID,
	}

	// 1. สร้าง Schedule หลัก
	if err := tx.Create(&schedule).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "บันทึกข้อมูลหลักไม่สำเร็จ"})
	}

	// 2. ตรวจสอบอาหารซ้ำ + สร้าง Item
	foodMap := make(map[uint]bool)
	for _, itemInput := range input.Items {
		var stock models.FoodStock
		if err := tx.First(&stock, itemInput.FoodID).Error; err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("ไม่พบข้อมูลอาหาร ID %d", itemInput.FoodID)})
		}
		if itemInput.Amount <= 0 {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "ปริมาณอาหารต้องมากกว่า 0"})
		}
		if itemInput.Amount > stock.Amount {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("ปริมาณอาหาร %s เกินกว่าสต็อกที่มี (เหลือ %.2f)", stock.Name, stock.Amount),
			})
		}
		if foodMap[itemInput.FoodID] {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "มีอาหารซ้ำกันในรายการ"})
		}
		foodMap[itemInput.FoodID] = true

		item := models.FeedingScheduleItem{
			ScheduleID: schedule.ID,
			FoodID:     itemInput.FoodID,
			Amount:     itemInput.Amount,
		}
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "บันทึกข้อมูลรายการอาหารไม่สำเร็จ"})
		}
	}

	// 3. Commit Transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "ไม่สามารถบันทึกข้อมูลได้"})
	}

	// 4. โหลดข้อมูลพร้อม Preload กลับมา
	config.DB.Preload("Items.FoodStock").Preload("Creator").First(&schedule, schedule.ID)

	return c.Status(fiber.StatusCreated).JSON(schedule)
}

// GetAllFeedingSchedules ดึงตารางเวลาทั้งหมด
func GetAllFeedingSchedules(c *fiber.Ctx) error {
	var schedules []models.FeedingSchedule
	if err := config.DB.Preload("Items.FoodStock").Preload("Creator").Order("created_at desc").Find(&schedules).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch schedules"})
	}
	return c.JSON(schedules)
}

// GetFeedingScheduleByID ดึงตารางเวลาตาม ID
func GetFeedingScheduleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var schedule models.FeedingSchedule
	if err := config.DB.Preload("Items.FoodStock").Preload("Creator").First(&schedule, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}
	return c.JSON(schedule)
}

// UpdateFeedingSchedule อัปเดตตารางเวลา
func UpdateFeedingSchedule(c *fiber.Ctx) error {
	id := c.Params("id")
	scheduleID, _ := strconv.Atoi(id)

	var input dto.FeedingScheduleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. ค้นหา Schedule ที่มีอยู่
	var schedule models.FeedingSchedule
	if err := tx.First(&schedule, scheduleID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Schedule not found"})
	}

	// 2. อัปเดตข้อมูลหลักของ Schedule
	schedule.Name = input.Name
	schedule.ScheduledTime = input.ScheduledTime
	schedule.IsActive = input.IsActive
	schedule.Note = input.Note
	if err := tx.Save(&schedule).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update schedule header"})
	}

	// 3. ลบ Items เก่าทั้งหมด
	if err := tx.Where("schedule_id = ?", scheduleID).Delete(&models.FeedingScheduleItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete old items"})
	}

	foodMap := make(map[uint]bool)
	for _, itemInput := range input.Items {
		if itemInput.Amount <= 0 {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "ปริมาณอาหารต้องมากกว่า 0"})
		}
		if foodMap[itemInput.FoodID] {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "มีอาหารซ้ำกันในรายการ"})
		}
		foodMap[itemInput.FoodID] = true

		var stock models.FoodStock
		if err := tx.First(&stock, itemInput.FoodID).Error; err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": fmt.Sprintf("ไม่พบข้อมูลอาหาร ID %d", itemInput.FoodID)})
		}
		if itemInput.Amount > stock.Amount {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("ปริมาณอาหาร %s เกินกว่าสต็อกที่มี (เหลือ %.2f)", stock.Name, stock.Amount),
			})
		}

		item := models.FeedingScheduleItem{
			ScheduleID: uint(scheduleID),
			FoodID:     itemInput.FoodID,
			Amount:     itemInput.Amount,
		}
		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create new items"})
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	config.DB.Preload("Items.FoodStock").Preload("Creator").First(&schedule, schedule.ID)
	return c.JSON(schedule)
}

// DeleteFeedingSchedule ลบตารางเวลา
func DeleteFeedingSchedule(c *fiber.Ctx) error {
	id := c.Params("id")

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. ลบ Items ก่อน
	if err := tx.Where("schedule_id = ?", id).Unscoped().Delete(&models.FeedingScheduleItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete items"})
	}

	// 2. ลบ Schedule หลัก
	result := tx.Unscoped().Delete(&models.FeedingSchedule{}, id)
	if result.Error != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete schedule"})
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Schedule not found"})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Feeding schedule deleted successfully"})
}
