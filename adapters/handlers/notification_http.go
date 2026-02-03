package handlers

import (
	"backend/adapters/dto"
	"backend/usecases"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type HttpNotificationHandler struct {
	service *usecases.NotificationService
}

func NewHttpNotificationHandler(service *usecases.NotificationService) *HttpNotificationHandler {
	return &HttpNotificationHandler{service: service}
}

func (h *HttpNotificationHandler) CreateNotification(c *fiber.Ctx) error {
	var req dto.CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	if err := h.service.CreateNotification(req.Type, req.Title, req.Message); err != nil {
		if errors.Is(err, usecases.ErrInvalidNotification) {
			return c.Status(400).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to create notification",
		})
	}

	return c.Status(201).JSON(dto.MessageResponse{
		Message: "Notification created successfully",
	})

}

func (h *HttpNotificationHandler) GetAllNotifications(c *fiber.Ctx) error {
	notifications, err := h.service.GetAllNotifications()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch notifications",
		})
	}

	unreadCount, _ := h.service.GetUnreadCount()
	response := dto.ToNotificationListResponse(notifications, unreadCount)

	return c.JSON(response)

}

func (h *HttpNotificationHandler) GetNotificationByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid notification id",
		})
	}

	noti, err := h.service.GetNotificationByID(uint(id))
	if err != nil {
		if errors.Is(err, usecases.ErrNotificationNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "notification not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch notification",
		})
	}
	response := dto.ToNotificationResponse(noti)
	return c.JSON(response)
}

func (h *HttpNotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	count, err := h.service.GetUnreadCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to get unread count",
		})
	}
	return c.JSON(dto.UnreadCountResponse{Count: count})
}

func (h *HttpNotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid notification id",
		})
	}
	if err := h.service.MarkAsRead(uint(id)); err != nil {
		if errors.Is(err, usecases.ErrNotificationNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "notification not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to mark notification as read",
		})
	}
	return c.JSON(dto.MessageResponse{Message: "Notification marked as read"})
}

func (h *HttpNotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	if err := h.service.MarkAllAsRead(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to mark all notifications as read",
		})
	}

	return c.JSON(dto.MessageResponse{Message: "All notifications marked as read"})
}

func (h *HttpNotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid notification id",
		})
	}

	if err := h.service.DeleteNotification(uint(id)); err != nil {
		if errors.Is(err, usecases.ErrNotificationNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "notification not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete notification",
		})
	}

	return c.JSON(fiber.Map{"message": "delete succes",
		"id": id})
}
