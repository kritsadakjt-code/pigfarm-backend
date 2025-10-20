package dto

import "time"

type ExpenseInput struct {
	Date     string  `json:"date" validate:"required"`
	Category string  `json:"category" validate:"required,oneof= ค่านํ้า ค่าไฟ ค่ายาวัคซีน ค่าอาหาร ค่าซ่อมบํารุง ค่าพนักงาน อื่นๆ"`
	Amount   float64 `json:"amount" `
	Note     string  `json:"note"`
}

type ExpenseUpdate struct {
	Date     *string  `json:"date,omitempty"`
	Category *string  `json:"category,omitempty" validate:"omitempty,oneof= ค่านํ้า ค่าไฟ ค่ายาวัคซีน ค่าอาหาร ค่าซ่อมบํารุง ค่าพนักงาน อื่นๆ"`
	Amount   *float64 `json:"amount,omitempty" `
	Note     *string  `json:"note,omitempty"`
}

type ExpenseResponse struct {
	ID       uint      `json:"id"`
	Date     time.Time `json:"date"`
	Category string    `json:"category"`
	Amount   float64   `json:"amount"`
	Note     string    `json:"note"`
}
