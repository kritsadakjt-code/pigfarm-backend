package usecases

import (
	"backend/dto"
	"backend/models"
	"math"
	"time"
)

type StockService struct {
	stockRepo StockRepository
}

func NewStockService(stockRepo StockRepository) *StockService {
	return &StockService{stockRepo: stockRepo}
}

func (s *StockService) mapToResponse(f *models.FoodStock) *dto.FoodStockResponse {
	return &dto.FoodStockResponse{
		ID:          f.ID,
		FoodTypeID:  f.FoodTypeID,
		FoodName:    f.FoodType.Name,
		FoodType:    f.FoodType.Type,
		DateTime:    f.DateTime,
		Amount:      f.Amount,
		Note:        f.Note,
		CreatedName: f.Creator.FullName,
		UpdatedName: f.Updater.FullName,
	}
}
func (s *StockService) CreateFoodStock(input *dto.FoodStockInput, userID uint) (*dto.FoodStockResponse, error) {
	if input.Amount <= 0 {
		return nil, ErrInvalidStockAmount
	}
	exists, err := s.stockRepo.CheckDuplicate(input.FoodTypeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrFoodStockAlreadyExists
	}
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}

	parsedDate, err := time.ParseInLocation("2006-01-02 15:04", input.DateTime, loc)
	if err != nil {
		return nil, ErrInvalidDate
	}
	if parsedDate.After(time.Now()) {
		return nil, ErrCantFuture
	}

	foodStock := models.FoodStock{
		FoodTypeID: input.FoodTypeID,
		DateTime:   parsedDate,
		Amount:     input.Amount,
		Note:       input.Note,
		CreatedBy:  userID,
		UpdatedBy:  userID,
	}

	if err := s.stockRepo.Create(&foodStock); err != nil {
		return nil, err
	}

	resp, err := s.stockRepo.GetByID(foodStock.ID)
	if err != nil {
		return nil, err
	}
	return s.mapToResponse(resp), nil

}

func (s *StockService) UpdateFoodStock(id uint, userID uint, input dto.FoodStockUpdate) (*dto.FoodStockResponse, error) {
	stock, err := s.stockRepo.GetByID(id)
	if err != nil {
		return nil, ErrFoodNotFound
	}
	if input.DateTime != nil {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			return nil, err
		}
		parsedDate, err := time.ParseInLocation("2006-01-02 15:04", *input.DateTime, loc)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if parsedDate.After(time.Now()) {
			return nil, ErrCantFuture
		}
		stock.DateTime = parsedDate
	}

	if input.Amount != nil {
		if *input.Amount <= 0 {
			return nil, ErrInvalidStockAmount
		}
		stock.Amount = *input.Amount
	}

	if input.Note != nil {
		stock.Note = *input.Note
	}

	stock.UpdatedBy = userID
	if err := s.stockRepo.Update(stock); err != nil {
		return nil, err
	}

	// ดึงข้อมูลมาเเสดงใหม่
	resp, err := s.stockRepo.GetByID(id)
	if err != nil {
		return nil, ErrFoodNotFound
	}
	return s.mapToResponse(resp), err
}

func (s *StockService) GetAllPagi(input dto.ParamFoodStock) (*dto.FoodStockPagiResp, error) {
	if input.Page <= 0 {
		input.Page = 1
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	stocks, total, err := s.stockRepo.GetAllPagi(input)
	if err != nil {
		return nil, err
	}

	var resp []dto.FoodStockResponse
	for _, f := range stocks {
		resp = append(resp, *s.mapToResponse(&f))
	}
	lastPage := int(math.Ceil(float64(total) / float64(input.Limit)))

	return &dto.FoodStockPagiResp{
		Data:     resp,
		Total:    total,
		Page:     input.Page,
		LastPage: lastPage,
		Limit:    input.Limit,
	}, nil
}

func (s *StockService) DeleteFoodStock(id uint) error {
	_, err := s.stockRepo.GetByID(id)
	if err != nil {
		return ErrFoodNotFound
	}
	isUsed, err := s.stockRepo.IsUsedInFeeding(id)
	if err != nil {
		return err
	}
	if isUsed {
		return ErrFoodStockUsed
	}
	// check ว่ามีการใช้มั้ย
	return s.DeleteFoodStock(id)
}
