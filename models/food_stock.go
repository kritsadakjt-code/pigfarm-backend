package models

import (
	"time"

	"gorm.io/gorm"
)

// type FoodStock struct {
// 	gorm.Model
// 	Name      string    `gorm:"size=100;not null; unique" json:"name"` // อาหารลูกหมู อาหารหมูขุน อาหารพ่อเเม่พันธุ์ อาหารหมูท้อง อาหารหมูให้นม วิตามินรวม เเร่ธาตุรวม โพรไบโอติก
// 	Type      string    `gorm:"size=50;not null" json:"type"`          // อาหารหลัก อาหารเสริม
// 	Amount    float64   `gorm:"not null" json:"amount"`
// 	DateTime  time.Time `gorm:"not null" json:"date_time"`
// 	Note      string    `gorm:"size=100;null" json:"note"`
// 	CreatedBy uint      `gorm:"not null" json:"created_by" `
// 	UpdatedBy uint      `gorm:"not null" json:"updated_by" `

// 	Creator User `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
// 	Updater User `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`
// }

type FoodStock struct {
	gorm.Model
	FoodTypeID uint      `gorm:"not null" json:"food_type_id"`
	Amount     float64   `gorm:"not null" json:"amount"`
	DateTime   time.Time `gorm:"not null" json:"date_time"`
	Note       string    `gorm:"size=100;null" json:"note"`
	CreatedBy  uint      `gorm:"not null" json:"created_by" `
	UpdatedBy  uint      `gorm:"not null" json:"updated_by" `

	FoodType FoodType `gorm:"foreignKey:FoodTypeID;references:ID" json:"food_type"`
	Creator  User     `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	Updater  User     `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`
}

func (f *FoodStock) IsLowStock(quantity float64) bool {
	return f.Amount <= quantity && f.Amount > 0
}

func (f *FoodStock) IsCriticalStock() bool {
	return f.Amount <= 5 && f.Amount > 0
}
