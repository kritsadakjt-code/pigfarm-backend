package models

import (
	"time"

	"gorm.io/gorm"
)

type Health struct {
	gorm.Model
	PigID     uint      `gorm:"not null" json:"pig_id" `       // อ้างอิงหมู
	Date      time.Time `gorm:"not null" json:"date" `         // วันที่บันทึก
	Type      string    `gorm:"size=50;not null" json:"type" ` // ประเภทเช่น วัคซีน, ตรวจสุขภาพ
	Detail    string    `gorm:"size=100; null" json:"detail" ` // รายละเอียด
	Note      string    `gorm:"size=100; null" json:"note" `
	CreatedBy uint      `gorm:"not null" json:"created_by" `
	UpdatedBy uint      `gorm:"not null" json:"updated_by" `

	Pig     Pig  `gorm:"foreignKey:PigID;references:ID" json:"pig"`
	Creator User `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	Updater User `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`
}
