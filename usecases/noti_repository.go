package usecases

import (
	"backend/entities"
	"backend/models"
	"time"
)

type NotificationRepository interface {
	Create(notification *entities.Notification) error
	GetAll() ([]entities.Notification, error)
	GetByID(id uint) (*entities.Notification, error)
	GetUnreadCount() (int64, error)
	MarkAsRead(id uint) error
	MarkAllAsRead() error
	Delete(id uint) error
	ExistsToday(notiType, keyword string) (bool, error)
}

type FoodStockRepository interface {
	GetLowStock(quantity float64) ([]models.FoodStock, error)
}

type BreedingRepository interface {
	GetUpcomingBirths(startDate, endDate time.Time) ([]models.Breeding, error)
}
