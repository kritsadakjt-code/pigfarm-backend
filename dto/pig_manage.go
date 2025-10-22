package dto

import "time"

type PigInput struct {
	CodePrefix string  `json:"code_prefix" validate:"required,len=1"`
	Name       string  `json:"name"`
	Breed      string  `json:"breed" validate:"required,oneof=ลาร์จไวท์ แลนด์เรซ ดูร็อก"`
	Gender     string  `json:"gender" validate:"required,oneof=ผู้ เมีย"`
	Type       string  `json:"type" validate:"required,oneof=พ่อพันธุ์ เเม่พันธุ์ หมูขุน ลูกหมู"`
	BirthDate  string  `json:"birth_date" `
	Weight     float64 `json:"weight" `
	Status     string  `json:"status" validate:"required,oneof=อุ้มท้อง พร้อมผสม ให้นมลูก กำลังขุน กำลังเลี้ยง ขายเเล้ว พร้อมขาย"`
}

type PigResponse struct {
	ID          uint      `json:"id"`
	CodeName    string    `json:"code_name"`
	Name        string    `json:"name"`
	Breed       string    `json:"breed"`
	Gender      string    `json:"gender"`
	Type        string    `json:"type"`
	BirthDate   time.Time `json:"birth_date"`
	Weight      float64   `json:"weight"`
	Status      string    `json:"status"`
	CreatedName string    `json:"created_name"`
	CreatedRole string    `json:"created_role"`
	UpdatedName string    `json:"updated_name"`
	UpdatedRole string    `json:"updated_role"`
}

type PigUpdate struct {
	CodeName  *string  `json:"code_name"`
	Name      *string  `json:"name"`
	Breed     *string  `json:"breed" validate:"omitempty,oneof=ลาร์จไวท์ แลนด์เรซ ดูร็อก"`
	Gender    *string  `json:"gender" validate:"omitempty,oneof=ผู้ เมีย"`
	Type      *string  `json:"type" validate:"omitempty,oneof=พ่อพันธุ์ เเม่พันธุ์ หมูขุน ลูกหมู"`
	BirthDate *string  `json:"birth_date,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
	Status    *string  `json:"status" validate:"omitempty,oneof=อุ้มท้อง พร้อมผสม ให้นมลูก กำลังขุน กำลังเลี้ยง ขายเเล้ว พร้อมขาย ตายเเล้ว"`
}
