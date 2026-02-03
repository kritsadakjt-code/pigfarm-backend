package repositories

import (
	"backend/models"
	"backend/usecases"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FeedingScheduleGormRepo struct {
	db *gorm.DB
}

func NewFeedingScheduleGormRepo(db *gorm.DB) usecases.FeedingRepositoryScheduler {
	return &FeedingScheduleGormRepo{db: db}
}

func (r *FeedingScheduleGormRepo) GetSchedulesByTime(timeStr string) ([]models.FeedingSchedule, error) {
	var schedules []models.FeedingSchedule
	err := r.db.Preload("Items").
		Where("is_active = ? AND scheduled_time = ?", true, timeStr).
		Find(&schedules).Error
	if err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *FeedingScheduleGormRepo) DeductStockAndLogFeeding(items []usecases.ItemToFeed, dateTime time.Time, note string, createdBy uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			// 1. ค้นหา Food Stock และ Lock แถวข้อมูลไว้เพื่อป้องกัน Race Condition
			var food models.FoodStock
			if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&food, item.FoodID).Error; err != nil {
				return fmt.Errorf("food stock ID %d not found", item.FoodID)
			}

			// 2. เช็คว่ามีของพอหรือไม่
			if food.Amount < item.Amount {
				return fmt.Errorf("not enough stock for %s (required: %.2f, available: %.2f)", food.FoodType.Name, item.Amount, food.Amount)
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
