package usecases

import (
	"backend/dto"
	"backend/models"
)

type PigRepository interface {
	Create(*models.Pig) error
	GenerateNextCode(pattern string) (*models.Pig, error)
	GetByID(id uint) (*models.Pig, error)
	Update(id uint, update map[string]interface{}) error
	FindAllPagination(param dto.PigParam) ([]models.Pig, int64, error)
	Delete(id uint) error
	IsUsedInBreeding(pigID uint) (bool, error)
}
