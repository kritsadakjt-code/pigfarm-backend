package models

import (
	"time"

	"gorm.io/gorm"
)

type Feeding struct {
	gorm.Model
	FoodID    uint          `gorm:"not null" json:"food_id"`
	DateTime  time.Time     `gorm:"not null" json:"date_time"` // วันเวลาที่ให้อาหาร
	Amount    float64       `gorm:"not null" json:"amount" `   // ปริมาณ (kg)
	Note      string        `gorm:"size:100; null" json:"note"`
	CreatedBy uint          `gorm:"not null" json:"created_by"`
	UpdatedBy uint          `gorm:"not null" json:"updated_by"`
	Items     []FeedingItem `gorm:"foreignKey:FeedingID;references:ID;" json:"items"`

	FoodStock FoodStock `gorm:"foreignKey:FoodID;references:ID" json:"food_stock"`
	Creator   User      `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	Updater   User      `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`
}

type FeedingItem struct {
	gorm.Model
	FeedingID uint `gorm:"not null" json:"feeding_id"` // FK ไป Feeding
	PigID     uint `gorm:"not null" json:"pig_id"`     // FK ไป Pig

	Pig Pig `gorm:"foreignKey:PigID;references:ID" json:"pig"`
}
