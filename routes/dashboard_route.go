package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoute(app *fiber.App) {
	dashboard := app.Group("dashboard", middlewares.CheckJWT())
	dashboard.Get("/", controllers.GetDashboard)
	dashboard.Get("/income", controllers.GetIncomeByMonthRange)
	dashboard.Get("/expense", controllers.GetExpenseByMonthRange)
}
