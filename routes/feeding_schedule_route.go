package routes

import (
	"backend/controllers"
	"backend/middlewares"

	"github.com/gofiber/fiber/v2"
)

func FeedingScheduleRoute(app *fiber.App) {
	// สร้างกลุ่ม Route ใหม่สำหรับจัดการตารางเวลา
	// ใช้ middlewares.IsAdminOnly เพื่อจำกัดสิทธิ์ให้เฉพาะ "owner"
	schedule := app.Group("/feeding-schedules", middlewares.CheckJWT(), middlewares.IsAdminOnly)

	schedule.Post("/", controllers.CreateFeedingSchedule)
	schedule.Get("/", controllers.GetAllFeedingSchedules)
	schedule.Get("/:id", controllers.GetFeedingScheduleByID)
	schedule.Put("/:id", controllers.UpdateFeedingSchedule)
	schedule.Delete("/:id", controllers.DeleteFeedingSchedule)
}
