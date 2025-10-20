package routes

import (
	"backend/adapters/handlers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func NotificationRoutes(app *fiber.App, handler *handlers.HttpNotificationHandler) {
	notis := app.Group("/notification", middlewares.CheckJWT())
	notis.Get("/", handler.GetAllNotifications)
	notis.Get("/unread-count", handler.GetUnreadCount)
	notis.Get("/:id", handler.GetNotificationByID)
	notis.Post("/", handler.CreateNotification)
	notis.Patch("/:id/read", handler.MarkAsRead)
	notis.Patch("/read-all", handler.MarkAllAsRead)
	notis.Delete("/:id", handler.DeleteNotification)
}
