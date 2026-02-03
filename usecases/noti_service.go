package usecases

import (
	"backend/entities"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	ErrNotificationNotFound = errors.New("notification not found")
	ErrInvalidNotification  = errors.New("invalid notification data")
)

type NotificationService struct {
	noti          NotificationRepository
	foodStockRepo FoodStockRepository
	breedingRepo  BreedingRepository
}

func NewNotificationService(
	noti NotificationRepository,
	foodStockRepo FoodStockRepository,
	breedingRepo BreedingRepository) *NotificationService {
	return &NotificationService{
		noti:          noti,
		foodStockRepo: foodStockRepo,
		breedingRepo:  breedingRepo,
	}
}

func (s *NotificationService) CreateNotification(notiType, title, message string) error {
	if notiType == "" || title == "" || message == "" {
		return ErrInvalidNotification
	}

	notification := entities.NewNotification(notiType, title, message)

	// ถ้าไม่ใช้ constructor ของ NewNotification ทําตรงๆ เเบบนี้ได้
	// notification2 := entities.Notification{
	// 	Type:    notiType,
	// 	Title:   title,
	// 	Message: message,
	// 	IsRead:  false,
	// }

	if err := s.noti.Create(notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	return nil

}

func (s *NotificationService) GetAllNotifications() ([]entities.Notification, error) {
	notifications, err := s.noti.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	return notifications, nil

}

func (s *NotificationService) GetNotificationByID(id uint) (*entities.Notification, error) {
	if id == 0 {
		return nil, ErrInvalidNotification
	}
	notification, err := s.noti.GetByID(id)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) GetUnreadCount() (int64, error) {
	count, err := s.noti.GetUnreadCount()
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}
	return count, nil
}

func (s *NotificationService) MarkAsRead(id uint) error {
	if id == 0 {
		return ErrInvalidNotification
	}

	if err := s.noti.MarkAsRead(id); err != nil {
		return fmt.Errorf("failed to mark as read: %w", err)
	}

	return nil
}

func (s *NotificationService) MarkAllAsRead() error {
	if err := s.noti.MarkAllAsRead(); err != nil {
		return fmt.Errorf("failed to mark all as read: %w", err)
	}
	return nil
}

func (s *NotificationService) DeleteNotification(id uint) error {
	if id == 0 {
		return ErrInvalidNotification
	}

	if err := s.noti.Delete(id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

func (s *NotificationService) CheckFoodStock() error {
	foods, err := s.foodStockRepo.GetLowStock(10)
	var title string
	if err != nil {
		return fmt.Errorf("failed to get low stock foods: %w", err)
	}
	for _, food := range foods {
		// เช็คว่ามี notification วันนี้ยัง
		exists, err := s.noti.ExistsToday("food_low", food.FoodTypeName)
		if err != nil {
			fmt.Printf("Error checking notification exists: %v\n", err)
			continue
		}
		if exists {
			continue
		}
		if food.IsCriticalStock() {
			title = "อาหารใกล้หมดอย่างเร่งด่วน"
		} else {
			title = "อาหารใกล้หมด"
		}

		message := fmt.Sprintf("%s เหลือเพียง %.2f Kg", food.FoodTypeName, food.Amount)
		if err := s.CreateNotification("food_low", title, message); err != nil {
			fmt.Printf("Error checking notification exists: %v\n", err)
			continue
		}
	}
	return nil
}

func (s *NotificationService) CheckBreeding() error {
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 0, 7)

	breedings, err := s.breedingRepo.GetUpcomingBirths(startDate, endDate)
	if err != nil {
		return fmt.Errorf("failed to get upcoming births: %w", err)
	}

	for _, breeding := range breedings {
		if !breeding.IsPregnant() {
			continue
		}
		exists, err := s.noti.ExistsToday(entities.TypeBirthDue, breeding.MotherCodename)
		if err != nil {
			return err
		}
		if exists {
			continue
		}

		daysUntil := breeding.DaysUntilBirth()
		title := "แจ้งเตือนกำหนดคลอด"
		if daysUntil <= 3 {
			title = "ใกล้กำหนดคลอดมาก"
		}
		if daysUntil <= 7 {
			title = "ใกล้กำหนดคลอด"
		}

		message := fmt.Sprintf("%s จะคลอดในอีก %d วัน (วันที่ %s)",
			breeding.MotherCodename,
			daysUntil,
			breeding.ExpectedBirth.Format("02/01/2006"))

		if err := s.CreateNotification(entities.TypeBirthDue, title, message); err != nil {
			return err
		}
		log.Println("Test", message, breeding.MotherCodename)
	}
	return nil

}

func (s *NotificationService) NotificationExistsToday(notiType, keyword string) (bool, error) {
	return s.noti.ExistsToday(notiType, keyword)
}
