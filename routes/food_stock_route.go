package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FoodStockRoute(app *fiber.App) {
	foodStock := app.Group("/food-stock", middlewares.CheckJWT())
	foodStock.Post("/", controllers.CreateFoodStock)
	foodStock.Get("/search", controllers.SearchFoodStock)
	foodStock.Get("/", controllers.GetAllFoodStock)
	foodStock.Get("/:id", controllers.GetFoodStockByID)
	foodStock.Put("/:id", controllers.UpdateFoodStock)
	foodStock.Delete("/:id", controllers.DeleteFoodStock)
}
