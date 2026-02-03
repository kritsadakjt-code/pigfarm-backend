package handlers

import (
	"backend/usecases"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	dashboardService *usecases.DashboardService
}

func NewDashboardHandler(dashboardService *usecases.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetDashboard(c *fiber.Ctx) error {
	role := c.Locals("role").(string)
	if role != "owner" && role != "employee" {
		return fiber.NewError(403, "Unauthorized")
	}

	resp, err := h.dashboardService.GetDashboardData(role)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

func (h *DashboardHandler) GetIncomeByMonthRange(c *fiber.Ctx) error {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		return fiber.NewError(400, "กรุณาระบุ start และ end เช่น ?start=2025-01&end=2025-06")
	}

	// เเปลงเวลาเป็น 2006-01
	layout := "2006-01"
	startDate, err1 := time.Parse(layout, startStr)
	endDate, err2 := time.Parse(layout, endStr)
	if err1 != nil || err2 != nil {
		return fiber.NewError(400, "รูปแบบเดือนไม่ถูกต้อง (ต้องเป็น YYYY-MM)")
	}
	// บวกเพิ่ม 1 เดือนโดยนับจนถึงวันที่ 1 ของเดือนถัดไป เช่น 10 ม.ค. - 1 ก.พ.
	endDate = endDate.AddDate(0, 1, 0)

	incomes, err := h.dashboardService.GetIncomeByMonthRange(startDate, endDate)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"start":  startStr,
		"end":    endStr,
		"income": incomes,
	})
}

func (h *DashboardHandler) GetExpenseByMonthRange(c *fiber.Ctx) error {
	startStr := c.Query("start")
	endStr := c.Query("end")

	if startStr == "" || endStr == "" {
		return fiber.NewError(400, "กรุณาระบุ start และ end เช่น ?start=2025-01&end=2025-06")
	}

	layout := "2006-01"
	startDate, err1 := time.Parse(layout, startStr)
	endDate, err2 := time.Parse(layout, endStr)
	if err1 != nil || err2 != nil {
		return fiber.NewError(400, "รูปแบบเดือนไม่ถูกต้อง (ต้องเป็น YYYY-MM)")
	}

	endDate = endDate.AddDate(0, 1, 0)

	expenses, err := h.dashboardService.GetExpenseByMonthRange(startDate, endDate)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"start":    startStr,
		"end":      endStr,
		"expenses": expenses,
	})
}
