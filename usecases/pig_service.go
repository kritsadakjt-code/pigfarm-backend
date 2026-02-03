package usecases

import (
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"math"
	"time"
)

type PigService struct {
	pigRepo PigRepository
}

func NewPigService(pigRepo PigRepository) *PigService {
	return &PigService{pigRepo: pigRepo}
}

func validateCreateLogic(input dto.PigInput) error {
	if input.Weight <= 0 {
		return ErrInvalidWeight
	}

	if input.Gender == "ผู้" && input.Type == "เเม่พันธุ์" {
		return ErrMaleAsMother
	}

	if input.Gender == "เมีย" && input.Type == "พ่อพันธุ์" {
		return ErrFemaleAsFather
	}

	if input.Gender == "ผู้" &&
		(input.Status == "อุ้มท้อง" || input.Status == "ให้นมลูก") {
		return ErrMaleInvalidStatus
	}
	invalidStatus := map[string]bool{
		"อุ้มท้อง": true,
		"ให้นมลูก": true,
		"พร้อมผสม": true,
	}
	if (input.Type == "ลูกหมู" || input.Type == "หมูขุน") && invalidStatus[input.Status] {
		return fmt.Errorf("%w: %s", ErrPigTypeInvalidStatus, input.Status)
	}
	return nil

}

func (s *PigService) CreatePig(input dto.PigInput, userID uint) (*models.Pig, error) {
	if err := validateCreateLogic(input); err != nil {
		return nil, err
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	parsedDate, err := time.ParseInLocation("2006-01-02", input.BirthDate, loc)
	if err != nil {
		return nil, ErrInvalidDate
	}
	if parsedDate.After(time.Now()) {
		return nil, ErrCantFuture
	}

	now := time.Now()
	year := now.Format("06") // "06" คือ format สำหรับปี 2 ตัวท้าย (เช่น 25)
	_, week := now.ISOWeek()
	weekStr := fmt.Sprintf("%02d", week) // "%02d" คือ format ให้มี 0 นำหน้าถ้าเป็นเลขหลักเดียว

	// 3. ค้นหาลำดับเลขถัดไป (001-999)

	// สร้าง Pattern สำหรับค้นหา เช่น "2001-4225" (หาหมูพันธุ์ดูร็อกที่เกิดในสัปดาห์ที่ 42 ปี 25)
	pattern := fmt.Sprintf("%s%%-%s%s", input.CodePrefix, weekStr, year)
	lastPig, err := s.pigRepo.GenerateNextCode(pattern)

	nextNumber := 1
	if err == nil && lastPig != nil {
		var lastNum int
		// ดึงตัวเลขลำดับจากรหัสล่าสุด เช่น "D-005-42-25" จะได้เลข 5
		fmt.Sscanf(lastPig.CodeName, input.CodePrefix+"%d", &lastNum)
		nextNumber = lastNum + 1 // บวก 1 เพื่อเป็นเลขถัดไป
	}

	newCodeName := fmt.Sprintf("%s%03d-%s%s", input.CodePrefix, nextNumber, weekStr, year)
	// Set default status ถ้า user เลือก Type บางอย่างมา

	status := input.Status
	if input.Type == "ลูกหมู" {
		status = "กำลังเลี้ยง"
	} else if input.Type == "หมูขุน" {
		status = "กำลังขุน"
	}

	newPig := models.Pig{
		CodeName:  newCodeName,
		Name:      input.Name,
		Breed:     input.Breed,
		Gender:    input.Gender,
		Type:      input.Type,
		BirthDate: parsedDate,
		Weight:    input.Weight,
		Status:    status,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := s.pigRepo.Create(&newPig); err != nil {
		return nil, err
	}
	return &newPig, nil
}

func validateUpdateLogic(req dto.PigUpdate) error {
	// เช็ค Gender vs Type
	if req.Gender != nil && req.Type != nil {
		if *req.Gender == "ผู้" && *req.Type == "เเม่พันธุ์" {
			return ErrMaleAsMother
		}
		if *req.Gender == "เมีย" && *req.Type == "พ่อพันธุ์" {
			return ErrFemaleAsFather
		}
	}

	// เช็ค Gender vs Status
	if req.Gender != nil && req.Status != nil {
		if *req.Gender == "ผู้" {
			if *req.Status == "อุ้มท้อง" || *req.Status == "ให้นมลูก" {
				return ErrMaleInvalidStatus
			}
		}
	}

	// เช็ค Type vs Status
	if req.Type != nil && req.Status != nil {
		invalidStatus := map[string]bool{"อุ้มท้อง": true, "ให้นมลูก": true, "พร้อมผสม": true}
		if (*req.Type == "ลูกหมู" || *req.Type == "หมูขุน") && invalidStatus[*req.Status] {
			return fmt.Errorf("%w: %s", ErrPigTypeInvalidStatus, *req.Status)
		}
	}

	// เช็ค Weight
	if req.Weight != nil && *req.Weight <= 0 {
		return ErrInvalidWeight
	}

	return nil
}
func (s *PigService) UpdatePig(id uint, req dto.PigUpdate, user_id uint) (*models.Pig, error) {
	pig, err := s.pigRepo.GetByID(id)
	if err != nil {
		return nil, ErrPigNotFound
	}
	if err := validateUpdateLogic(req); err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Breed != nil {
		updates["breed"] = *req.Breed
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.Weight != nil {
		updates["weight"] = *req.Weight
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if req.BirthDate != nil {
		loc, err := time.LoadLocation("Asia/Bangkok")
		if err != nil {
			return nil, errors.New("failed to load time zone")
		}
		parsedDate, err := time.ParseInLocation("2006-01-02", *req.BirthDate, loc)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if parsedDate.After(time.Now()) {
			return nil, ErrCantFuture
		}
		updates["birth_date"] = parsedDate
	}
	updates["updated_by"] = user_id

	if err := s.pigRepo.Update(id, updates); err != nil {
		return nil, ErrFailedUpdate
	}

	return pig, nil

}

func (s *PigService) FindAllPagination(param dto.PigParam) (*dto.PigPaginationResponse, error) {
	// set Defualt
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	pigs, total, err := s.pigRepo.FindAllPagination(param)
	if err != nil {
		return nil, err
	}

	var resp []dto.PigResponse
	for _, p := range pigs {
		resp = append(resp, dto.PigResponse{
			ID:        p.ID,
			CodeName:  p.CodeName,
			Name:      p.Name,
			Breed:     p.Breed,
			Gender:    p.Gender,
			Type:      p.Type,
			BirthDate: p.BirthDate,
			Weight:    p.Weight,
			Status:    p.Status,
			// เช็ค nil เผื่อบางเคสไม่มี User (เช่น Seed มา)
			CreatedName: p.Creator.FullName,
			UpdatedName: p.Updater.FullName,
		})
	}

	// คํานวณหน้าสุดท้าย
	lastPage := int(math.Ceil(float64(total) / float64(param.Limit)))

	return &dto.PigPaginationResponse{
		Data:     resp,
		Total:    total,
		Page:     param.Page,
		Limit:    param.Limit,
		LastPage: lastPage,
	}, nil
}

func (s *PigService) GetPigByID(id uint) (*models.Pig, error) {
	pig, err := s.pigRepo.GetByID(id)
	if err != nil {
		return nil, ErrPigNotFound
	}
	return pig, nil
}

func (s *PigService) DeletePig(id uint) error {
	_, err := s.pigRepo.GetByID(id)
	if err != nil {
		return ErrPigNotFound
	}
	isUsed, err := s.pigRepo.IsUsedInBreeding(id)
	if err != nil {
		return err
	}
	if isUsed {
		return ErrIsUsedInBreeding
	}
	if err := s.pigRepo.Delete(id); err != nil {
		return err
	}
	return nil
}
