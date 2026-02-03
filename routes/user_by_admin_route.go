package routes

import (
	"backend/adapters/handlers"
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func UserByAdminRoute(app *fiber.App, handler *handlers.HttpUserHandlers) {
	admin := app.Group("/users", middlewares.CheckJWT(), middlewares.IsAdminOnly)
	// ?search="full_name"
	// สําหรับ pagination
	// admin.Get("/", handler.FindAllPagination)
	admin.Get("/search", controllers.SearchUser)
	admin.Post("/", handler.CreateUserByAdmin)
	admin.Get("/", handler.GetAllUsers)
	admin.Get("/:id", handler.FindByID)
	admin.Put("/:id", handler.UpdateUserByAdmin)
	admin.Delete("/:id", handler.DeleteUser)
	admin.Put("/approve/:id", handler.ApproveUser)
	admin.Put("/reject/:id", handler.RejectUser)
}
