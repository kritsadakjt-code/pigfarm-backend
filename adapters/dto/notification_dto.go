package dto

import (
	"backend/entities"
	"time"
)

type CreateNotificationRequest struct {
	Type    string `json:"type" validate:"required,oneof=food_low birth_due"`
	Title   string `json:"title" validate:"required,min=1,max=255"`
	Message string `json:"message" validate:"required,min=1"`
}

type NotificationResponse struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	Total         int                    `json:"total"`
	UnreadCount   int64                  `json:"unread_count"`
}

type UnreadCountResponse struct {
	Count int64 `json:"count"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// convert
func ToNotificationResponse(entity *entities.Notification) NotificationResponse {
	return NotificationResponse{
		ID:        entity.ID,
		Type:      entity.Type,
		Title:     entity.Title,
		Message:   entity.Message,
		IsRead:    entity.IsRead,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}
}

func ToNotificationListResponse(entities []entities.Notification, unreadCount int64) NotificationListResponse {
	responses := make([]NotificationResponse, len(entities))
	for i, entity := range entities {
		responses[i] = ToNotificationResponse(&entity)
	}

	return NotificationListResponse{
		Notifications: responses,
		Total:         len(entities),
		UnreadCount:   unreadCount,
	}
}
