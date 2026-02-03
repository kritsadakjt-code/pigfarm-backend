package routes

import (
	"backend/adapters/handlers"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PigRoute(app *fiber.App, handler *handlers.HttpPigHandler) {
	pig := app.Group("/pigs", middlewares.CheckJWT())
	// pig.Get("/", handler.FindAllPagination)
	pig.Post("/", handler.CreatePig)
	pig.Get("/search", controllers.SearchPigs)
	pig.Get("/", controllers.GetAllPigs)
	pig.Get("/:id", handler.GetPigByID)
	pig.Put("/:id", handler.UpdatePig)
	pig.Delete("/:id", handler.DeletePig)

}
