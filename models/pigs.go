package models

import (
	"time"

	"gorm.io/gorm"
)

//	type Pig struct {
//		gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt
//		Name       string    `json:"name" gorm:"type:varchar(100);not null"`
//		Breed      string    `json:"breed" gorm:"type:varchar(50);not null"`  // ลาร์จไวท์, แลนด์เรซ, ดูร็อก
//		Gender     string    `json:"gender" gorm:"type:varchar(10);not null"` // ผู้, เมีย
//		Type       string    `json:"type" gorm:"type:varchar(20);not null"`   // พ่อพันธุ์, แม่พันธุ์, หมูขุน, ลูกหมู
//		BirthDate  time.Time `json:"birth_date" gorm:"not null"`              // วันเกิด
//		Weight     float64   `json:"weight"`                                  // น้ําหนัก (kg)
//		Status     string    `json:"status" gorm:"size:20;not null"`          // อุ้มท้อง, พร้อมผสม, ให้นมลูก, กำลังขุน, กำลังเลี้ยง, ขายเเล้ว
//	}

type Pig struct {
	gorm.Model
	CodeName  string    `gorm:"size:100;unique; not null" json:"code_name"`
	Name      string    `gorm:"size:50" json:"name"`
	Breed     string    `gorm:"size:50;not null" json:"breed"`  // ลาร์จไวท์, แลนด์เรซ, ดูร็อก
	Gender    string    `gorm:"size:10;not null" json:"gender"` // ผู้, เมีย
	Type      string    `gorm:"size:20;not null" json:"type"`   // พ่อพันธุ์, แม่พันธุ์, หมูขุน, ลูกหมู
	BirthDate time.Time `gorm:"not null" json:"birth_date"`     // วันเกิด
	Weight    float64   `gorm:"not null" json:"weight"`         // น้ำหนัก (kg)
	Status    string    `gorm:"size:20;not null" json:"status"` // อุ้มท้อง, พร้อมผสม, ให้นมลูก, กำลังขุน, กำลังเลี้ยง, ขายแล้ว
	CreatedBy uint      `json:"created_by" gorm:"not null"`
	UpdatedBy uint      `json:"updated_by" gorm:"not null"`

	Creator User `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	Updater User `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`
}
