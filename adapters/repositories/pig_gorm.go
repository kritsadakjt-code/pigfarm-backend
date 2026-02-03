package repositories

import (
	"backend/dto"
	"backend/models"
	"backend/usecases"

	"gorm.io/gorm"
)

type PigGormRepository struct {
	db *gorm.DB
}

func NewPigGormRepository(db *gorm.DB) usecases.PigRepository {
	return &PigGormRepository{db: db}
}

func (r *PigGormRepository) GenerateNextCode(pattern string) (*models.Pig, error) {
	var pig models.Pig
	if err := r.db.First(&pig, "code_name LIKE ?", pattern).Error; err != nil {
		return nil, err
	}
	return &pig, nil
}

func (r *PigGormRepository) Create(pig *models.Pig) error {
	return r.db.Create(pig).Error
}

func (r *PigGormRepository) GetByID(id uint) (*models.Pig, error) {
	var pig models.Pig
	if err := r.db.Preload("Creator").Preload("Updater").First(&pig, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &pig, nil
}

func (r *PigGormRepository) Update(id uint, update map[string]interface{}) error {
	result := r.db.Model(&models.Pig{}).Where("id = ?", id).Updates(update)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// func (r *PigGormRepository) FindAllPagination(param dto.PigParam) ([]models.Pig, int64, error) {
// 	var pig []models.Pig
// 	var total int64
// 	db := r.db.Model(&models.Pig{})

// 	// ค้นหา เช่น codeName Name
// 	if param.Search != "" {
// 		keyword := "%" + param.Search + "%"
// 		db = db.Where("code_name ILIKE ? OR name ILIKE ?", keyword, keyword)
// 	}

// 	// ค้นหาเจาะจง
// 	if param.Status != "" {
// 		db = db.Where("status = ?", param.Status)
// 	}
// 	// นับจํานวนข้อมูลทั้งหมดเพื่อเอามาคํานวณ
// 	if err := db.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}

// 	// คํานวณหน้าที่ต้องข้าม
// 	offset := (param.Page - 1) * param.Limit

// 	if err := db.Preload("Creator").Preload("Updater").
// 		Offset(offset).Limit(param.Limit).Order("id DESC").
// 		Find(&pig).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	return pig, total, nil
// }

func (r *PigGormRepository) FindAllPagination(param dto.PigParam) ([]models.Pig, int64, error) {
	// var pig []models.Pig
	// var total int64

	// db := r.db.Model(&models.Pig{})
	// if param.Search != "" {
	// 	keyword := "%" + param.Search + "%"
	// 	db = db.Where("code_name ILIKE ? OR name LIKE ?", keyword, keyword)
	// }
	// if param.Status != "" {
	// 	db = db.Where("status", param.Status)
	// }
	// if err := db.Count(&total).Error; err != nil {
	// 	return nil, 0, err
	// }

	// // หาหน้าที่ต้องข้าม
	// offset := (param.Page - 1) * param.Limit
	// if err := db.Preload("Creator").Preload("Updater").Offset(offset).Limit(param.Limit).Order("id DESC").Find(&pig).Error; err != nil {
	// 	return nil, 0, err
	// }
	// return pig, total, nil

	var pig []models.Pig
	var total int64

	db := r.db.Model(&models.Pig{})
	if param.Search != "" {
		keyword := "%" + param.Search + "%"
		db = db.Where("code_name ILIKE ? OR name ILIKE ?", keyword, keyword)
	}
	if param.Status != "" {
		db = db.Where("status ILIKE ?", param.Status)
	}
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (param.Page - 1) * param.Limit
	if err := db.Preload("Creator").Preload("Updater").Offset(offset).Limit(param.Limit).Order("id DESC").Find(&pig).Error; err != nil {
		return nil, 0, err
	}
	return pig, total, nil

}

func (r *PigGormRepository) Delete(id uint) error {
	if err := r.db.Model(&models.Pig{}).Unscoped().Delete("id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (r *PigGormRepository) IsUsedInBreeding(pigID uint) (bool, error) {
	// var count int64
	// err := r.db.Model(&models.Breeding{}).Where("father_id = ? OR mother_id = ?", pigID, pigID).Count(&count).Error
	// return count > 0, err

	var count int64
	err := r.db.Model(&models.Breeding{}).Where("father_id = ? OR mother_id = ?", pigID, pigID).Count(&count).Error
	return count > 0, err

}
