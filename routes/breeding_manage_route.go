package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func BreedingRoute(app *fiber.App) {
	breeding := app.Group("breeding", middlewares.CheckJWT())
	breeding.Post("/", controllers.CreateBreeding)
	breeding.Get("/search", controllers.SearchBreeding)
	breeding.Get("/", controllers.GetAllBreeding)
	breeding.Get("/:id", controllers.GetBreedingByID)
	breeding.Put("/:id", controllers.UpdateBreeding)
	breeding.Delete("/:id", controllers.DeleteBreeding)
}
