package controllers

import (
	"backend/config"
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func ForgetPassword(c *fiber.Ctx) error {
	type Request struct {
		Email string `json:"email"`
	}
	var input Request
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// check email
	user := &models.User{}

	err := config.DB.Where("email = ?", input.Email).First(user).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// ถ้าเป็น error อื่นที่ไม่ใช่ "ไม่เจอ" ก็ควร log ไว้
			fmt.Println("Database error:", err)
		}
		return c.Status(200).JSON(fiber.Map{"message": "หากมีบัญชีที่ใช้อีเมลนี้อยู่ ระบบได้ส่งลิงก์สำหรับรีเซ็ตรหัสผ่านไปให้แล้ว"})
	}

	//gen token ใหม่
	token := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(30 * time.Minute)
	user.ResetToken = token
	user.ResetTokenExpiry = &expiry

	if err := config.DB.Save(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update user"})
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	resetURL := fmt.Sprintf("%s/reset-password/%s", frontendURL, token)

	// ส่ง link reset to email
	err = utils.SentResetPasswordFromEmail(user.Email, resetURL)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send email"})
	}

	return c.JSON(fiber.Map{"message": "หากมีบัญชีที่ใช้อีเมลนี้อยู่ ระบบได้ส่งลิงก์สำหรับรีเซ็ตรหัสผ่านไปให้แล้ว"})
}

// // reset from email
func ResetPassword(c *fiber.Ctx) error {
	// ดึง token จาก URL
	token := c.Params("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing token"})
	}
	// ไว้รับ new password ที่ส่งเข้ามา json
	type Input struct {
		NewPassword string `json:"new_password"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	// check user with token from URL
	user := &models.User{}
	err := config.DB.Where("reset_token = ? AND reset_token_expiry > ?", token, time.Now()).First(user).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token is invalid or has expired"})
	}
	// update new password
	hash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to hash password"})
	}
	user.Password = string(hash)
	// ล้าง token เพื่อให้ใช้ได้เเค่ครั้งเดียว
	user.ResetToken = ""
	user.ResetTokenExpiry = nil
	// update
	if err := config.DB.Save(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update password"})
	}
	return c.JSON(fiber.Map{"message": "Password reset successfully"})
}

// change password จาก ระบบ
func ChangePassword(c *fiber.Ctx) error {
	// // ถ้าจะเอาไปใช้ต่อ เช่น บันทึก user_id ลงใน db
	// user_id_str := fmt.Sprintf("%v", c.Locals("user_id"))
	// // เเปลงเป็น uint64
	// user_id64, err := strconv.ParseUint(user_id_str, 10, 64)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	// }
	// user_id := uint(user_id64)

	user_id := c.Locals("user_id")

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	var user models.User
	if err := config.DB.Where("id = ?", user_id).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "user not found"})
	}
	// compare current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "current password incorrect"})
	}
	// gen newPassword
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to generate password"})
	}
	user.Password = string(hashed)
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to update password"})
	}

	return c.JSON(fiber.Map{"message": "password update successful"})
}
