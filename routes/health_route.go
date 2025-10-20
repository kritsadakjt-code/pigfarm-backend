package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func HealthRoute(app *fiber.App) {
	health := app.Group("/health", middlewares.CheckJWT())
	health.Post("/", controllers.CreateHealth)
	health.Get("/search", controllers.SearchHealth)
	health.Get("/", controllers.GetAllHealth)
	health.Get("/:id", controllers.GetHealthByID)
	health.Put("/:id", controllers.UpdateHealth)
	health.Delete("/:id", controllers.DeleteHealth)
}
