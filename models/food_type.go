package models

import "gorm.io/gorm"

type FoodType struct {
	gorm.Model
	Name string `gorm:"size:50;not null;unique" json:"name"`
	Type string `gorm:"size:50;not null" json:"type"`
}
