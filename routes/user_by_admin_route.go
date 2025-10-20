package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserByAdminRoute(app *fiber.App) {
	admin := app.Group("/users", middlewares.CheckJWT(), middlewares.IsAdminOnly)
	// ?search="full_name"
	admin.Get("/search", controllers.SearchUser)
	admin.Post("/", controllers.CreateUser)
	admin.Get("/", controllers.GetAllUsers)
	admin.Get("/:id", controllers.GetUserByID)
	admin.Put("/:id", controllers.UpdateUser)
	admin.Delete("/:id", controllers.DeleteUser)
	admin.Put("/approve/:id", controllers.ApproveUser)
	admin.Put("/reject/:id", controllers.RejectUser)
}
