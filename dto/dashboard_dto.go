package dto

type DashboardResponse struct {
	TotalPigs           int64              `json:"total_pigs"`
	Fathers             int64              `json:"fathers"`
	Mothers             int64              `json:"mothers"`
	FatteningPigs       int64              `json:"fattening_pigs"`
	Piglets             int64              `json:"piglets"`
	BreedingWaiting     int64              `json:"breeding_waiting"`
	BreedingSuccess     int64              `json:"breeding_success"`
	BreedingFail        int64              `json:"breeding_fail"`
	MainFoodKg          float64            `json:"main_food_kg"`
	SupplementFoodKg    float64            `json:"supplement_food_kg"`
	MonthlyExpenseChart map[string]float64 `json:"monthly_expense_chart"` // ค่าใช้จ่ายรายเดือนตามประเภท
	MonthlyIncome       float64            `json:"monthly_income"`
	DailyFeedingSummary map[string]float64 `json:"daily_feeding_summary"` // อาหารที่ให้ต่อวัน
}
