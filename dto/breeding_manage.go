package dto

import "time"

type BreedingInput struct {
	FatherID     uint   `json:"father_id" validate:"required"`
	MotherID     uint   `json:"mother_id" validate:"required"`
	BreedingDate string `json:"breeding_date" validate:"required"`
	// Status       string `json:"status" validate:"required,oneof=รอผล"`
	// Result       string `json:"result" validate:"required,oneof=รอผล"`
	Note string `json:"note,omitempty"`
}

type BreedingUpdate struct {
	// *time.Time เเบบนี้รับเวลาจากการเลือก ปฏิทินไม่ได้ต้องเป็น string คิดว่า
	BreedingDate *string `json:"breeding_date,omitempty"`
	Status       *string `json:"status" validate:"omitempty,oneof=อุ้มท้อง ผสมไม่ติด เเท้ง คลอดเเล้ว"`
	// result update auto
	// Result *string `json:"result" validate:"omitempty,oneof=สําเร็จ ไม่สําเร็จ"`
	Note *string `json:"note,omitempty"`
}

type BreedingResponse struct {
	ID             uint      `json:"id"`
	FatherID       uint      `json:"father_id"`
	MotherID       uint      `json:"mother_id"`
	FatherCodename string    `json:"father_codename"`
	MotherCodename string    `json:"mother_codename"`
	BreedingDate   time.Time `json:"breeding_date"`
	ExpectedBirth  time.Time `json:"expected_birth"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	Note           string    `json:"note"`
	CreatedName    string    `json:"created_name"`
	CreatedRole    string    `json:"created_role"`
	UpdatedName    string    `json:"updated_name"`
	UpdatedRole    string    `json:"updated_role"`
}
