package routes

import (
	"backend/adapters/handlers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoute(app *fiber.App, handler *handlers.DashboardHandler) {
	dashboard := app.Group("dashboard", middlewares.CheckJWT())
	dashboard.Get("/", handler.GetDashboard)
	dashboard.Get("/income", handler.GetIncomeByMonthRange)
	dashboard.Get("/expense", handler.GetExpenseByMonthRange)
}
