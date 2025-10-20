package dto

import "time"

type FoodStockInput struct {
	Name     string  `json:"name" validate:"required,oneof=อาหารลูกหมู อาหารหมูขุน อาหารพ่อเเม่พันธุ์ อาหารหมูท้อง อาหารหมูให้นม วิตามินรวม เเร่ธาตุรวม โพรไบโอติก"`
	Type     string  `json:"type" validate:"required,oneof=อาหารหลัก อาหารเสริม"`
	DateTime string  `json:"date_time" validate:"required"`
	Amount   float64 `json:"amount"`
	Note     string  `json:"note"`
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
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	DateTime    time.Time `json:"date_time"`
	Amount      float64   `json:"amount"`
	Note        string    `json:"note"`
	CreatedName string    `json:"created_name"`
	CreatedRole string    `json:"created_role"`
	UpdatedName string    `json:"updated_name"`
	UpdatedRole string    `json:"updated_role"`
}
