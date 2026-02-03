package repositories

import (
	"backend/dto"
	"backend/models"
	"backend/usecases"
	"time"

	"gorm.io/gorm"
)

type BreedingGormRepo struct {
	db *gorm.DB
}

func NewBreedingGormRepo(db *gorm.DB) usecases.BreedingRepo {
	return &BreedingGormRepo{db: db}
}

func (r *BreedingGormRepo) Create(breeding *models.Breeding) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(breeding).Error; err != nil {
			return err
		}
		if err := tx.Model(&models.Pig{}).Where("id = ?", breeding.MotherID).Update("status", "อุ้มท้อง").Error; err != nil {
			return err
		}
		return nil

	})
}

func (r *BreedingGormRepo) CheckBreedingAlready(fatherID, motherID uint, date time.Time) (bool, error) {
	var count int64
	err := r.db.Model(&models.Breeding{}).Where("father_id = ? AND mother_id = ? AND DATE(breeding_date) = DATE(?)", fatherID, motherID, date).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *BreedingGormRepo) GetByID(id uint) (*models.Breeding, error) {
	var breeding models.Breeding
	err := r.db.Preload("Mother").Preload("Father").Preload("Creator").Preload("Updater").First(&breeding, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &breeding, nil
}

func (r *BreedingGormRepo) UpdateBreeding(id uint, updates map[string]interface{}, motherID uint, newStatusMother string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.Breeding{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}

		if newStatusMother != "" && motherID != 0 {
			if err := tx.Model(&models.Pig{}).Where("id = ?", motherID).Update("status", newStatusMother).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *BreedingGormRepo) GetAll(param dto.BreedingParam) ([]models.Breeding, int64, error) {
	var breeding []models.Breeding
	var total int64
	db := r.db.Model(&models.Breeding{}) // ค่อยมา preload Creator Updaterถ้าไม่ได้
	if param.Search != "" {
		keyword := "%" + param.Search + "%"
		db = db.Joins("LEFT JOIN users AS creator ON breeding.created_by = creator.id").
			Joins("LEFT JOIN users AS updater ON breeding.updated_by = updater.id").
			Where("status ILIKE ? OR creator.full_name ILIKE ? OR updater.full_name ILIKE", keyword, keyword, keyword)
	}
	if param.Status != "" {
		db = db.Where("status = ?", param.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offSet := (param.Page - 1) * param.Limit
	err := db.Preload("Updater").Preload("Creator").Offset(offSet).Limit(param.Limit).Order("id DESC").Find(&breeding).Error
	if err != nil {
		return nil, 0, err
	}
	return breeding, total, nil
}

func (r *BreedingGormRepo) Delete(id uint, motherID uint, reStatus bool) error {

	return r.db.Transaction(func(tx *gorm.DB) error {
		if reStatus && motherID != 0 {
			if err := tx.Model(&models.Pig{}).Where("id = ?", motherID).Update("status", "พร้อมผสม").Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&models.Breeding{}).Unscoped().Delete("id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}
