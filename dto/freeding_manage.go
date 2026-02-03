package dto

import "time"

type FeedingInput struct {
	FoodID   uint    `json:"food_id" validate:"required"`
	PigIDs   []uint  `json:"pig_ids" validate:"omitempty,dive,gt=0"`
	DateTime string  `json:"date_time" validate:"required"`
	Amount   float64 `json:"amount"` // kg
	Note     string  `json:"note"`
}

type FeedingUpdate struct {
	FoodID   *uint    `json:"food_id,omitempty"`
	PigIDs   *[]uint  `json:"pig_ids" validate:"omitempty,dive,gt=0"`
	DateTime *string  `json:"date_time,omitempty"`
	Amount   *float64 `json:"amount,omitempty"` // kg
	Note     *string  `json:"note,omitempty"`   // หมายเหตุ

}

type FeedingResponse struct {
	ID          uint      `json:"id"`
	FoodID      uint      `json:"food_id"`
	FoodName    string    `json:"food_name"`
	PigIDs      []uint    `json:"pig_ids"` // <-- เพิ่ม
	PigCodeName []string  `json:"pig_code_name"`
	DateTime    time.Time `json:"date_time"`
	Amount      float64   `json:"amount"`
	Note        string    `json:"note"`
	CreatedName string    `json:"created_name"`
	CreatedRole string    `json:"created_role"`
	UpdatedName string    `json:"updated_name"`
	UpdatedRole string    `json:"updated_role"`
}

type ParamFeeding struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
	FoodID uint   `query:"food_id"` // กรองดูเฉพาะการให้อาหารชนิดนี้
}

type FeedingPagiResp struct {
	Data     []FeedingResponse `json:"data"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
	Total    int64             `json:"total"`
	LastPage int               `json:"last_page"`
}
