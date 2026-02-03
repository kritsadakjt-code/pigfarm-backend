package routes

import (
	"backend/adapters/handlers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ExpenseRoute(app *fiber.App, handler *handlers.HttpExpenseHandler) {
	expense := app.Group("/expense", middlewares.CheckJWT(), middlewares.IsAdminOnly)
	expense.Post("/", handler.CreateExpense)
	expense.Get("/search", handler.SearchExpense)
	expense.Get("/", handler.GetAllExpense)
	expense.Get("/:id", handler.GetExpenseByID)
	expense.Put("/:id", handler.UpdateExpense)
	expense.Delete("/:id", handler.DeleteExpense)

}
