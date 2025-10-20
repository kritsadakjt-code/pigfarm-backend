package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreatePigSale(c *fiber.Ctx) error {
	var input dto.PigSaleInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// ดึงข้อมูลหมูที่เลือก
	var pigs []models.Pig
	if err := tx.Where("id IN ?", input.PigIDs).Find(&pigs).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Pigs not found"})
	}
	if len(pigs) == 0 {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "No pigs found for given IDs"})
	}

	// เเยกหมูพร้อมขายเเละไม่พร้อม
	var readyPigs, notReadyPigs []models.Pig
	for _, pig := range pigs {
		if pig.Status == "พร้อมขาย" {
			readyPigs = append(readyPigs, pig)
		} else {
			notReadyPigs = append(notReadyPigs, pig)
		}
	}

	// ไม่มีหมูตัวไหนพร้อมขาย
	if len(readyPigs) == 0 {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "No pigs are ready for sale"})
	}

	// auto gen saleCode
	dateStr := parsedDate.Format("20060102")

	var lastSale models.PigSale
	err = tx.Where("date::date = ?", parsedDate.Format("2006-01-02")).Order("id DESC").First(&lastSale).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate sale code"})
	}
	seq := 1
	if lastSale.SaleCode != "" {
		var lastSeq int
		fmt.Sscanf(lastSale.SaleCode, "SALE-"+dateStr+"-%03d", &lastSeq)
		seq = lastSeq + 1
	}
	saleCode := fmt.Sprintf("SALE-%s-%03d", dateStr, seq)
	// สร้างเฉพาะหมูที่พร้อมขาย
	sale := &models.PigSale{
		SaleCode: saleCode,
		// Pigs:       readyPigs,
		Date:       parsedDate,
		TotalPrice: input.TotalPrice,
		Amount:     len(readyPigs),
		Buyer:      input.Buyer,
		Note:       input.Note,
	}

	if err := tx.Create(&sale).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create Pig Sale"})
	}

	// Save PigSaleItems
	var items []models.PigSaleItem
	for _, pig := range readyPigs {
		items = append(items, models.PigSaleItem{
			PigSaleID: sale.ID,
			PigID:     pig.ID,
		})
	}
	if err := tx.Create(&items).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create Pig Sale Items"})
	}

	// อัปเดต status เฉพาะหมูที่พร้อมขาย
	var readyIDs []uint
	for _, p := range readyPigs {
		readyIDs = append(readyIDs, p.ID)
	}
	if err := tx.Model(&models.Pig{}).Where("id IN ?", readyIDs).Update("status", "ขายเเล้ว").Error; err != nil {
		tx.Rollback()
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update pig status"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Transaction commit failed"})
	}

	pigIDs := []uint{}
	pigCodeNames := []string{}
	for _, pig := range readyPigs {
		pigIDs = append(pigIDs, pig.ID)
		pigCodeNames = append(pigCodeNames, pig.CodeName)
	}

	notReadyNames := []string{}
	for _, pig := range notReadyPigs {
		notReadyNames = append(notReadyNames, pig.CodeName)
	}

	resp := dto.PigSaleResponse{
		ID:          sale.ID,
		SaleCode:    sale.SaleCode,
		PigIDs:      pigIDs,
		PigCodeName: pigCodeNames,
		Date:        sale.Date,
		Amount:      len(readyPigs),
		TotalPrice:  sale.TotalPrice,
		Buyer:       sale.Buyer,
		Note:        sale.Note,
	}
	// ต่อ note หมูที่ไม่พร้อมขาย
	if len(notReadyPigs) > 0 {

		resp.Note += " | Pigs not ready : " + strings.Join(notReadyNames, ",")
	}

	return c.Status(201).JSON(resp)

}

func SearchPigSale(c *fiber.Ctx) error {
	startStr := strings.TrimSpace(c.Query("start"))
	endStr := strings.TrimSpace(c.Query("end"))
	saleCode := strings.TrimSpace(c.Query("salecode"))
	buyer := strings.TrimSpace(c.Query("buyer"))
	note := strings.TrimSpace(c.Query("note"))
	var pigSale []models.PigSale
	db := config.DB.Model(&models.PigSale{}).Preload("Items.Pig")

	// ส่งมาแค่ start หรือ end อย่างเดียว
	if (startStr == "") != (endStr == "") {
		return c.Status(400).JSON(fiber.Map{"error": "Both start and end date must be provided"})
	}

	// มี start เเละ end
	if startStr != "" && endStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startStr)
		endDate, err2 := time.Parse("2006-01-02", endStr)
		if err1 != nil || err2 != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}

		// ให้วลาครอบคลุมตั้งแต่ 00:00:00 - 23:59:59
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, time.UTC)
		db = db.Where("date BETWEEN ? AND ?", startDate, endDate)

	}

	if saleCode != "" {
		db = db.Where("sale_code ILIKE ?", "%"+saleCode+"%")
	}
	if buyer != "" {
		db = db.Where("buyer ILIKE ?", "%"+buyer+"%")
	}
	if note != "" {
		db = db.Where("note ILIKE ?", "%"+note+"%")
	}

	db = db.Order("date DESC")

	if err := db.Find(&pigSale).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	var resp []dto.PigSaleResponse
	for _, p := range pigSale {
		var pigIDs []uint
		var pigCodeNames []string

		for _, item := range p.Items {
			if item.Pig.ID != 0 {
				pigIDs = append(pigIDs, item.Pig.ID)
				pigCodeNames = append(pigCodeNames, item.Pig.CodeName)
			}
		}

		resp = append(resp, dto.PigSaleResponse{
			ID:          p.ID,
			SaleCode:    p.SaleCode,
			PigIDs:      pigIDs,
			PigCodeName: pigCodeNames,
			Date:        p.Date,
			Amount:      len(p.Items),
			TotalPrice:  p.TotalPrice,
			Buyer:       p.Buyer,
			Note:        p.Note,
		})
	}

	return c.JSON(resp)

}

func GetAllPigSale(c *fiber.Ctx) error {
	var pigSale []models.PigSale
	if err := config.DB.Preload("Items.Pig").Find(&pigSale).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Fail to fecth PigSale"})
	}

	var resp []dto.PigSaleResponse
	for _, p := range pigSale {
		var pigIDs []uint
		var pigCodeNames []string

		for _, item := range p.Items {
			if item.Pig.ID != 0 {
				pigIDs = append(pigIDs, item.Pig.ID)
				pigCodeNames = append(pigCodeNames, item.Pig.CodeName)
			}
		}

		resp = append(resp, dto.PigSaleResponse{
			ID:          p.ID,
			SaleCode:    p.SaleCode,
			PigIDs:      pigIDs,
			PigCodeName: pigCodeNames,
			Date:        p.Date,
			Amount:      len(p.Items),
			TotalPrice:  p.TotalPrice,
			Buyer:       p.Buyer,
			Note:        p.Note,
		})
	}

	return c.JSON(resp)
}

func GetPigSaleByID(c *fiber.Ctx) error {
	id := c.Params("id")
	pigSale_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	pigSale := &models.PigSale{}
	if err := config.DB.Preload("Items.Pig").First(pigSale, "id = ?", pigSale_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "not found pig sale ID"})
	}

	var pigIDs []uint
	var pigCodeNames []string
	for _, item := range pigSale.Items {
		if item.Pig.ID != 0 {
			pigIDs = append(pigIDs, item.Pig.ID)
			pigCodeNames = append(pigCodeNames, item.Pig.CodeName)
		}
	}

	resp := dto.PigSaleResponse{
		ID:          pigSale.ID,
		SaleCode:    pigSale.SaleCode,
		PigIDs:      pigIDs,
		PigCodeName: pigCodeNames,
		Date:        pigSale.Date,
		Amount:      len(pigSale.Items),
		TotalPrice:  pigSale.TotalPrice,
		Buyer:       pigSale.Buyer,
		Note:        pigSale.Note,
	}

	return c.JSON(resp)

}

func UpdatePigSale(c *fiber.Ctx) error {
	id := c.Params("id")
	pigSaleID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var input dto.PigSaleUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	tx := config.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// โหลด PigSale พร้อม Items เดิม
	pigSale := &models.PigSale{}
	if err := tx.Preload("Items").First(pigSale, "id = ?", pigSaleID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "PigSale not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	updates := make(map[string]interface{})

	// อัปเดต Date
	if input.Date != nil {
		parsedDate, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format"})
		}
		if parsedDate.After(time.Now()) {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Date cannot be in the future"})
		}
		updates["date"] = parsedDate
	}

	if len(*input.PigIDs) > 0 {
		//  คืนสถานะหมูเก่า
		var oldPigIDs []uint
		for _, item := range pigSale.Items {
			oldPigIDs = append(oldPigIDs, item.PigID)
		}
		if len(oldPigIDs) > 0 {
			if err := tx.Model(&models.Pig{}).Where("id IN ?", oldPigIDs).
				Update("status", "พร้อมขาย").Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to reset old pig status"})
			}
		}

		//  Soft delete ของ PigSaleItems เก่า
		if err := tx.Where("pig_sale_id = ?", pigSale.ID).Delete(&models.PigSaleItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete old PigSaleItems"})
		}

		//  ดึงหมูใหม่ที่ต้องขาย
		var pigs []models.Pig
		if err := tx.Where("id IN ?", *input.PigIDs).Find(&pigs).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
		if len(pigs) != len(*input.PigIDs) {
			tx.Rollback()
			return c.Status(404).JSON(fiber.Map{"error": "Some pigs not found"})
		}

		//  เช็คสถานะพร้อมขาย
		var readyPigs []models.Pig
		for _, pig := range pigs {
			if pig.Status == "พร้อมขาย" {
				readyPigs = append(readyPigs, pig)
			}
		}
		if len(readyPigs) == 0 {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "No pigs are ready for sale"})
		}

		var newItems []models.PigSaleItem
		for _, pig := range readyPigs {
			newItems = append(newItems, models.PigSaleItem{
				PigSaleID: pigSale.ID,
				PigID:     pig.ID,
			})
		}
		if err := tx.Create(&newItems).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create new PigSaleItems"})
		}

		var readyIDs []uint
		for _, pig := range readyPigs {
			readyIDs = append(readyIDs, pig.ID)
		}
		if err := tx.Model(&models.Pig{}).Where("id IN ?", readyIDs).
			Update("status", "ขายเเล้ว").Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update pig status"})
		}

		updates["amount"] = len(readyIDs)
	}

	if input.TotalPrice != nil {
		updates["total_price"] = *input.TotalPrice
	}
	if input.Buyer != nil {
		updates["buyer"] = *input.Buyer
	}
	if input.Note != nil {
		updates["note"] = *input.Note
	}

	// อัปเดต PigSale
	if err := tx.Model(&pigSale).Updates(updates).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update pig sale"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Transaction commit failed"})
	}

	// โหลดข้อมูลใหม่พร้อม Pig
	if err := config.DB.Preload("Items.Pig").First(pigSale, "id = ?", pigSaleID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to reload pig sale"})
	}

	// เตรียม response
	var pigIDs []uint
	var pigCodeNames []string
	for _, item := range pigSale.Items {
		pigIDs = append(pigIDs, item.PigID)
		pigCodeNames = append(pigCodeNames, item.Pig.CodeName)
	}

	resp := dto.PigSaleResponse{
		ID:          pigSale.ID,
		SaleCode:    pigSale.SaleCode,
		PigIDs:      pigIDs,
		PigCodeName: pigCodeNames,
		Date:        pigSale.Date,
		Amount:      len(pigIDs),
		TotalPrice:  pigSale.TotalPrice,
		Buyer:       pigSale.Buyer,
		Note:        pigSale.Note,
	}

	return c.JSON(resp)
}

func DeletePigSale(c *fiber.Ctx) error {
	id := c.Params("id")
	pigSaleID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// ✅ 1. เริ่มต้น Transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction"})
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// ✅ 2. ค้นหา PigSale พร้อมกับรายการหมู (Items) ที่เกี่ยวข้อง
	pigSale := &models.PigSale{}
	if err := tx.Preload("Items").First(pigSale, "id = ?", pigSaleID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "PigSale not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// ✅ 3. รวบรวม ID ของหมูทั้งหมดที่อยู่ในการขายครั้งนี้
	var pigIDsToReset []uint
	if len(pigSale.Items) > 0 {
		for _, item := range pigSale.Items {
			pigIDsToReset = append(pigIDsToReset, item.PigID)
		}
	}

	// ✅ 4. คืนสถานะหมูทั้งหมดกลับไปเป็น "พร้อมขาย"
	// ทำขั้นตอนนี้ก่อน เพื่อให้แน่ใจว่าถ้าเกิดปัญหากับการลบ เรายังไม่ได้ลบอะไรไป
	if len(pigIDsToReset) > 0 {
		if err := tx.Model(&models.Pig{}).Where("id IN ?", pigIDsToReset).Update("status", "พร้อมขาย").Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to reset pig status"})
		}
	}

	// ✅ 5. ลบรายการลูก (PigSaleItems) ทั้งหมดที่เกี่ยวข้องกับการขายนี้
	// GORM ที่มีการตั้งค่า association on delete: cascade อาจจะจัดการส่วนนี้ให้เอง
	// แต่เพื่อความแน่นอน เราสามารถลบเองได้
	if len(pigSale.Items) > 0 {
		if err := tx.Where("pig_sale_id = ?", pigSaleID).Delete(&models.PigSaleItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to delete PigSale items"})
		}
	}

	// ✅ 6. ลบบันทึกการขายหลัก (PigSale)
	// ใช้ Unscoped() เพื่อลบแบบถาวร (Hard Delete) หากตาราง PigSale มี gorm.DeletedAt
	if err := tx.Unscoped().Delete(pigSale).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete PigSale"})
	}

	// ✅ 7. ยืนยันการเปลี่ยนแปลงทั้งหมด
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Transaction commit failed"})
	}

	return c.JSON(fiber.Map{
		"message": "PigSale deleted successfully and pig statuses have been reset",
		"id":      pigSaleID,
	})
}

// func DeletePigSale(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	pigSaleID, err := strconv.Atoi(id)
// 	if err != nil {
// 		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
// 	}

// 	// tx := config.DB.Begin()
// 	// defer func() {
// 	// 	if r := recover(); r != nil {
// 	// 		tx.Rollback()
// 	// 	}
// 	// }()
// 	// if tx.Error != nil {
// 	// 	return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction"})
// 	// }

// 	pigSale := &models.PigSale{}
// 	if err := config.DB.First(pigSale, "id = ?", pigSaleID).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return c.Status(404).JSON(fiber.Map{"error": "PigSale not found"})
// 		}
// 		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
// 	}

// 	// ลบ PigSaleItem ทั้งหมดก่อน
// 	// if len(pigSale.Items) > 0 {
// 	// 	if err := tx.Where("pig_sale_id = ?", pigSaleID).Delete(&models.PigSaleItem{}).Error; err != nil {
// 	// 		tx.Rollback()
// 	// 		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete PigSale items"})
// 	// 	}
// 	// }

// 	// ลบ PigSale
// 	if err := config.DB.Unscoped().Delete(pigSale).Error; err != nil {
// 		// tx.Rollback()
// 		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete PigSale"})
// 	}

// 	// if err := tx.Commit().Error; err != nil {
// 	// 	return c.Status(500).JSON(fiber.Map{"error": "Transaction commit failed"})
// 	// }

// 	return c.JSON(fiber.Map{
// 		"message": "PigSale deleted successfully",
// 		"id":      pigSaleID,
// 	})
// }
