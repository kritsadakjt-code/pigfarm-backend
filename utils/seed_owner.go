package utils

import (
	"backend/config"
	"backend/models"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedOwnerUser(db *gorm.DB) {
	var existingOwner models.User
	err := config.DB.Where("email = ? AND role = ?", "owner@owner.com", "owner").First(&existingOwner).Error
	if err == nil {
		log.Println("Owner user already exists, skipping creation")
		return
	}
	// สร้างรหัสผ่าน default
	password := "12345678"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	owner := models.User{
		FullName: "Owner Admin",
		Email:    "owner@owner.com",
		Password: string(hashed),
		Phone:    "0999999999",
		Role:     "owner",
		Status:   "approved",
	}

	if err := db.Create(&owner).Error; err != nil {
		log.Fatal("Failed to create owner user:", err)
	}
	log.Println("Owner user created successfully")
}
