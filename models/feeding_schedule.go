package models

import "gorm.io/gorm"

//
type FeedingSchedule struct {
	gorm.Model
	Name          string `gorm:"size:100;not null" json:"name"`
	ScheduledTime string `gorm:"size:5;not null;index" json:"scheduled_time"`
	IsActive      bool   `gorm:"default:true" json:"is_active"`
	Note          string `gorm:"size:255" json:"note"`
	CreatedBy     uint   `gorm:"not null" json:"created_by"`

	// Relation: หนึ่ง Schedule มีได้หลาย Items
	Items   []FeedingScheduleItem `gorm:"foreignKey:ScheduleID" json:"items"`
	Creator User                  `gorm:"foreignKey:CreatedBy"`
}

// ตารางย่อยเก็บรายการอาหาร
type FeedingScheduleItem struct {
	gorm.Model
	ScheduleID uint    `gorm:"not null" json:"schedule_id"` // FK ไปยัง feeding_schedules
	FoodID     uint    `gorm:"not null" json:"food_id"`     // FK ไปยัง food_stocks
	Amount     float64 `gorm:"not null" json:"amount"`

	FoodStock FoodStock `gorm:"foreignKey:FoodID" json:"food_stock"`
}
