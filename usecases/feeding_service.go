package usecases

import (
	"backend/dto"
	"backend/models"
	"errors"
	"math"
	"time"

	"gorm.io/gorm"
)

type FeedingService struct {
	feedingRepo FeedingRepository
}

func NewFeedingService(feedingRepo FeedingRepository) *FeedingService {
	return &FeedingService{feedingRepo: feedingRepo}
}

func (s *FeedingService) CreateFeeding(input dto.FeedingInput, userID uint) (*dto.FeedingResponse, error) {
	if input.Amount < 0 {
		return nil, ErrFoodNotZero
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

	// เตรียม model เพื่อส่งข้อมูลเช็คเเละบันทึก
	feeding := models.Feeding{
		FoodID:    input.FoodID,
		DateTime:  parsedDate,
		Amount:    input.Amount,
		Note:      input.Note,
		CreatedBy: userID,
		UpdatedBy: userID,
	}
	validPigIDs, err := s.feedingRepo.Create(&feeding, input.PigIDs)
	if err != nil {
		return nil, err
	}

	createdFeeding, err := s.feedingRepo.GetById(feeding.ID)
	if err != nil {
		return nil, err
	}

	resp := dto.FeedingResponse{
		ID:          createdFeeding.ID,
		FoodID:      createdFeeding.FoodID,
		PigIDs:      validPigIDs,
		DateTime:    createdFeeding.DateTime,
		Amount:      createdFeeding.Amount,
		Note:        createdFeeding.Note,
		CreatedName: createdFeeding.Creator.FullName,
		UpdatedName: createdFeeding.Updater.FullName,
	}

	return &resp, nil
}

func (s *FeedingService) UpdateFeeding(id uint, input dto.FeedingUpdate, userID uint) (*dto.FeedingResponse, error) {
	existing, err := s.feedingRepo.GetById(id)
	if err != nil {
		return nil, ErrFeedingNotFound
	}
	// prepare ก่อนเพื่อไม่ให้ค่าเป็นค่าว่าง เพราะจะส่งค่าไปเปรียบเทียบ
	updatedFeeding := models.Feeding{
		// ID:        existing.ID,
		FoodID:    existing.FoodID,
		DateTime:  existing.DateTime,
		Amount:    existing.Amount,
		Note:      existing.Note,
		CreatedBy: existing.CreatedBy, // คนสร้างคนเดิม
		UpdatedBy: userID,             // คนแก้คือคนปัจจุบัน
	}
	// ถ้ามีการส่งค่าใหม่มา
	if input.FoodID != nil {
		updatedFeeding.FoodID = *input.FoodID
	}
	if input.Amount != nil {
		if *input.Amount <= 0 {
			return nil, ErrInvalidInput
		}
		updatedFeeding.Amount = *input.Amount
	}
	if input.Note != nil {
		updatedFeeding.Note = *input.Note
	}
	if input.DateTime != nil {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			return nil, err
		}
		parsedDate, err := time.ParseInLocation("2006-01-02 15:04", *input.DateTime, loc)
		if err != nil {
			return nil, err
		}
		if parsedDate.After(time.Now()) {
			return nil, ErrCantFuture
		}
		updatedFeeding.DateTime = parsedDate
	}

	var newPigIDs []uint
	if input.PigIDs != nil {
		newPigIDs = *input.PigIDs
	} else {
		// ถ้าไม่ส่งหมูใหม่มาให้ใช้ค่าเดิม
		for _, item := range existing.Items {
			newPigIDs = append(newPigIDs, item.PigID)
		}
	}
	if err := s.feedingRepo.Update(id, &updatedFeeding, newPigIDs); err != nil {
		return nil, err
	}
	return s.GetFeedingByID(id)
}

func (s *FeedingService) GetAllFeedingPagination(param dto.ParamFeeding) (*dto.FeedingPagiResp, error) {
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	feedings, total, err := s.feedingRepo.GetAll(param)
	if err != nil {
		return nil, err
	}
	var resp []dto.FeedingResponse
	for _, f := range feedings {
		// ดึงข้อมูลหมู
		var pigIDs []uint
		var pigCodes []string
		for _, item := range f.Items {
			pigIDs = append(pigIDs, item.PigID)
			if item.Pig.CodeName != "" {
				pigCodes = append(pigCodes, item.Pig.CodeName)
			}
		}
		resp = append(resp, dto.FeedingResponse{
			ID:          f.ID,
			FoodID:      f.FoodID,
			FoodName:    f.FoodStock.FoodType.Name, // ต้องมั่นใจว่า Preload FoodStock มาแล้ว
			PigIDs:      pigIDs,
			PigCodeName: pigCodes,
			DateTime:    f.DateTime,
			Amount:      f.Amount,
			Note:        f.Note,
			CreatedName: f.Creator.FullName,
			CreatedRole: f.Creator.Role,
			UpdatedName: f.Updater.FullName,
			UpdatedRole: f.Updater.Role,
		})
	}
	lastPage := int(math.Ceil(float64(total)) / float64(param.Limit))
	return &dto.FeedingPagiResp{
		Data:     resp,
		Total:    total,
		Page:     param.Page,
		LastPage: lastPage,
		Limit:    param.Limit,
	}, nil
}

func (s *FeedingService) GetFeedingByID(id uint) (*dto.FeedingResponse, error) {
	feeding, err := s.feedingRepo.GetById(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrFeedingNotFound
		}
		return nil, err
	}

	// map ข้อมูลหมู
	var pigIDs []uint
	var pigCodes []string

	for _, item := range feeding.Items {
		pigIDs = append(pigIDs, item.PigID)
		if item.Pig.CodeName != "" {
			pigCodes = append(pigCodes, item.Pig.CodeName)
		}
	}

	resp := dto.FeedingResponse{
		ID:          feeding.ID,
		FoodID:      feeding.FoodID,
		FoodName:    feeding.FoodStock.FoodType.Name, // ได้จาก Preload FoodStock
		PigIDs:      pigIDs,
		PigCodeName: pigCodes,
		DateTime:    feeding.DateTime,
		Amount:      feeding.Amount,
		Note:        feeding.Note,
		CreatedName: feeding.Creator.FullName,
		CreatedRole: feeding.Creator.Role,
		UpdatedName: feeding.Updater.FullName,
		UpdatedRole: feeding.Updater.Role,
	}

	return &resp, nil

}
func (s *FeedingService) DeleteFeeding(id uint) error {
	err := s.feedingRepo.Delete(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFeedingNotFound
		}
		return err
	}
	return nil

}
