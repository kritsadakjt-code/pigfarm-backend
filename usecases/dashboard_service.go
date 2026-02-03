package usecases

import (
	"backend/dto"
	"time"
)

type DashboardService struct {
	DashboardRepo DashboardRepository
}

func NewDashboardService(DashboardRepo DashboardRepository) *DashboardService {
	return &DashboardService{DashboardRepo: DashboardRepo}
}

func (s *DashboardService) GetDashboardData(role string) (*dto.DashboardResponse, error) {
	pigCounts, err := s.DashboardRepo.GetPigCounts()
	if err != nil {
		return nil, err
	}
	// หาจํานวนหมูทั้งหมด
	var total int64
	for _, pigCount := range pigCounts {
		total += pigCount
	}
	// สถิติการผสมพันธุ์
	breedingStats, err := s.DashboardRepo.GetBreedingStat()
	if err != nil {
		return nil, err
	}
	// ปริมาณอาหาร
	foodStockStats, err := s.DashboardRepo.GetFoodStockStats()
	if err != nil {
		return nil, err
	}
	// ค่าใช้จ่าย
	expense := make(map[string]float64)
	income := make(map[string]float64)
	if role == "owner" {
		// ดึงจากเวลาปัจจุบัน ถ้าเวลาปัจจุบันไม่มีข้อมูลใน db ผลลัพธ์จะไม่ออก
		now := time.Now()
		// วันที่ 1 ของเดือนนี้
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		// วันที่ 1 ของเดือนถัดไป
		endOfMonth := startOfMonth.AddDate(0, 1, 0)
		var err error
		expense, err = s.DashboardRepo.GetMonthlyExpense(startOfMonth, endOfMonth)
		if err != nil {
			return nil, err
		}
		income, err = s.DashboardRepo.GetMonthlyIncome(startOfMonth, endOfMonth)
		if err != nil {
			return nil, err
		}

	}

	feedingStats, err := s.DashboardRepo.GetDairyFeeding()
	if err != nil {
		return nil, err
	}

	resp := dto.DashboardResponse{
		TotalPigs:           total,
		Fathers:             pigCounts["พ่อพันธุ์"],
		Mothers:             pigCounts["เเม่พันธุ์"],
		FatteningPigs:       pigCounts["หมูขุน"],
		Piglets:             pigCounts["ลูกหมู"],
		BreedingWaiting:     breedingStats["รอผล"],
		BreedingSuccess:     breedingStats["สําเร็จ"],
		BreedingFail:        breedingStats["ไม่สําเร็จ"],
		MainFoodKg:          foodStockStats["อาหารหลัก"],
		SupplementFoodKg:    foodStockStats["อาหารเสริม"],
		MonthlyExpenseChart: expense,
		MonthlyIncome:       income,
		DailyFeedingSummary: feedingStats,
	}
	return &resp, err
}

func (s *DashboardService) GetIncomeByMonthRange(start, end time.Time) ([]dto.MonthlyIncome, error) {
	return s.DashboardRepo.GetIncomeRange(start, end)
}

func (s *DashboardService) GetExpenseByMonthRange(start, end time.Time) ([]dto.MonthlyExpense, error) {
	return s.DashboardRepo.GetExpenseRange(start, end)
}
