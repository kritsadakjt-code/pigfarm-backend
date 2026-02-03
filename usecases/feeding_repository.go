package usecases

import (
	"backend/dto"
	"backend/models"
)

type FeedingRepository interface {
	Create(feeding *models.Feeding, pigIDs []uint) ([]uint, error)
	GetById(id uint) (*models.Feeding, error)
	GetAll(param dto.ParamFeeding) ([]models.Feeding, int64, error)
	Delete(id uint) error
	Update(id uint, newFeeding *models.Feeding, newPigIDs []uint) error
}
