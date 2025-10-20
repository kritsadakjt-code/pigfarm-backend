package models

import (
	"time"

	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Date     time.Time `gorm:"not null" json:"date"`             // วันที่
	Category string    `gorm:"size:50;not null" json:"category"` // ประเภทค่าใช้จ่าย
	Amount   float64   `gorm:"not null" json:"amount"`           // จำนวนเงิน
	Note     string    `gorm:"size:100;null" json:"note"`        // รายละเอียดเพิ่มเติม
}
