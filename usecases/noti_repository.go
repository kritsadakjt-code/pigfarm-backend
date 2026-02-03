package usecases

import (
	"backend/entities"
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
	GetLowStock(quantity float64) ([]entities.FoodStock, error)
}

type BreedingRepository interface {
	GetUpcomingBirths(startDate, endDate time.Time) ([]entities.Breeding, error)
}
