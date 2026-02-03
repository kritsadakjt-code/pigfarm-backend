package models

import (
	"time"

	"gorm.io/gorm"
)

type Breeding struct {
	gorm.Model
	FatherID      uint      `json:"father_id" gorm:"not null"`      // FK -> pigs.id
	MotherID      uint      `json:"mother_id" gorm:"not null"`      // FK -> pigs.id
	BreedingDate  time.Time `json:"breeding_date" gorm:"not null"`  // วันผสมพันธุ์
	ExpectedBirth time.Time `json:"expected_birth" gorm:"not null"` // วันคาดการณ์คลอด
	Status        string    `json:"status" gorm:"size:20;not null"` // รอผล อุ้มท้อง เเท้ง ผสมไม่ติด
	Result        string    `json:"result" gorm:"size:20;not null"` // สําเร็จ ไม่สําเร็จ
	Note          string    `json:"note" gorm:"size:255"`
	CreatedBy     uint      `json:"created_by" gorm:"not null"`
	UpdatedBy     uint      `json:"updated_by" gorm:"not null"`

	// Relations
	Father Pig `gorm:"foreignKey:FatherID;references:ID" json:"father"`
	Mother Pig `gorm:"foreignKey:MotherID;references:ID" json:"mother"`

	Creator User `gorm:"foreignKey:CreatedBy;references:ID" json:"creator"`
	Updater User `gorm:"foreignKey:UpdatedBy;references:ID" json:"updater"`
}

// func (b *Breeding) DaysUntilBirth() int {
// 	return int(time.Until(b.ExpectedBirth).Hours() / 24)
// }
// func (b *Breeding) IsPregnant() bool {
// 	return b.Status == "อุ้มท้อง"
// }
