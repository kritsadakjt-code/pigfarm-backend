// package controllers

// import (
// 	"backend/config"
// 	"backend/models"
// 	"backend/utils"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// )

// func RequestChangeEmail(c *fiber.Ctx) error {
// 	userIDStr := fmt.Sprintf("%v", c.Locals("user_id"))
// 	userID, err := strconv.ParseUint(userIDStr, 10, 64)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user_id"})
// 	}

// 	var input struct {
// 		NewEmail string `json:"new_email"`
// 	}
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
// 	}

// 	var user models.User
// 	if err := config.DB.First(&user, "id = ?", userID).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
// 	}

// 	// check email ซ้ำ
// 	var count int64
// 	if err := config.DB.Model(&models.User{}).Where("email = ?", input.NewEmail).Count(&count).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to check email"})
// 	}
// 	if count > 0 {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email already exists"})
// 	}

// 	// generate token
// 	token := uuid.NewString()
// 	exp := time.Now().Add(30 * time.Minute)

// 	user.PendingEmail = input.NewEmail
// 	user.EmailVerificationToken = token
// 	user.EmailVerificationExpiry = &exp

// 	if err := config.DB.Save(&user).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
// 	}

// 	frontendURL := os.Getenv("FRONTEND_URL")
// 	verificationURL := fmt.Sprintf("%s/verify-email/change/%s", frontendURL, token)

// 	if err := utils.SendEmailVerification(input.NewEmail, verificationURL); err != nil {
// 		fmt.Println("Error sending email:", err)
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send email"})
// 	}

// 	return c.JSON(fiber.Map{"message": "Verification email sent"})
// }

// func VerifyEmailForChange(c *fiber.Ctx) error {
// 	token := c.Params("token")
// 	if token == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token missing"})
// 	}

// 	var user models.User
// 	if err := config.DB.Where("email_verification_token = ?", token).First(&user).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "invalid token"})
// 	}

// 	if user.EmailVerificationExpiry == nil || user.EmailVerificationExpiry.Before(time.Now()) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
// 	}

// 	user.Email = user.PendingEmail
// 	user.PendingEmail = ""
// 	user.EmailVerificationToken = ""
// 	user.EmailVerificationExpiry = nil

// 	if err := config.DB.Save(&user).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update email"})
// 	}

// 	return c.JSON(fiber.Map{"message": "Email updated successfully"})
// }

// func VerifyEmailForRegister(c *fiber.Ctx) error {
// 	token := c.Params("token")
// 	if token == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token missing"})
// 	}

// 	var user models.User
// 	if err := config.DB.Where("email_verification_token = ? AND email_verified_register IS NULL", token).First(&user).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid or already used token"})
// 	}

// 	if user.EmailVerificationExpiry == nil || user.EmailVerificationExpiry.Before(time.Now()) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "token expired"})
// 	}

// 	now := time.Now()
// 	user.EmailVerifiedRegister = &now
// 	user.Status = "pending"

// 	user.EmailVerificationToken = ""
// 	user.EmailVerificationExpiry = nil

// 	if err := config.DB.Save(&user).Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify email"})
// 	}

// 	frontendURL := os.Getenv("FRONTEND_URL")
// 	return c.Redirect(fmt.Sprintf("%s/email-verified-success", frontendURL))
// }

package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func RequestChangeEmail(c *fiber.Ctx) error {
	user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	user_id, err := strconv.ParseUint(user_id_str, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user_id"})
	}

	var input struct {
		NewEmail string `json:"new_email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	var user models.User
	if err := config.DB.First(&user, "id = ?", user_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}

	// check ว่า email ซํ้าหรือไม่
	var count int64
	if err := config.DB.Model(&user).Where("email = ?", input.NewEmail).Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to check email"})
	}
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "email already exists"})
	}

	// generate token for reset email
	token := uuid.NewString()
	exp := time.Now().Add(30 * time.Minute)
	user.PendingEmail = input.NewEmail
	user.EmailVerificationToken = token
	user.EmailVerificationExpiry = &exp
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update"})
	}
	frontendURL := os.Getenv("FRONTEND_URL_DEV")

	// ในฟังก์ชัน RequestChangeEmail
	verificationURL := fmt.Sprintf("%s/verify-email/change/%s", frontendURL, token)
	err = utils.SendEmailVerification(input.NewEmail, verificationURL)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send email"})
	}
	return c.JSON(fiber.Map{"message": "verification email sent"})
}

func VerifyEmailForChange(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "token missing"})
	}
	var user models.User
	err := config.DB.Where("email_verification_token = ?", token).First(&user).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "invalid token"})
	}

	if user.EmailVerificationExpiry == nil || user.EmailVerificationExpiry.Before(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "token expired"})
	}

	// update email
	user.Email = user.PendingEmail
	user.PendingEmail = ""
	user.EmailVerificationToken = ""
	user.EmailVerificationExpiry = nil
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update"})
	}
	return c.JSON(fiber.Map{"message": "email updated successfully"})
}

func VerifyEmailForRegister(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Verification token is missing"})
	}

	var user models.User
	// ค้นหา user จาก token และต้องแน่ใจว่ายังไม่เคยถูกยืนยันมาก่อน
	err := config.DB.Where("email_verification_token = ? AND email_verified_register IS NULL", token).First(&user).Error
	if err != nil {
		// ถ้าไม่เจอ อาจเป็นเพราะ token ผิด หรือถูกใช้ไปแล้ว
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invalid or already used token"})
	}

	if user.EmailVerificationExpiry == nil || user.EmailVerificationExpiry.Before(time.Now()) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token has expired"})
	}

	now := time.Now()
	user.EmailVerifiedRegister = &now // 1. ประทับเวลาว่าอีเมลถูกยืนยันแล้ว
	user.Status = "pending"

	user.EmailVerificationToken = ""
	user.EmailVerificationExpiry = nil

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify email"})
	}

	// แนะนำให้ Redirect ผู้ใช้ไปหน้า "ยืนยันสำเร็จ" บน Frontend
	// frontendURL := os.Getenv("FRONTEND_URL")
	// return c.Redirect(fmt.Sprintf("%s/email-verified-success", frontendURL))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email verified successfully. Please wait for admin approval."})
}
