package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func PigSaleRoute(app *fiber.App) {
	pig := app.Group("pigsale", middlewares.CheckJWT(), middlewares.IsAdminOnly)
	pig.Post("/", controllers.CreatePigSale)
	pig.Get("/search", controllers.SearchPigSale)
	pig.Get("/", controllers.GetAllPigSale)
	pig.Get("/:id", controllers.GetPigSaleByID)
	pig.Put("/:id", controllers.UpdatePigSale)
	pig.Delete("/:id", controllers.DeletePigSale)
}
