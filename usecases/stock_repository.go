package usecases

import (
	"backend/dto"
	"backend/models"
)

type StockRepository interface {
	CheckDuplicate(foodTypeID uint) (bool, error)
	Create(stock *models.FoodStock) error
	GetByID(id uint) (*models.FoodStock, error)
	Update(stock *models.FoodStock) error
	GetAllPagi(input dto.ParamFoodStock) ([]models.FoodStock, int64, error)
	Delete(id uint) error
	IsUsedInFeeding(foodID uint) (bool, error)
}
