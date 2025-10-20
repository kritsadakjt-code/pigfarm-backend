package dto

import "time"

type PigSaleInput struct {
	PigIDs []uint `json:"pig_ids" validate:"required,dive,gt=0"`
	Date   string `json:"date" validate:"required"`
	// Amount     int       `json:"amount" validate:"required,gt=0"`
	TotalPrice float64 `json:"total_price" validate:"required,gt=0"`
	Buyer      string  `json:"buyer" validate:"required"`
	Note       string  `json:"note"`
}

type PigSaleUpdate struct {
	PigIDs *[]uint `json:"pig_ids" validate:"omitempty,dive,gt=0"`
	Date   *string `json:"date"`
	// Amount     *int       `json:"amount" validate:"omitempty,gt=0"`
	TotalPrice *float64 `json:"total_price" validate:"omitempty,gt=0"`
	Buyer      *string  `json:"buyer"`
	Note       *string  `json:"note"`
}

type PigSaleResponse struct {
	ID          uint      `json:"id"`
	SaleCode    string    `json:"sale_code"`
	PigIDs      []uint    `json:"pig_ids"`
	PigCodeName []string  `json:"pig_code_name"`
	Date        time.Time `json:"date"`
	Amount      int       `json:"amount"`
	TotalPrice  float64   `json:"total_price"`
	Buyer       string    `json:"buyer"`
	Note        string    `json:"note,omitempty"`
}
