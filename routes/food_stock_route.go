package routes

import (
	"backend/adapters/handlers"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FoodStockRoute(app *fiber.App, handler *handlers.HttpStockHandler) {
	foodStock := app.Group("/food-stock", middlewares.CheckJWT())
	foodStock.Post("/", handler.CreateFoodStock)
	// foodStock.Get("/", handler.GetAllPagi)
	foodStock.Get("/search", controllers.SearchFoodStock)
	foodStock.Get("/", controllers.GetAllFoodStock)
	foodStock.Get("/:id", controllers.GetFoodStockByID)
	foodStock.Put("/:id", handler.UpdateFoodStock)
	foodStock.Delete("/:id", handler.DeleteFoodStock)
}
