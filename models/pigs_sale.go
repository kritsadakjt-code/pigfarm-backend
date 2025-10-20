package models

import (
	"time"

	"gorm.io/gorm"
)

type PigSale struct {
	gorm.Model
	SaleCode   string        `gorm:"size:100;not null;unique" json:"sale_code"`
	Date       time.Time     `gorm:"not null" json:"date"`
	Amount     int           `gorm:"not null" json:"amount"`
	TotalPrice float64       `gorm:"not null" json:"total_price"`
	Buyer      string        `gorm:"size:100;not null" json:"buyer"`
	Note       string        `gorm:"size:100;null" json:"note"`
	Items      []PigSaleItem `gorm:"foreignKey:PigSaleID;references:ID;constraint:OnDelete:CASCADE;" json:"items"`
	// Items      []PigSaleItem `gorm:"foreignKey:PigSaleID;references:ID;constraint:OnDelete:CASCADE" json:"items"`
}

type PigSaleItem struct {
	gorm.Model
	PigSaleID uint `gorm:"not null" json:"pig_sale_id"` //FK ไป PigSale
	PigID     uint `gorm:"not null" json:"pig_id"`      //FK ไป Pig

	Pig Pig `gorm:"foreignKey:PigID;references:ID" json:"pig"`
}
