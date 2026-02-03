package dto

import "time"

type FoodStockInput struct {
	FoodTypeID uint    `json:"food_type_id" validate:"required"`
	DateTime   string  `json:"date_time" validate:"required"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
}

type FoodStockUpdate struct {
	// Name     *string  `json:"name" validate:"omitempty,oneof=อาหารลูกหมู อาหารหมูขุน อาหารพ่อเเม่พันธุ์ อาหารหมูท้อง อาหารหมูให้นม วิตามินรวม เเร่ธาตุรวม โพรไบโอติก"`
	// Type     *string  `json:"type,omitempty" validate:"omitempty,oneof=อาหารหลัก อาหารเสริม"`
	DateTime *string  `json:"date_time,omitempty"`
	Amount   *float64 `json:"amount,omitempty"`
	Note     *string  `json:"note,omitempty"`
}

type FoodStockResponse struct {
	ID          uint      `json:"id"`
	FoodTypeID  uint      `json:"food_type_id"`
	FoodName    string    `json:"food_name"`
	FoodType    string    `json:"food_type"`
	DateTime    time.Time `json:"date_time"`
	Amount      float64   `json:"amount"`
	Note        string    `json:"note"`
	CreatedName string    `json:"created_name"`
	CreatedRole string    `json:"created_role"`
	UpdatedName string    `json:"updated_name"`
	UpdatedRole string    `json:"updated_role"`
}

type ParamFoodStock struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

type FoodStockPagiResp struct {
	Data     []FoodStockResponse `json:"data"`
	Total    int64               `json:"total"`
	Page     int                 `query:"page"`
	LastPage int                 `query:"last_page"`
	Limit    int                 `query:"limit"`
}
