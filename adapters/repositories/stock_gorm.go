package repositories

import (
	"backend/dto"
	"backend/models"
	"backend/usecases"

	"gorm.io/gorm"
)

type StockGormRepo struct {
	db *gorm.DB
}

func NewStockGormRepo(db *gorm.DB) usecases.StockRepository {
	return &StockGormRepo{db: db}

}
func (r *StockGormRepo) CheckDuplicate(foodTypeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.FoodStock{}).Where("food_type_id = ? ", foodTypeID).Count(&count).Error
	return count > 0, err
}

func (r *StockGormRepo) Create(stock *models.FoodStock) error {
	return r.db.Create(stock).Error
}

func (r *StockGormRepo) GetByID(id uint) (*models.FoodStock, error) {
	var stock models.FoodStock
	err := r.db.Preload("FoodType").Preload("Creator").Preload("Updater").First(&stock, id).Error
	return &stock, err
}

func (r *StockGormRepo) Update(stock *models.FoodStock) error {
	return r.db.Save(stock).Error
}

func (r *StockGormRepo) GetAllPagi(input dto.ParamFoodStock) ([]models.FoodStock, int64, error) {
	var stock []models.FoodStock
	var total int64
	db := r.db.Model(&models.FoodStock{})
	db = db.Joins("LEFT JOIN food_types ON food_stocks.food_type_id = food_types.id")
	if input.Search != "" {
		keyword := "%" + input.Search + "%"
		db = db.Where("food_types.name ILIKE ? OR food_types.type OR CAST(food_stocks.amount AS TEXT) ILIKE ?", keyword, keyword, keyword)

	}
	// นับจํานวนก่อนตัด
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offSet := (input.Page - 1) * input.Limit
	db = db.Offset(offSet).Limit(input.Limit)

	err := db.Preload("Creator").Preload("Updater").Preload("FoodType").Order("food_stocks.id").Find(&stock).Error

	return stock, total, err
}

func (r *StockGormRepo) Delete(id uint) error {
	stock := &models.FoodStock{}

	return r.db.Unscoped().Delete(stock, id).Error
}

// เช็คเพื่่อ เเจ้ง user ว่ามีมีการใช้ stock นี้อยู่
func (r *StockGormRepo) IsUsedInFeeding(stockID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Feeding{}).Where("food_id = ?", stockID).Count(&count).Error
	return count > 0, err
}
