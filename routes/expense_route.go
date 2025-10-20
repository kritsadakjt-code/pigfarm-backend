package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ExpenseRoute(app *fiber.App) {
	expense := app.Group("/expense", middlewares.CheckJWT(), middlewares.IsAdminOnly)
	expense.Post("/", controllers.CreateExpense)
	expense.Get("/search", controllers.SearchExpense)
	expense.Get("/", controllers.GetAllExpense)
	expense.Get("/:id", controllers.GetExpenseByID)
	expense.Put("/:id", controllers.UpdateExpense)
	expense.Delete("/:id", controllers.DeleteExpense)

}
