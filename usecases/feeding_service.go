package usecases

import (
	"backend/models"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FeedingService struct {
	DB *gorm.DB
}

func NewFeedingService(db *gorm.DB) *FeedingService {
	return &FeedingService{DB: db}
}

// ItemToFeed คือ struct สำหรับรายการอาหารที่จะป้อน
type ItemToFeed struct {
	FoodID uint
	Amount float64
}

// CreateFeedingLogsForItems คือ Logic หลักที่รับ "หลายรายการ" และทำงานใน Transaction เดียว
func (s *FeedingService) CreateFeedingLogsForItems(items []ItemToFeed, dateTime time.Time, note string, createdBy uint) error {
	return s.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			// 1. ค้นหา Food Stock และ Lock แถวข้อมูลไว้เพื่อป้องกัน Race Condition
			var food models.FoodStock
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&food, item.FoodID).Error; err != nil {
				return fmt.Errorf("food stock ID %d not found", item.FoodID)
			}

			// 2. เช็คว่ามีของพอหรือไม่
			if food.Amount < item.Amount {
				return fmt.Errorf("not enough stock for %s (required: %.2f, available: %.2f)", food.Name, item.Amount, food.Amount)
			}

			// 3. อัปเดต (ลด) Food Stock
			food.Amount -= item.Amount
			if err := tx.Save(&food).Error; err != nil {
				return err
			}

			// 4. สร้าง Feeding record ใหม่สำหรับรายการนี้
			feedingLog := models.Feeding{
				FoodID:    item.FoodID,
				Amount:    item.Amount,
				DateTime:  dateTime,
				Note:      note,
				CreatedBy: createdBy,
				UpdatedBy: createdBy, // ตอนสร้างครั้งแรกให้เป็นคนเดียวกัน
			}
			if err := tx.Create(&feedingLog).Error; err != nil {
				return err
			}
		}

		// ถ้าทุกอย่างผ่าน จะ Commit Transaction ให้อัตโนมัติ
		return nil
	})
}
