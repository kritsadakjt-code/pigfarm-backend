package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FoodTypeRoute(app *fiber.App) {
	food_type := app.Group("/food-types", middlewares.CheckJWT())
	food_type.Post("/", controllers.CreateFoodType)
	food_type.Get("/", controllers.GetAllFoodType)
	food_type.Put("/:id", controllers.UpdateFoodType)
	food_type.Delete("/:id", controllers.DeleteFoodType)
}
