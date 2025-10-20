package controllers

import (
	"backend/config"
	"backend/dto"
	"backend/models"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateExpense(c *fiber.Ctx) error {

	var input dto.ExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input " + err.Error()})
	}

	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
	}
	if parsedDate.After(time.Now()) {
		return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
	}

	if input.Amount <= 0 {
		return c.Status(400).JSON(fiber.Map{"error": "amount must be more than 0"})
	}

	expense := &models.Expense{
		Date:     parsedDate,
		Category: input.Category,
		Amount:   input.Amount,
		Note:     input.Note,
	}

	if err := config.DB.Create(expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create expense"})
	}

	resp := dto.ExpenseResponse{
		ID:       expense.ID,
		Date:     expense.Date,
		Category: expense.Category,
		Amount:   expense.Amount,
		Note:     expense.Note,
	}

	return c.Status(201).JSON(resp)
}

func SearchExpense(c *fiber.Ctx) error {
	keyword := c.Query("keyword")

	var expenses []models.Expense
	db := config.DB.Model(&models.Expense{})

	if keyword != "" {
		kw := "%" + keyword + "%"
		db = db.Where("category ILIKE ?", kw)
	}

	if err := db.Find(&expenses).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch expenses"})
	}

	var resp []dto.ExpenseResponse
	for _, e := range expenses {
		resp = append(resp, dto.ExpenseResponse{
			ID:       e.ID,
			Date:     e.Date,
			Category: e.Category,
			Amount:   e.Amount,
			Note:     e.Note,
		})
	}

	return c.JSON(resp)
}

func GetAllExpense(c *fiber.Ctx) error {
	var expenses []models.Expense

	// ดึงข้อมูลทั้งหมด
	if err := config.DB.Find(&expenses).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "cannot fetch expenses"})
	}

	var resp []dto.ExpenseResponse
	for _, e := range expenses {
		resp = append(resp, dto.ExpenseResponse{
			ID:       e.ID,
			Date:     e.Date,
			Category: e.Category,
			Amount:   e.Amount,
			Note:     e.Note,
		})
	}

	return c.JSON(resp)
}

func GetExpenseByID(c *fiber.Ctx) error {
	id := c.Params("id")
	expendID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var expense models.Expense
	if err := config.DB.First(&expense, "id = ?", expendID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Expense ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	resp := dto.ExpenseResponse{
		ID:       expense.ID,
		Date:     expense.Date,
		Category: expense.Category,
		Amount:   expense.Amount,
		Note:     expense.Note,
	}

	return c.JSON(resp)
}

func UpdateExpense(c *fiber.Ctx) error {
	id := c.Params("id")
	expendID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}
	var expense models.Expense
	if err := config.DB.First(&expense, "id = ?", expendID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Expense ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	var input dto.ExpenseUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input " + err.Error()})
	}
	if err := validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	updates := make(map[string]interface{})

	if input.Date != nil {
		parsedDate, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid date format, expected YYYY-MM-DD"})
		}
		if parsedDate.After(time.Now()) {
			return c.Status(400).JSON(fiber.Map{"error": "Time cannot be in the future"})
		}
		updates["date"] = parsedDate
	}
	if input.Category != nil {
		updates["category"] = *input.Category
	}
	if input.Amount != nil {
		if *input.Amount <= 0 {
			return c.Status(400).JSON(fiber.Map{"error": "Amount must be greater than 0"})
		}
		updates["amount"] = *input.Amount
	}
	if input.Note != nil {
		updates["note"] = *input.Note
	}

	if err := config.DB.Model(&expense).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Fail to update Expense"})
	}

	resp := dto.ExpenseResponse{
		ID:       expense.ID,
		Date:     expense.Date,
		Category: expense.Category,
		Amount:   expense.Amount,
		Note:     expense.Note,
	}

	return c.JSON(resp)
}

func DeleteExpense(c *fiber.Ctx) error {
	id := c.Params("id")
	expendID, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID: " + err.Error()})
	}

	var expense models.Expense
	if err := config.DB.First(&expense, "id = ?", expendID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Not Found Expense ID"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	if err := config.DB.Delete(&expense).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete Expense"})
	}

	return c.JSON(fiber.Map{
		"message": "Expense deleted successfully",
		"id":      expense.ID,
	})
}
