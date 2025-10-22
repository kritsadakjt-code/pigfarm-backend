package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(c *fiber.Ctx) error {
	var input dto.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	input.Email = strings.ToLower(input.Email)
	existingEmail := &models.User{}
	err := config.DB.Where("email = ?", input.Email).First(existingEmail).Error
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already exists"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}
	verificationToken := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(30 * time.Minute)
	user := models.User{
		FullName:                input.FullName,
		Email:                   strings.ToLower(input.Email),
		Password:                string(hashed),
		Phone:                   input.Phone,
		Role:                    input.Role,
		Status:                  "pending_verification_email",
		EmailVerificationToken:  verificationToken,
		EmailVerificationExpiry: &expiry,
	}

	if err = config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to register"})
	}

	// frontendURL := os.Getenv("FRONTEND_URL_DEV")
	frontendURL := os.Getenv("FRONTEND_URL")
	verificationURL := fmt.Sprintf("%s/verify-email/register/%s", frontendURL, verificationToken)

	if err := utils.SendEmailVerification(input.Email, verificationURL); err != nil {
		fmt.Println("Error sending verification email:", err)
		// Don't fail the whole request, but log the error
	}

	return c.Status(200).JSON(fiber.Map{"message": "Registration successful", "status": user.Status})

}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	// check email
	input.Email = strings.ToLower(input.Email)

	user := &models.User{}
	err := config.DB.Where("email = ?", input.Email).First(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Email not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// compare hash password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Incorrect password"})
	}

	// --- ส่วนที่แก้ไข ---
	// ตรวจสอบสถานะต่างๆ ก่อนที่จะเช็ค approved
	fmt.Println(user.Status)
	switch user.Status {
	case "pending_verification_email":
		return c.Status(403).JSON(fiber.Map{"error": "กรุณายืนยันอีเมลก่อน"})
	case "pending":
		return c.Status(403).JSON(fiber.Map{"error": "Account not approved yet"})
	case "rejected":
		return c.Status(403).JSON(fiber.Map{"error": "Your account was rejected"})
	case "approved":
		// ถ้าเป็น approved ให้ผ่านไปทำขั้นตอนต่อไป (gen token)
		break
	}

	// gen token (จะทำเฉพาะกรณี status เป็น "approved")
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// return
	return c.JSON(fiber.Map{
		"token":    token,
		"fullName": user.FullName,
		"email":    user.Email,
		"role":     user.Role,
	})
}
