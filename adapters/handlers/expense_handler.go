package handlers

import (
	"backend/dto"
	"backend/usecases"
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// สร้าง validate เพื่อง่ายต่อการ test case
type HttpExpenseHandler struct {
	service  *usecases.ExpenseService
	validate *validator.Validate
}

func NewHttpExpenseHandler(service *usecases.ExpenseService, validate *validator.Validate) *HttpExpenseHandler {
	return &HttpExpenseHandler{
		service:  service,
		validate: validate,
	}
}

func (h *HttpExpenseHandler) CreateExpense(c *fiber.Ctx) error {
	var input dto.ExpenseInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input" + err.Error()})
	}
	if err := h.validate.Struct(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	expense, err := h.service.CreateExpense(input)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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

func (h *HttpExpenseHandler) GetAllExpense(c *fiber.Ctx) error {
	expense, err := h.service.GetAllExpense()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch expense" + err.Error()})
	}

	var resp []dto.ExpenseResponse
	for _, e := range expense {
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

func (h *HttpExpenseHandler) UpdateExpense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ID"})
	}

	var input dto.ExpenseUpdate
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.validate.Struct(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	expense, err := h.service.UpdateExpense(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Expense not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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

func (h *HttpExpenseHandler) GetExpenseByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ID"})
	}
	expense, err := h.service.GetFindByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Expense not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
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

func (h *HttpExpenseHandler) SearchExpense(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	expenses, err := h.service.SearchExpense(keyword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	var resp []dto.ExpenseResponse
	for _, e := range expenses {
		resp = append(resp, dto.ExpenseResponse{
			ID:       e.ID,
			Category: e.Category,
			Amount:   e.Amount,
			Note:     e.Note,
		})
	}
	return c.JSON(resp)
}

func (h *HttpExpenseHandler) DeleteExpense(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid ID"})
	}
	err = h.service.DeleteExpense(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "Expense not found"})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted successfully"})

}
