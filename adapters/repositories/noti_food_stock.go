package repositories

import (
	"backend/models"
	"backend/usecases"
	"fmt"

	"gorm.io/gorm"
)

type GormFoodStockRepository struct {
	db *gorm.DB
}

func NewGormFoodStockRepository(db *gorm.DB) usecases.FoodStockRepository {
	return &GormFoodStockRepository{db: db}
}

func (r *GormFoodStockRepository) GetLowStock(quantity float64) ([]models.FoodStock, error) {
	var foods []models.FoodStock
	result := r.db.Where("amount <= ? AND amount > 0", quantity).Find(&foods)
	if result.Error != nil {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return foods, nil
}
