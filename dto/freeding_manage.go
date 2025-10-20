package dto

import "time"

type FeedingInput struct {
	FoodID   uint    `json:"food_id" validate:"required"`
	DateTime string  `json:"date_time" validate:"required"`
	Amount   float64 `json:"amount"` // kg
	Note     string  `json:"note"`
}

type FeedingUpdate struct {
	FoodID   *uint    `json:"food_id,omitempty"`
	DateTime *string  `json:"date_time,omitempty"`
	Amount   *float64 `json:"amount,omitempty"` // kg
	Note     *string  `json:"note,omitempty"`   // หมายเหตุ

}

type FeedingResponse struct {
	ID          uint      `json:"id"`
	FoodID      uint      `json:"food_id"`
	FoodName    string    `json:"food_name"`
	DateTime    time.Time `json:"date_time"`
	Amount      float64   `json:"amount"`
	Note        string    `json:"note"`
	CreatedName string    `json:"created_name"`
	CreatedRole string    `json:"created_role"`
	UpdatedName string    `json:"updated_name"`
	UpdatedRole string    `json:"updated_role"`
}
