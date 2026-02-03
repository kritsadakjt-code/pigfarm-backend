package usecases

import (
	"backend/dto"
	"time"
)

type DashboardRepository interface {
	GetPigCounts() (map[string]int64, error)
	GetBreedingStat() (map[string]int64, error)
	GetFoodStockStats() (map[string]float64, error)
	GetMonthlyExpense(start, end time.Time) (map[string]float64, error)
	GetMonthlyIncome(start, ernd time.Time) (map[string]float64, error)
	GetDairyFeeding() (map[string]float64, error)

	GetIncomeRange(start, end time.Time) ([]dto.MonthlyIncome, error)
	GetExpenseRange(start, end time.Time) ([]dto.MonthlyExpense, error)
}
