package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetDashboard(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role != "owner" && role != "employee" {
		return c.Status(403).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var monthlyIncome float64
	// ดึงจํานวนหมู
	var total, father, mother, fattening, piglet int64
	pig := models.Pig{}
	config.DB.Model(pig).Count(&total)
	config.DB.Model(pig).Where("type = ?", "พ่อพันธุ์").Count(&father)
	config.DB.Model(pig).Where("type = ?", "เเม่พันธุ์").Count(&mother)
	config.DB.Model(pig).Where("type = ?", "หมูขุน").Count(&fattening)
	config.DB.Model(pig).Where("type = ?", "ลูกหมู").Count(&piglet)

	// ดึงการผสมพันธุ์
	var success, fail, waitResult int64
	config.DB.Model(&models.Breeding{}).Where("result = ?", "รอผล").Count(&waitResult)
	config.DB.Model(&models.Breeding{}).Where("result = ?", "สําเร็จ").Count(&success)
	config.DB.Model(&models.Breeding{}).Where("result = ?", "ไม่สําเร็จ").Count(&fail)

	// ดึงปริมาณอาหาร
	var mainFood, supplementFood float64
	config.DB.Model(&models.FoodStock{}).Select("COALESCE(SUM(amount),0)").Where("type = ?", "อาหารหลัก").Scan(&mainFood)
	config.DB.Model(&models.FoodStock{}).Select("COALESCE(SUM(amount),0)").Where("type = ?", "อาหารเสริม").Scan(&supplementFood)

	// ดึงค่าใช้จ่ายเดือนนี้
	expenseMap := map[string]float64{}
	if role == "owner" {
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0)

		type ExpenseResult struct {
			Category string
			Total    float64
		}
		var results []ExpenseResult
		config.DB.Model(&models.Expense{}).Select("category, SUM(amount) as total").
			Where("date >= ? AND date <= ?", startOfMonth, endOfMonth).
			Group("category").
			Scan(&results)

		for _, r := range results {
			expenseMap[r.Category] = r.Total
		}
		// คำนวณผลรวมของ total_price จากตาราง pig_sales ในเดือนปัจจุบัน
		config.DB.Model(&models.PigSale{}).
			Select("COALESCE(SUM(total_price), 0)").
			Where("date >= ? AND date < ?", startOfMonth, endOfMonth).
			Scan(&monthlyIncome)

	}
	// สรุปอาหารที่ให้ต่อวัน
	type DairyFeeding struct {
		Date  string
		Total float64
	}

	var DairyFeedings []DairyFeeding
	config.DB.Model(&models.Feeding{}).
		Select("DATE(date_time) as date, SUM(amount) as total").
		Group("DATE(date_time)").
		Order("DATE(date_time)").
		Scan(&DairyFeedings)
	feedingSummary := map[string]float64{}
	for _, f := range DairyFeedings {
		feedingSummary[f.Date] = f.Total
	}

	resp := dto.DashboardResponse{
		TotalPigs:           total,
		Fathers:             father,
		Mothers:             mother,
		FatteningPigs:       fattening,
		Piglets:             piglet,
		BreedingWaiting:     waitResult,
		BreedingSuccess:     success,
		BreedingFail:        fail,
		MainFoodKg:          mainFood,
		SupplementFoodKg:    supplementFood,
		MonthlyIncome:       monthlyIncome,
		MonthlyExpenseChart: expenseMap,
		DailyFeedingSummary: feedingSummary,
	}
	return c.JSON(resp)
}

func GetIncomeByMonthRange(c *fiber.Ctx) error {
	start := c.Query("start") // เช่น 2025-01
	end := c.Query("end")     // เช่น 2025-06

	if start == "" || end == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "กรุณาระบุ start และ end เช่น ?start=2025-01&end=2025-06",
		})
	}

	// แปลงเดือนให้เป็นวันที่เริ่มและสิ้นสุดจริง
	startDate, err1 := time.Parse("2006-01", start)
	endDate, err2 := time.Parse("2006-01", end)
	if err1 != nil || err2 != nil {
		return c.Status(400).JSON(fiber.Map{"error": "รูปแบบเดือนไม่ถูกต้อง (ต้องเป็น YYYY-MM)"})
	}

	// บวก 1 เดือนจาก end เพื่อให้ครอบคลุมเดือนสุดท้าย
	endDate = endDate.AddDate(0, 1, 0)

	type MonthlyIncome struct {
		Month string  `json:"month"`
		Total float64 `json:"total"`
	}

	var incomes []MonthlyIncome

	config.DB.Model(&models.PigSale{}).
		Select(`
            TO_CHAR(date, 'Mon YYYY') as month,
            EXTRACT(MONTH FROM date) as month_num,
            COALESCE(SUM(total_price), 0) as total
        `).
		Where("date >= ? AND date < ?", startDate, endDate).
		Group("month, month_num").
		Order("month_num").
		Scan(&incomes)

	return c.JSON(fiber.Map{
		"start":  start,
		"end":    end,
		"income": incomes,
	})
}

func GetExpenseByMonthRange(c *fiber.Ctx) error {
	start := c.Query("start") // เช่น 2025-01
	end := c.Query("end")     // เช่น 2025-06

	if start == "" || end == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "กรุณาระบุ start และ end เช่น ?start=2025-01&end=2025-06",
		})
	}

	startDate, err1 := time.Parse("2006-01", start)
	endDate, err2 := time.Parse("2006-01", end)
	if err1 != nil || err2 != nil {
		return c.Status(400).JSON(fiber.Map{"error": "รูปแบบเดือนไม่ถูกต้อง (ต้องเป็น YYYY-MM)"})
	}

	endDate = endDate.AddDate(0, 1, 0) // ครอบคลุมเดือนสุดท้าย

	type MonthlyExpense struct {
		Month    string  `json:"month"`
		Category string  `json:"category"`
		Total    float64 `json:"total"`
	}

	var expenses []MonthlyExpense

	config.DB.Model(&models.Expense{}).
		Select(`
            TO_CHAR(date, 'Mon YYYY') as month,
            category,
            COALESCE(SUM(amount),0) as total
        `).
		Where("date >= ? AND date < ?", startDate, endDate).
		Group("month, category").
		Order("MIN(date) ASC").
		Scan(&expenses)

	return c.JSON(fiber.Map{
		"start":    start,
		"end":      end,
		"expenses": expenses,
	})
}
