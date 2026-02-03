package routes

import (
	"backend/adapters/handlers"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FeedingRoute(app *fiber.App, handler *handlers.FeedingHttpHandler) {
	feeding := app.Group("/feeding", middlewares.CheckJWT())
	feeding.Post("/", handler.CreateFeeding)
	// feeding.Get("/", handler.GetAllFeedingPagination)
	feeding.Get("/search", controllers.SearchFeeding)
	feeding.Get("/", controllers.GetAllFeeding)
	feeding.Get("/:id", handler.GetFeedingByID)
	feeding.Put("/:id", handler.UpdateFeeding)
	feeding.Delete("/:id", handler.DeleteFeeding)

}
