package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	user := &models.User{}
	if err := config.DB.First(user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	// Map user -> DTO
	res := dto.ProfileResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role,
	}

	return c.JSON(res)
}

// UpdateProfile
func UpdateProfile(c *fiber.Ctx) error {
	userIDStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid user_id"})
	}

	user := &models.User{}
	if err := config.DB.First(user, "id = ?", userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	var input dto.UpdateProfileRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := validate.Struct(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := make(map[string]interface{})
	if input.FullName != nil {
		updates["full_name"] = input.FullName
	}
	if input.Phone != nil {
		updates["phone"] = input.Phone
	}

	if err := config.DB.Model(user).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Update failed"})
	}

	res := dto.ProfileResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Role:     user.Role,
	}

	return c.JSON(res)
}
