package repositories

import (
	"backend/dto"
	"backend/models"
	"backend/usecases"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FeedingGormRepository struct {
	db *gorm.DB
}

func NewFeedingGormRepository(db *gorm.DB) usecases.FeedingRepository {
	return &FeedingGormRepository{db: db}
}

func (r *FeedingGormRepository) Create(feeding *models.Feeding, pigIDs []uint) ([]uint, error) {
	var validPigIDs []uint
	err := r.db.Transaction(func(tx *gorm.DB) error {
		// lock เเย่งกันตัดสต็อก
		var food models.FoodStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&food, "id = ?", feeding.FoodID).Error; err != nil {
			return err
		}
		if food.Amount < feeding.Amount {
			return usecases.ErrNotEnoughFood
		}
		// ตัดสต็อก
		food.Amount -= feeding.Amount
		if err := tx.Save(&food).Error; err != nil {
			return err
		}

		// ค้นหาหมูที่ ไม่ตาย ไม่ขาย จาก pigIDs ที่ส่งมา
		var validPigs []models.Pig
		if len(pigIDs) > 0 {
			if err := tx.Where("id IN ? AND status NOT IN ?", pigIDs, []string{"ตายเเล้ว", "ขายเเล้ว"}).Find(&validPigs).Error; err != nil {
				return err
			}
		}
		if len(validPigs) == 0 {
			return usecases.ErrNoValidPigs
		}
		if err := tx.Create(feeding).Error; err != nil {
			return err
		}

		// สร้าง items
		var items []models.FeedingItem
		for _, p := range validPigs {
			items = append(items, models.FeedingItem{
				FeedingID: feeding.ID,
				PigID:     p.ID,
			})
			validPigIDs = append(validPigIDs, p.ID)
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil

	})
	// for _, v := range validPigIDs {
	// 	fmt.Println(v)
	// }
	// fmt.Println(validPigIDs)
	return validPigIDs, err

}

func (r *FeedingGormRepository) Update(id uint, newFeeding *models.Feeding, newPigIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// ดึงข้อมูล feeding ตัวเก่า เพื่อเอาค่าเดิมที่บันทึกไปกลับมาคืนสต็อก
		var oldFeeding models.Feeding
		if err := tx.First(&oldFeeding, id).Error; err != nil {
			return err
		}

		// คืน stock
		var oldStock models.FoodStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&oldStock, oldFeeding.FoodID).Error; err != nil {
			return err
		}

		oldStock.Amount += oldFeeding.Amount
		if err := tx.Save(&oldStock).Error; err != nil {
			return err
		}

		// ตัดสต็อกใหม่
		var newStock models.FoodStock
		if oldFeeding.FoodID == newFeeding.FoodID {
			newStock = oldStock // oleStock ถูกคืนเเล้ว
		} else {
			// กรณีเปลี่ยนชนิดอาหารใหม่ ไปดึง stock ตัวใหม่มา
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&newStock, newFeeding.FoodID).Error; err != nil {
				return err
			}

		}
		if newStock.Amount < newFeeding.Amount {
			return usecases.ErrNotEnoughFood
		}
		newStock.Amount -= newFeeding.Amount
		if err := tx.Save(&newStock).Error; err != nil {
			return err
		}
		// make sure ให้ id ถูกต้อง
		newFeeding.ID = id
		// update table
		if err := tx.Model(newFeeding).Updates(newFeeding).Error; err != nil {
			return err
		}

		// update items
		// ลบ item เดิมทิ้ง
		if err := tx.Unscoped().Where("feeding_id = ?", id).Delete(&models.FeedingItem{}).Error; err != nil {
			return err
		}

		// เริ่มสร้าง items ใหม่
		if len(newPigIDs) > 0 {
			var validPigs []models.Pig
			if err := tx.Where("id IN ? AND status NOT IN ?", newPigIDs, []string{"ตายเเล้ว", "ขายเเล้ว"}).Find(&validPigs).Error; err != nil {
				return err
			}
			// ถ้า user เลือกหมูมา เเต่หมูตาย ขายหมดเเล้ว
			if len(validPigs) == 0 {
				return usecases.ErrNoValidPigs
			}
			var items []models.FeedingItem
			for _, p := range validPigs {
				items = append(items, models.FeedingItem{
					FeedingID: id,
					PigID:     p.ID,
				})

			}
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil

	})
}
func (r *FeedingGormRepository) GetById(id uint) (*models.Feeding, error) {

	var feeding models.Feeding
	err := r.db.Preload("FoodStock").
		Preload("Creator").Preload("Updater").
		Preload("Items.Pig").First(&feeding, id).Error
	return &feeding, err
}

func (r *FeedingGormRepository) GetAll(param dto.ParamFeeding) ([]models.Feeding, int64, error) {
	var feedings []models.Feeding
	var total int64
	db := r.db.Model(&models.Feeding{})

	if param.FoodID != 0 {
		db = db.Where("feedings.food_id = ?", param.FoodID)
	}
	if param.Search != "" {
		keyword := "%" + param.Search + "%"
		db = db.Joins("LEFT JOIN food_stocks ON feedings.food_id = food_stocks.id").
			Joins("LEFT JOIN users AS creator ON feedings.created_by = creator.id").
			Joins("LEFT JOIN users AS updater ON feedings.updated_by = updater.id").
			Joins("LEFT JOIN feeding_items ON feedings.id = feeding_items.feeding_id").
			Joins("LEFT JOIN pigs ON feeding_items.pig_id = pigs.id").
			Where(`
			feedings.note ILIKE ? OR
			food_stocks.name ILIKE ? OR
			pigs.code_name ILIKE ? OR
			creator.full_name ILIKE ? OR 
			updater.full_name ILIKE ?
			`, keyword, keyword, keyword, keyword, keyword).
			Group("feedings.id") // รวมข้อมูลที่ซํ้ากันกเป็นชุดเดียว

	}

	// นับข้อมูลที่ feedings.id ซํ้ากัน เพื่อเอาจํานวนข้อมูลไปคํานวณ page
	if err := db.Session(&gorm.Session{}).Distinct("feedings.id").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (param.Page - 1) * param.Limit

	err := db.Preload("FoodStock").
		Preload("Creator").
		Preload("Updater").
		Preload("Items.Pig").
		Order("feedings.date_time DESC").
		Offset(offset).
		Limit(param.Limit).
		Find(&feedings).Error

	return feedings, total, err

}

func (r *FeedingGormRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var feeding models.Feeding
		if err := tx.First(&feeding, "id = ?", id).Error; err != nil {
			return err
		}
		// lock ป้องกันการเเย่งกัน
		var food models.FoodStock
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&food, "id = ?", feeding.FoodID).Error; err != nil {
			return err
		}
		food.Amount += feeding.Amount
		if err := tx.Save(&food).Error; err != nil {
			return err
		}

		// ลบรายการลูก
		if err := tx.Unscoped().Where("feeding_id = ?", id).Delete(&models.FeedingItem{}).Error; err != nil {
			return err
		}
		// ลบรายการเเม่
		if err := tx.Unscoped().Delete(&feeding).Error; err != nil {
			return err
		}
		return nil
	})
}
