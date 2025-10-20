package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PigRoute(app *fiber.App) {
	pig := app.Group("/pigs", middlewares.CheckJWT())

	pig.Post("/", controllers.CreatePig)
	pig.Get("/search", controllers.SearchPigs)
	pig.Get("/", controllers.GetAllPigs)
	pig.Get("/:id", controllers.GetPigByID)
	pig.Put("/:id", controllers.UpdatePig)
	pig.Delete("/:id", controllers.DeletePig)
}
