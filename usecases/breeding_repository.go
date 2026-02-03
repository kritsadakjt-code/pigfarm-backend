package usecases

import (
	"backend/dto"
	"backend/models"
	"time"
)

type BreedingRepo interface {
	Create(breeding *models.Breeding) error
	CheckBreedingAlready(fatherID, motherID uint, date time.Time) (bool, error)
	GetByID(id uint) (*models.Breeding, error)
	UpdateBreeding(id uint, updates map[string]interface{}, motherID uint, newStatusMother string) error
	GetAll(param dto.BreedingParam) ([]models.Breeding, int64, error)
	Delete(id uint, motherID uint, reStatus bool) error
}
