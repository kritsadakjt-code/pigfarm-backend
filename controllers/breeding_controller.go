package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateBreeding(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	var input dto.BreedingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	// check format input
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// check father mother in database
	father := &models.Pig{}
	if err := config.DB.First(father, "id = ?", input.FatherID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "father pig not found"})
	}
	mother := &models.Pig{}
	if err := config.DB.First(mother, "id = ?", input.MotherID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "mother pig not found"})
	}

	// พ่อเเม่ต้องไม่ใช่ตัวเดียวกัน
	if input.FatherID == input.MotherID {
		return c.Status(400).JSON(fiber.Map{"error": "Father and mother cannot be the same pig"})
	}
	// check gender type
	if father.Gender != "ผู้" || father.Type != "พ่อพันธุ์" {
		return c.Status(400).JSON(fiber.Map{"error": "Father pig must be male breeder"})
	}
	if mother.Gender != "เมีย" || mother.Type != "เเม่พันธุ์" {
		return c.Status(400).JSON(fiber.Map{"error": "Mother pig must be female breeder"})
	}
	// check ห้ามผสมซํ้าในวันเดียวกัน
	var count int64
	err = config.DB.Model(&models.Breeding{}).Where("father_id = ? AND mother_id = ? AND breeding_date = ?", input.FatherID, input.MotherID, input.BreedingDate).Count(&count).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to check breeding"})
	}
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "breeding record already exists"})
	}
	// status ต้องพร้อมผสม
	if father.Status != "พร้อมผสม" {
		return c.Status(400).JSON(fiber.Map{"error": "Father pig is not ready for breeding"})
	}
	if mother.Status != "พร้อมผสม" {
		return c.Status(400).JSON(fiber.Map{"error": "Mother pig is not ready for breeding"})
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")
	parsedDate, err := time.ParseInLocation("2006-01-02", input.BreedingDate, loc)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Birth date cannot be in the future"})
	}
	// calculate expected birth
	expectedBirth := parsedDate.AddDate(0, 0, 114)
	breeding := models.Breeding{
		FatherID:      input.FatherID,
		MotherID:      input.MotherID,
		BreedingDate:  parsedDate,
		ExpectedBirth: expectedBirth,
		Status:        "รอผล",
		Result:        "รอผล",
		Note:          input.Note,
		CreatedBy:     user_id,
		UpdatedBy:     user_id,
	}
	// transaction
	tx := config.DB.Begin()
	// check error กรณีที่ไม่ได้เกิดจาก create save
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Create(&breeding).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "cannot create breeding"})
	}

	// update status mother
	mother.Status = "อุ้มท้อง"
	if err := tx.Save(mother).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update mother status"})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	if err := config.DB.Preload("Creator").Preload("Updater").First(&breeding, "id = ?", breeding.ID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch preload breeding"})
	}
	resp := dto.BreedingResponse{
		ID:             breeding.ID,
		FatherID:       breeding.FatherID,
		MotherID:       breeding.MotherID,
		FatherCodename: father.CodeName,
		MotherCodename: mother.CodeName,
		BreedingDate:   breeding.BreedingDate,
		ExpectedBirth:  breeding.ExpectedBirth,
		Status:         breeding.Status,
		Result:         breeding.Result,
		Note:           breeding.Note,
		CreatedName:    breeding.Creator.FullName,
		UpdatedName:    breeding.Updater.FullName,
	}
	return c.Status(201).JSON(resp)
}

func UpdateBreeding(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}
	user_id := uint(user_id64)

	id := c.Params("id")

	var input dto.BreedingUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// ✅ 1. เริ่มต้น Transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// ✅ 2. ค้นหา Breeding record พร้อมกับข้อมูลแม่หมู (Mother) ภายใน Transaction
	var breeding models.Breeding
	if err := tx.Preload("Mother").First(&breeding, "id = ?", id).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Not found breeding ID"})
	}

	updates := make(map[string]interface{})
	if input.BreedingDate != nil {
		loc, _ := time.LoadLocation("Asia/Bangkok")
		parsedDate, err := time.ParseInLocation("2006-01-02", *input.BreedingDate, loc)
		if err != nil {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}
		if parsedDate.After(time.Now()) {
			tx.Rollback()
			return c.Status(400).JSON(fiber.Map{"error": "Breeding date cannot be in the future"})
		}
		updates["breeding_date"] = parsedDate
		updates["expected_birth"] = parsedDate.AddDate(0, 0, 114)
	}

	// ✅ 3. ตรรกะหลัก: อัปเดตสถานะ Breeding และสถานะแม่หมู
	if input.Status != nil {
		updates["status"] = *input.Status

		// กำหนดสถานะของแม่หมูใหม่
		var newMotherStatus string

		switch *input.Status {
		case "อุ้มท้อง":
			updates["result"] = "รอผล"
			newMotherStatus = "อุ้มท้อง"
		case "ผสมไม่ติด":
			updates["result"] = "ไม่สําเร็จ"
			newMotherStatus = "พร้อมผสม" // คืนสถานะให้พร้อมผสมใหม่
		case "เเท้ง":
			updates["result"] = "ไม่สําเร็จ"
			newMotherStatus = "พักท้อง" // ต้องพักฟื้นก่อนผสมใหม่
		case "คลอดเเล้ว":
			updates["result"] = "สําเร็จ"
			newMotherStatus = "ให้นมลูก"
		}

		// อัปเดตสถานะของแม่หมู ถ้ามีการเปลี่ยนแปลง
		if newMotherStatus != "" && breeding.Mother.Status != newMotherStatus {
			if err := tx.Model(&breeding.Mother).Update("status", newMotherStatus).Error; err != nil {
				tx.Rollback()
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update mother's status"})
			}
		}
	}

	if input.Note != nil {
		updates["note"] = *input.Note
	}
	updates["updated_by"] = user_id

	// ✅ 4. อัปเดต Breeding record
	if err := tx.Model(&breeding).Updates(updates).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Fail to update breeding"})
	}

	// ✅ 5. ยืนยันการเปลี่ยนแปลงทั้งหมด
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	// โหลดข้อมูลทั้งหมดอีกครั้งเพื่อส่งกลับ (รวม Father ที่อาจจะยังไม่ได้โหลด)
	if err := config.DB.Preload("Father").Preload("Mother").Preload("Updater").First(&breeding, "id = ?", id).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to reload breeding data"})
	}

	resp := dto.BreedingResponse{
		ID:             breeding.ID,
		FatherID:       breeding.FatherID,
		MotherID:       breeding.MotherID,
		FatherCodename: breeding.Father.CodeName,
		MotherCodename: breeding.Mother.CodeName,
		BreedingDate:   breeding.BreedingDate,
		ExpectedBirth:  breeding.ExpectedBirth,
		Status:         breeding.Status,
		Result:         breeding.Result,
		Note:           breeding.Note,
		UpdatedName:    breeding.Updater.FullName,
	}

	return c.JSON(resp)
}

func SearchBreeding(c *fiber.Ctx) error {
	keyword := strings.TrimSpace(c.Query("keyword"))

	var breedings []models.Breeding
	db := config.DB.Model(&models.Breeding{}).Preload("Father").Preload("Mother").Preload("Creator").Preload("Updater").
		Joins("LEFT JOIN pigs father on father.id = breedings.father_id"). // join เพราะจะเอาข้อมูล filter ต่อ
		Joins("LEFT JOIN pigs mother on mother.id = breedings.mother_id").
		Joins("LEFT JOIN users creator on creator.id = breedings.created_by")

	if keyword != "" {
		kw := "%" + keyword + "%"
		if err := db.Where(
			"breedings.status ILIKE ? OR breedings.result ILIKE ? OR breedings.note ILIKE ? OR father.code_name ILIKE ? OR mother.code_name ILIKE ? OR creator.role ILIKE ?",
			kw, kw, kw, kw, kw, kw,
		).Find(&breedings).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "search breedings failed"})
		}
	} else {
		if err := db.Find(&breedings).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "cannot fetch breedings"})
		}

	}

	var resp []dto.BreedingResponse
	for _, b := range breedings {
		resp = append(resp, dto.BreedingResponse{
			ID:             b.ID,
			FatherID:       b.FatherID,
			MotherID:       b.MotherID,
			FatherCodename: b.Father.CodeName,
			MotherCodename: b.Mother.CodeName,
			BreedingDate:   b.BreedingDate,
			ExpectedBirth:  b.ExpectedBirth,
			Status:         b.Status,
			Result:         b.Result,
			Note:           b.Note,
			CreatedName:    b.Creator.FullName,
			CreatedRole:    b.Creator.Role,
			UpdatedName:    b.Updater.FullName,
			UpdatedRole:    b.Updater.Role,
		})
	}

	return c.JSON(resp)
}

func GetAllBreeding(c *fiber.Ctx) error {
	var breedings []models.Breeding
	if err := config.DB.Order("breeding_date DESC").Model(&models.Breeding{}).Preload("Father").Preload("Mother").Preload("Creator").Preload("Updater").Find(&breedings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch breeding data"})
	}

	var resp []dto.BreedingResponse
	for _, b := range breedings {
		resp = append(resp, dto.BreedingResponse{
			ID:             b.ID,
			FatherID:       b.FatherID,
			MotherID:       b.MotherID,
			FatherCodename: b.Father.CodeName,
			MotherCodename: b.Mother.CodeName,
			BreedingDate:   b.BreedingDate,
			ExpectedBirth:  b.ExpectedBirth,
			Status:         b.Status,
			Result:         b.Result,
			Note:           b.Note,
			CreatedName:    b.Creator.FullName,
			CreatedRole:    b.Creator.Role,
			UpdatedName:    b.Updater.FullName,
			UpdatedRole:    b.Updater.Role,
		})
	}
	return c.JSON(resp)
}

func GetBreedingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	breeding_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var breeding models.Breeding
	if err := config.DB.Model(&models.Breeding{}).Preload("Father").Preload("Mother").Preload("Creator").Preload("Updater").First(&breeding, "id = ?", breeding_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found breeding ID"})
	}

	resp := dto.BreedingResponse{
		ID:             breeding.ID,
		FatherID:       breeding.FatherID,
		MotherID:       breeding.MotherID,
		FatherCodename: breeding.Father.CodeName,
		MotherCodename: breeding.Mother.CodeName,
		BreedingDate:   breeding.BreedingDate,
		ExpectedBirth:  breeding.ExpectedBirth,
		Status:         breeding.Status,
		Result:         breeding.Result,
		Note:           breeding.Note,
		CreatedName:    breeding.Creator.FullName,
		CreatedRole:    breeding.Creator.Role,
		UpdatedName:    breeding.Updater.FullName,
		UpdatedRole:    breeding.Updater.Role,
	}

	return c.JSON(resp)

}

func DeleteBreeding(c *fiber.Ctx) error {
	id := c.Params("id")
	breeding_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var breeding models.Breeding
	if err := config.DB.Preload("Mother").Preload("Father").First(&breeding, "id = ?", breeding_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not found breeding ID"})
	}

	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// คืนสถานะแม่ถ้าเป็นอุ้มท้อง
	if breeding.Mother.Status == "อุ้มท้อง" {
		breeding.Mother.Status = "พร้อมผสม"
		if err := tx.Save(&breeding.Mother).Error; err != nil {
			tx.Rollback()
			return c.Status(500).JSON(fiber.Map{"error": "Failed to update mother status"})
		}
	}

	// คืนสถานะพ่อถ้าต้องการ (เช่น พ่ออาจถูก mark ว่าไม่พร้อมผสม)
	// if breeding.Father.Status != "พร้อมผสม" {
	// 	breeding.Father.Status = "พร้อมผสม"
	// 	if err := tx.Save(&breeding.Father).Error; err != nil {
	// 		tx.Rollback()
	// 		return c.Status(500).JSON(fiber.Map{"error": "Failed to update father status"})
	// 	}
	// }

	// ลบ record breeding
	if err := tx.Delete(&breeding).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete breeding"})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.JSON(fiber.Map{"message": "Breeding deleted and mother status restored"})
}
