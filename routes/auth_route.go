package routes

import (
	"backend/adapters/handlers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, handler *handlers.HttpUserHandlers) {
	auth := app.Group("/auth")
	auth.Post("/register", handler.CreateUser)
	auth.Post("/login", handler.Login)
	auth.Post("/resend-verify", handler.ResendVerifyEmail)

}
