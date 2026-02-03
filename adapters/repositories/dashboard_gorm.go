package repositories

import (
	"backend/dto"
	"backend/models"
	"backend/usecases"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type DashboardGormRepo struct {
	db *gorm.DB
}

func NewDashboardGormRepo(db *gorm.DB) usecases.DashboardRepository {
	return &DashboardGormRepo{db: db}
}
func (r *DashboardGormRepo) GetPigCounts() (map[string]int64, error) {
	type PigCountResult struct {
		Type  string
		Count int64
	}
	var results []PigCountResult
	err := r.db.Model(&models.Pig{}).
		Select("type, count(*) as count").
		Group("type").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Type] = result.Count
	}
	return counts, nil
}
func (r *DashboardGormRepo) GetBreedingStat() (map[string]int64, error) {
	type BreedingResult struct {
		Result string
		Count  int64
	}

	var results []BreedingResult
	err := r.db.Model(&models.Breeding{}).
		Select("result, count(*) as count").
		Group("result").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]int64)
	for _, r := range results {
		stats[r.Result] = r.Count
	}

	return stats, nil
}

func (r *DashboardGormRepo) GetFoodStockStats() (map[string]float64, error) {
	type FoodResult struct {
		Type  string
		Total float64
	}

	var results []FoodResult

	err := r.db.Model(&models.FoodStock{}).
		Select("food_types.type, COALESCE(SUM(food_stocks.amount), 0) as total").
		Joins("JOIN food_types ON food_types.id = food_stocks.food_type_id").
		Group("food_types.type").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	stats := make(map[string]float64)
	for _, r := range results {
		stats[r.Type] = r.Total
	}
	return stats, nil
}

func (r *DashboardGormRepo) GetMonthlyExpense(start, end time.Time) (map[string]float64, error) {
	type ExpenseResult struct {
		Category string
		Total    float64
	}

	var results []ExpenseResult
	err := r.db.Model(&models.Expense{}).
		Select("category, COALESCE(SUM(amount),0) as total").
		Where("date >= ? AND date < ?", start, end).
		Group("category").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	fmt.Println(results)
	stats := make(map[string]float64)
	for _, r := range results {
		stats[r.Category] = r.Total

	}

	return stats, nil
}

func (r *DashboardGormRepo) GetDairyFeeding() (map[string]float64, error) {
	type FeedingResult struct {
		Date  string
		Total float64
	}

	var results []FeedingResult
	err := r.db.Model(&models.Feeding{}).
		Select("TO_CHAR(date_time, 'YYYY-MM-DD') as date, COALESCE(SUM(amount), 0) as total ").
		Group("TO_CHAR(date_time, 'YYYY-MM-DD')").
		Order("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]float64)
	for _, r := range results {
		stats[r.Date] = r.Total
	}
	return stats, nil
}
func (r *DashboardGormRepo) GetMonthlyIncome(start, end time.Time) (map[string]float64, error) {
	type IncomeResult struct {
		Date  string
		Total float64
	}

	var results []IncomeResult
	err := r.db.Model(&models.PigSale{}).
		Select("date, COALESCE(sum(total_price),0) as total0").
		Where("date >= ? AND date < ?", start, end).
		Group("date").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	stats := make(map[string]float64)
	for _, r := range results {
		stats[r.Date] = r.Total
	}
	return stats, nil
}

func (r *DashboardGormRepo) GetIncomeRange(start, end time.Time) ([]dto.MonthlyIncome, error) {
	var incomes []dto.MonthlyIncome

	err := r.db.Model(&models.PigSale{}).
		Select("TO_CHAR(date, 'Mon YYYY') as month, COALESCE(sum(total_price),0) as total").
		Where("date >= ? AND date < ?", start, end).
		Group("TO_CHAR(date, 'Mon YYYY'), date_trunc('month', date)").
		Order("date_trunc('month', date) ASC").
		Scan(&incomes).Error
	return incomes, err
}

func (r *DashboardGormRepo) GetExpenseRange(start, end time.Time) ([]dto.MonthlyExpense, error) {
	var expense []dto.MonthlyExpense

	err := r.db.Model(&models.Expense{}).
		Select(`TO_CHAR(date, 'Mon YYYY') as month,
	category, 
	COALESCE(SUM(amount), 0) as total`).
		Where("date >= ? AND date < ?", start, end).
		Group("TO_CHAR(date, 'Mon YYYY'), category, date_trunc('month', date)"). // data_trunc ให้เรียงต่อจาก order_by
		// data_trunc เพื่อให้เรียนตามปฏิทิน ถ้าไม่มีจะเรียงจาก a-z เเละจะ order_by อะไรค่านั้นต้องอยู่ใน group_by ด้วย
		Order("date_trunc('month',date) ASC, category ASC").
		Scan(&expense).Error

	return expense, err

}
