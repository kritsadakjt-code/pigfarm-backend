package routes

import (
	"backend/adapters/handlers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func EmailPasswordRoute(app *fiber.App, handler *handlers.HttpUserHandlers) {
	password := app.Group("/api")
	password.Post("/forget-password", handler.ForgetPassword)
	password.Post("/reset-password/:token", handler.ResetPassword)
	password.Put("/change-password", middlewares.CheckJWT(), handler.ChangePassword)
	password.Put("/change-email-request", middlewares.CheckJWT(), handler.RequestChangeEmail)
	password.Get("/verify-email-change/:token", handler.VerifyEmailForChange)
	password.Get("/verify-email-register/:token", handler.VerifyEmailForRegister)
}
