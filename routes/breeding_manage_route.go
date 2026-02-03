package routes

import (
	"backend/adapters/handlers"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func BreedingRoute(app *fiber.App, handler *handlers.BreedingHttpHandler) {
	breeding := app.Group("breeding", middlewares.CheckJWT())
	// breeding.Get("/", handler.FindAllPagination)
	breeding.Post("/", handler.CreateBreeding)
	breeding.Get("/search", controllers.SearchBreeding)
	breeding.Get("/", controllers.GetAllBreeding)
	breeding.Get("/:id", controllers.GetBreedingByID)
	breeding.Put("/:id", handler.UpdateBreeding)
	breeding.Delete("/:id", handler.DeleteBreeding)
}
