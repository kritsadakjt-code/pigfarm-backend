package repositories

import (
	"backend/entities"
	"backend/mappers"
	"backend/models"
	"backend/usecases"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type GormNotificationRepository struct {
	db *gorm.DB
}

func NewGormNotificationRepository(db *gorm.DB) usecases.NotificationRepository {
	return &GormNotificationRepository{db: db}
}

func (r *GormNotificationRepository) Create(noti *entities.Notification) error {
	model := mappers.NotificationEntityToModel(noti)
	result := r.db.Create(&model)
	if result.Error != nil {
		return fmt.Errorf("database error: %w", result.Error)
	}
	// update entity เพื่อเอาไว้ response ข้อมูล
	noti.ID = strconv.FormatUint(uint64(model.ID), 10)
	noti.CreatedAt = model.CreatedAt
	noti.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *GormNotificationRepository) GetAll() ([]entities.Notification, error) {
	var model_noti []models.NotificationModel
	result := r.db.Order("created_at").Find(&model_noti)
	if result.Error != nil {
		return nil, fmt.Errorf("database error: %w", result.Error)
	}

	noti := make([]entities.Notification, len(model_noti))
	for i, model := range model_noti {
		noti[i] = mappers.NotificationModelToEntity(model)
	}
	return noti, nil

}

func (r *GormNotificationRepository) GetByID(id uint) (*entities.Notification, error) {
	var model_noti models.NotificationModel
	if err := r.db.First(&model_noti, id).Error; err != nil {
		return nil, err
	}
	noti := models.ToEntity(model_noti)
	return &noti, nil
}

func (r *GormNotificationRepository) GetUnreadCount() (int64, error) {
	var count int64
	if err := r.db.Model(&models.NotificationModel{}).Where("is_read = ?", false).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *GormNotificationRepository) MarkAsRead(id uint) error {
	result := r.db.Model(&models.NotificationModel{}).
		Where("id = ?", id).
		Update("is_read", true)
	if result.Error != nil {
		return fmt.Errorf("database error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return usecases.ErrNotificationNotFound
	}

	return nil
}

func (r *GormNotificationRepository) MarkAllAsRead() error {
	// result := r.db.Model(&models.NotificationModel{}).
	// 	Where("is_read = ?", false).
	// 	Update("is_read", true)

	// if result.Error != nil {
	// 	return fmt.Errorf("database error: %w", result.Error)
	// }

	// return nil
	err := r.db.Model(&models.NotificationModel{}).
		Where("is_read = ?", false).
		Update("is_read = ?", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *GormNotificationRepository) Delete(id uint) error {
	result := r.db.Unscoped().Delete(&models.NotificationModel{}, id)
	if result.Error != nil {
		return fmt.Errorf("database error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return usecases.ErrNotificationNotFound
	}

	return nil
}

func (r *GormNotificationRepository) ExistsToday(notiType, keyword string) (bool, error) {
	var count int64
	today := time.Now().Format("2006-01-02")

	result := r.db.Model(&models.NotificationModel{}).
		Where("type = ? AND message LIKE ? AND DATE(created_at) = ?",
			notiType, "%"+keyword+"%", today).Count(&count)
	if result.Error != nil {
		return false, fmt.Errorf("database error: %w", result.Error)
	}

	return count > 0, nil
}

// package repositorys

// import (
// 	"backend/models"
// 	"fmt"
// 	"time"

// 	"gorm.io/gorm"
// )

// type NotificationService struct {
// 	db *gorm.DB
// }

// func NewNotificationService(db *gorm.DB) *NotificationService {
// 	return &NotificationService{db: db}
// }

// func (ns *NotificationService) CreateNotification(notiType, title, message string) error {
// 	notification := models.Notification{
// 		Type:    notiType,
// 		Title:   title,
// 		Message: message,
// 	}
// 	return ns.db.Create(&notification).Error
// }

// // check ซั้าเเต่ละวัน
// func (ns *NotificationService) NotificationExistsToday(notiType, keyword string) bool {
// 	var count int64
// 	today := time.Now().Format("2006-01-02")
// 	ns.db.Model(&models.Notification{}).
// 		Where("type = ? AND message LIKE ? AND DATE(created_at) = ?", notiType, "%"+keyword+"%", today).
// 		Count(&count)
// 	return count > 0
// }

// // check food and create
// func (ns *NotificationService) CheckFoodLowStock(threshold float64) error {
// 	var foods []models.FoodStock
// 	// เอาเฉพาะ amount ที่น้อยกว่า minAmount
// 	if err := ns.db.Where("amount < ? AND amount > 0", threshold).Find(&foods).Error; err != nil {
// 		return err
// 	}

// 	for _, food := range foods {
// 		if !ns.NotificationExistsToday("food_low", food.Name) {
// 			title := "อาหารใกล้หมด"
// 			message := fmt.Sprintf("%s เหลือเพียง %.2f Kg", food.Name, food.Amount)
// 			if food.Amount < 5 {
// 				title = "อาหารใกล้หมดอย่างเร่งด่วน"
// 			}
// 			if err := ns.CreateNotification("food_low", title, message); err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

// // check date and create
// // dayAhead คือวัน ก่อนคลอดกี่วัน
// func (ns *NotificationService) CheckUpcomingBirths(daysAhead int) error {
// 	var breedings []models.Breeding
// 	startDate := time.Now()
// 	endDate := time.Now().AddDate(0, 0, daysAhead)

// 	if err := ns.db.Preload("Mother").
// 		Where("status = ? AND expected_birth BETWEEN ? AND ?", "อุ้มท้อง", startDate, endDate).
// 		Find(&breedings).Error; err != nil {
// 		return err
// 	}

// 	for _, breeding := range breedings {
// 		if !ns.NotificationExistsToday("birth_due", breeding.Mother.CodeName) {
// 			daysUntilBirth := int(time.Until(breeding.ExpectedBirth).Hours() / 24)
// 			title := "แจ้งเตือนกำหนดคลอด"
// 			if daysUntilBirth <= 7 {
// 				title = "ใกล้กำหนดคลอด"
// 			}
// 			if daysUntilBirth <= 3 {
// 				title = "ใกล้กำหนดคลอดมาก"

// 			}
// 			message := fmt.Sprintf("%s จะคลอดในอีก %d วัน (วันที่ %s)",
// 				breeding.Mother.CodeName, daysUntilBirth, breeding.ExpectedBirth.Format("02/01/2006"))
// 			ns.CreateNotification("birth_due", title, message)
// 		}
// 	}
// 	return nil
// }

// // เรียกตรวจสอบทั้งหมด
// func (ns *NotificationService) RunAllChecks() error {
// 	if err := ns.CheckFoodLowStock(10); err != nil {
// 		return err
// 	}
// 	if err := ns.CheckUpcomingBirths(7); err != nil {
// 		return err
// 	}

// 	return nil
// }
