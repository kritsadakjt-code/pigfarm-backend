package repositories

import (
	"backend/usecases"

	"gorm.io/gorm"
)

type PigSaleGormRepo struct {
	db *gorm.DB
}

func NewPigSaleGormRepo(db *gorm.DB) usecases.PigSaleRepository {
	return &PigSaleGormRepo{db: db}
}
