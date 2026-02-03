package models

import (
	"backend/entities"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type NotificationModel struct {
	gorm.Model
	Type    string `json:"type" gorm:"not null"` // "food_low", "birth_due", "health_alert"
	Title   string `json:"title" gorm:"not null"`
	Message string `json:"message" gorm:"not null"`
	IsRead  bool   `json:"is_read" gorm:"default:false"`
}

func (NotificationModel) TableName() string {
	return "notifications"
}

// convert: moldel -> entity
func ToEntity(model NotificationModel) entities.Notification {
	var deletedAt *time.Time
	if model.DeletedAt.Valid {
		deletedAt = &model.DeletedAt.Time
	}

	return entities.Notification{
		ID:        strconv.FormatUint(uint64(model.ID), 10),
		Type:      model.Type,
		Title:     model.Title,
		Message:   model.Message,
		IsRead:    model.IsRead,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: deletedAt,
	}
}

// convert: entity -> moldel
func ToModel(entity *entities.Notification) NotificationModel {
	model := NotificationModel{
		Type:    entity.Type,
		Title:   entity.Title,
		Message: entity.Message,
		IsRead:  entity.IsRead,
	}
	if entity.ID != "" {
		id, _ := strconv.ParseInt(entity.ID, 10, 64)
		model.ID = uint(id)
	}
	model.CreatedAt = entity.CreatedAt
	model.UpdatedAt = entity.UpdatedAt

	return model
}
