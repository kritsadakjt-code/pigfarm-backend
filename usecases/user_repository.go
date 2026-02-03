package usecases

import (
	"backend/dto"
	"backend/models"
	"time"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	Save(user *models.User) error
	FindByResetToken(token string, expiry time.Time) (*models.User, error)

	ExistEmail(email string) (bool, error)
	FindByEmailVerifyToken(token string, expiry time.Time) (*models.User, error)
	FindByEmailVerifyRegisterToken(token string) (*models.User, error)
	GetAll() ([]models.User, error)
	FindByID(id uint) (*models.User, error)
	Delete(id uint) error
	FindAllPagination(param dto.UserParam) ([]models.User, int64, error)
	UpdateStatus(id uint, status string) error
}
