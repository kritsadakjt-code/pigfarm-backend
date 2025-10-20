package repositories

import (
	"backend/models"
	"backend/usecases"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GormBreedingRepository struct {
	db *gorm.DB
}

func NewGormBreedingRepository(db *gorm.DB) usecases.BreedingRepository {
	return &GormBreedingRepository{db: db}
}

func (r *GormBreedingRepository) GetUpcomingBirths(startDate, endDate time.Time) ([]models.Breeding, error) {
	var breedings []models.Breeding
	result := r.db.Preload("Mother").Where("status = ? AND expected_birth BETWEEN ? AND ?", "อุ้มท้อง", startDate, endDate).Find(&breedings)
	if result.Error != nil {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	return breedings, nil
}
