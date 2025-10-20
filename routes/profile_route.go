package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(app *fiber.App) {
	profile := app.Group("/api")
	profile.Get("/profile", middlewares.CheckJWT(), controllers.GetProfile)
	profile.Put("/profile", middlewares.CheckJWT(), controllers.UpdateProfile)
}
