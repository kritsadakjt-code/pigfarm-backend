package mappers

import (
	"backend/entities"
	"backend/models"
	"strconv"
	"time"
)

func NotificationModelToEntity(model models.NotificationModel) entities.Notification {
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

func NotificationEntityToModel(entity *entities.Notification) models.NotificationModel {
	model := models.NotificationModel{
		Type:    entity.Type,
		Title:   entity.Title,
		Message: entity.Message,
		IsRead:  entity.IsRead,
	}
	if entity.ID != "" {
		id, _ := strconv.ParseUint(entity.ID, 10, 64)
		model.ID = uint(id)
	}

	model.CreatedAt = entity.CreatedAt
	model.UpdatedAt = entity.UpdatedAt
	return model
}

func NotificationToEntities(models []models.NotificationModel) []entities.Notification {
	entities := make([]entities.Notification, len(models))
	for i, m := range models {
		entities[i] = NotificationModelToEntity(m)
	}
	return entities
}
