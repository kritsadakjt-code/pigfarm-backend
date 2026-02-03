package usecases

import (
	"backend/dto"
	"backend/models"
	"math"
	"time"
)

type BreedingService struct {
	breedingRepo BreedingRepo
	pigRepo      PigRepository
}

func NewBreedingService(breedingRepo BreedingRepo, pigRepo PigRepository) *BreedingService {
	return &BreedingService{
		breedingRepo: breedingRepo,
		pigRepo:      pigRepo}
}

func (s *BreedingService) CreateBreeding(input dto.BreedingInput, userID uint) (*models.Breeding, error) {
	if input.FatherID == input.MotherID {
		return nil, ErrSamePig
	}
	father, err := s.pigRepo.GetByID(input.FatherID)
	if err != nil {
		return nil, ErrFatherNotFound
	}
	mother, err := s.pigRepo.GetByID(input.MotherID)
	if err != nil {
		return nil, ErrMotherNotFound
	}

	if father.Gender != "ผู้" || father.Type != "พ่อพันธุ์" {
		return nil, ErrInvalidFatherBreeder
	}
	if mother.Gender != "เมีย" || mother.Type != "เเม่พันธุ์" {
		return nil, ErrInvalidMotherBreeder
	}

	if father.Status != "พร้อมผสม" {
		return nil, ErrPigNotReady
	}
	if mother.Status != "พร้อมผสม" {
		return nil, ErrPigNotReady
	}
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	parsedDate, err := time.ParseInLocation("2006-01-02", input.BreedingDate, loc)
	if err != nil {
		return nil, ErrInvalidDate
	}
	if parsedDate.After(time.Now()) {
		return nil, ErrCantFuture
	}
	exist, err := s.breedingRepo.CheckBreedingAlready(input.FatherID, input.MotherID, parsedDate)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, ErrDuplicateBreeding
	}
	expectedBirth := parsedDate.AddDate(0, 0, 114)

	breeding := models.Breeding{
		FatherID:      input.FatherID,
		MotherID:      input.MotherID,
		BreedingDate:  parsedDate,
		ExpectedBirth: expectedBirth,
		Status:        "รอผล",
		Result:        "รอผล",
		Note:          input.Note,
		CreatedBy:     userID,
		UpdatedBy:     userID,
	}
	if err := s.breedingRepo.Create(&breeding); err != nil {
		return nil, err
	}
	return &breeding, nil
}

func (s *BreedingService) UpdateBreeding(req dto.BreedingUpdate, id uint, user_id uint) (*dto.BreedingResponse, error) {
	breeding, err := s.breedingRepo.GetByID(id)
	if err != nil {
		return nil, ErrBreedingNotFound
	}
	updates := make(map[string]interface{})
	var newStatusMother string

	if req.BreedingDate != nil {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			return nil, err
		}
		parsedDate, err := time.ParseInLocation("2006-01-02", *req.BreedingDate, loc)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if parsedDate.After(time.Now()) {
			return nil, ErrCantFuture
		}
		updates["breeding_date"] = parsedDate
		updates["expected_birth"] = parsedDate.AddDate(0, 0, 114)
	}
	if req.Status != nil {
		updates["status"] = *req.Status
		switch *req.Status {
		case "อุ้มท้อง":
			updates["result"] = "รอผล"
			newStatusMother = "อุ้มท้อง"
		case "ผสมไม่ติด":
			updates["result"] = "ไม่สําเร็จ"
			newStatusMother = "พร้อมผสม" // คืนสถานะให้พร้อมผสมใหม่
		case "เเท้ง":
			updates["result"] = "ไม่สําเร็จ"
			newStatusMother = "พักท้อง" // ต้องพักฟื้นก่อนผสมใหม่
		case "คลอดเเล้ว":
			updates["result"] = "สําเร็จ"
			newStatusMother = "ให้นมลูก"
		}
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}

	updates["created_by"] = user_id
	err = s.breedingRepo.UpdateBreeding(id, updates, breeding.MotherID, newStatusMother)
	if err != nil {
		return nil, err
	}

	updatedBreeding, err := s.breedingRepo.GetByID(id)
	if err != nil {
		return nil, ErrBreedingNotFound
	}

	return &dto.BreedingResponse{
		ID:             updatedBreeding.ID,
		FatherID:       updatedBreeding.FatherID,
		MotherID:       updatedBreeding.MotherID,
		FatherCodename: updatedBreeding.Father.CodeName,
		MotherCodename: updatedBreeding.Mother.CodeName,
		BreedingDate:   updatedBreeding.BreedingDate,
		ExpectedBirth:  updatedBreeding.ExpectedBirth,
		Status:         updatedBreeding.Status,
		Result:         updatedBreeding.Result,
		Note:           updatedBreeding.Note,
		CreatedName:    updatedBreeding.Creator.FullName,
		UpdatedName:    updatedBreeding.Updater.FullName,
	}, nil
}
func (s *BreedingService) GetAllBreeding(param dto.BreedingParam) (*dto.BreedingPaginationResponse, error) {
	if param.Page <= 0 {
		param.Page = 0
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}

	breedings, total, err := s.breedingRepo.GetAll(param)
	if err != nil {
		return nil, err
	}
	lastPage := int(math.Ceil(float64(total) / float64(param.Limit)))

	var resp []dto.BreedingResponse
	for _, b := range breedings {
		resp = append(resp, dto.BreedingResponse{
			ID:             b.ID,
			FatherID:       b.FatherID,
			MotherID:       b.MotherID,
			FatherCodename: b.Father.CodeName,
			MotherCodename: b.Mother.CodeName,
			BreedingDate:   b.BreedingDate,
			ExpectedBirth:  b.ExpectedBirth,
			Status:         b.Status,
			Result:         b.Result,
			Note:           b.Note,
			CreatedName:    b.Creator.FullName, // ระวัง nil pointer ถ้าไม่มีข้อมูล user
			CreatedRole:    b.Creator.Role,
			UpdatedName:    b.Updater.FullName,
			UpdatedRole:    b.Updater.Role,
		})
	}
	return &dto.BreedingPaginationResponse{
		Data:     resp,
		Total:    total,
		Page:     param.Page,
		LastPage: lastPage,
		Limit:    param.Limit,
	}, nil

}

func (s *BreedingService) DeleteBreeding(id uint) error {
	breeding, err := s.breedingRepo.GetByID(id)
	if err != nil {
		return ErrBreedingNotFound
	}
	reStatus := false
	if breeding.Mother.Status == "อุ้มท้อง" {
		reStatus = true
	}

	if err := s.breedingRepo.Delete(id, breeding.MotherID, reStatus); err != nil {
		return err
	}

	return nil
}
