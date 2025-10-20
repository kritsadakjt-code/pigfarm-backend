package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func EmailPasswordRoute(app *fiber.App) {
	password := app.Group("/api")
	password.Post("/forget-password", controllers.ForgetPassword)
	password.Post("/reset-password/:token", controllers.ResetPassword)
	password.Put("/change-password", middlewares.CheckJWT(), controllers.ChangePassword)
	password.Put("/change-email-request", middlewares.CheckJWT(), controllers.RequestChangeEmail)
	password.Get("/verify-email-change/:token", controllers.VerifyEmailForChange)
	password.Get("/verify-email-register/:token", controllers.VerifyEmailForRegister)
}
