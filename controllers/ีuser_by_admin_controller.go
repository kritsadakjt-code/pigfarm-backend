package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// search users
func SearchUser(c *fiber.Ctx) error {

	keyword := c.Query("keyword")

	var users []models.User
	// ชื่อที่ต้องการค้นหา
	if keyword != "" {
		// search
		err := config.DB.Where("full_name ILIKE ?", "%"+keyword+"%").Find(&users).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "cannot fetch users"})
		}
	} else {
		err := config.DB.Find(&users).Error
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "cannot fetch users"})
		}
	}

	// สร้าง slice เพื่อ map users
	var responses []dto.UserResponse
	for _, user := range users {
		responses = append(responses, dto.UserResponse{
			ID:       fmt.Sprintf("%d", user.ID),
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
			Role:     user.Role,
			Status:   user.Status,
		})
	}

	return c.JSON(responses)

}

func GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var user models.User
	err = config.DB.First(&user, "id = ?", user_id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Not Found User"})
	}

	resp := dto.UserResponse{
		ID:        fmt.Sprintf("%d", user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}
	return c.JSON(resp)
}

func GetAllUsers(c *fiber.Ctx) error {
	var users []models.User
	err := config.DB.Order("created_at DESC").Find(&users).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Failed to fetch users"})
	}
	var resp []dto.UserResponse
	for _, user := range users {
		resp = append(resp, dto.UserResponse{
			ID:        fmt.Sprintf("%d", user.ID),
			FullName:  user.FullName,
			Email:     user.Email,
			Phone:     user.Phone,
			Role:      user.Role,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
		})
	}

	return c.JSON(resp)
}

func CreateUser(c *fiber.Ctx) error {
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
	// input.Password = "12345678"
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to hash password"})
	}
	user := models.User{
		FullName: input.FullName,
		Email:    strings.ToLower(input.Email),
		Password: string(hashed),
		Phone:    input.Phone,
		Role:     input.Role,
		Status:   input.Status,
	}

	if err = config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to create user"})
	}

	return c.Status(201).JSON(fiber.Map{"message": "create user successful"})

}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	user := &models.User{}
	err = config.DB.First(user, "id = ?", user_id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var input struct {
		FullName *string `json:"full_name,omitempty"`
		Email    *string `json:"email,omitempty"`
		Phone    *string `json:"phone,omitempty"`
		Role     *string `json:"role" `
		Status   *string `json:"status" `
	}

	err = c.BodyParser(&input)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := make(map[string]interface{})
	if input.FullName != nil {
		updates["full_name"] = input.FullName
	}
	if input.Email != nil {
		updates["email"] = input.Email
	}
	if input.Phone != nil {
		updates["phone"] = input.Phone
	}
	if input.Role != nil {
		updates["role"] = input.Role
	}
	if input.Status != nil {
		updates["status"] = input.Status
	}

	err = config.DB.Model(user).Updates(updates).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}
	// response ด้วยข้อมูล user ที่ update เเล้ว
	resp := dto.UserResponse{
		ID:        fmt.Sprintf("%d", user.ID),
		FullName:  user.FullName,
		Email:     user.Email,
		Phone:     user.Phone,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
	}

	return c.JSON(resp)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	// check user ก่อนว่ามีไหม
	user := &models.User{}
	err = config.DB.First(user, "id = ?", user_id).Error
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	// config.DB.Unscoped().Delete(user).Error ลบจริงๆ
	err = config.DB.Unscoped().Delete(user).Error
	if err != nil {

		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{"message": "User deleted successfully"})
}

func ApproveUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	// check
	user := &models.User{}
	if err := config.DB.First(user, user_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User Not Found"})
	}

	// ถ้าเจอ
	user.Status = "approved"
	if err := config.DB.Save(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Approval failed"})
	}

	return c.JSON(fiber.Map{"message": "User approved Successful"})

}

func RejectUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	// check
	user := &models.User{}
	if err := config.DB.First(user, user_id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User Not Found"})
	}

	// ถ้าเจอ
	user.Status = "rejected"
	if err := config.DB.Save(user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Rejection failed"})
	}

	return c.JSON(fiber.Map{"message": "User Rejected"})

}
