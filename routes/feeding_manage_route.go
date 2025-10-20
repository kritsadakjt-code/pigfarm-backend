package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FeedingRoute(app *fiber.App) {
	feeding := app.Group("/feeding", middlewares.CheckJWT())
	feeding.Post("/", controllers.CreateFeeding)
	feeding.Get("/search", controllers.SearchFeeding)
	feeding.Get("/", controllers.GetAllFeeding)
	feeding.Get("/:id", controllers.GetFeedingByID)
	feeding.Put("/:id", controllers.UpdateFeeding)
	feeding.Delete("/:id", controllers.DeleteFeeding)

}
