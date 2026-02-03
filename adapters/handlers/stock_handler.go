package handlers

import (
	"backend/dto"
	"backend/usecases"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpStockHandler struct {
	stockService *usecases.StockService
	validate     *validator.Validate
}

func NewHttpStockHandler(stockService *usecases.StockService, validate *validator.Validate) *HttpStockHandler {
	return &HttpStockHandler{
		stockService: stockService,
		validate:     validate}
}

func (s *HttpStockHandler) CreateFoodStock(c *fiber.Ctx) error {
	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(400, "invalid userID")
	}
	var input dto.FoodStockInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "invalid input")
	}
	if err := s.validate.Struct(input); err != nil {
		return err
	}

	resp, err := s.stockService.CreateFoodStock(&input, uint(userId64))
	if err != nil {
		return err
	}

	return c.Status(201).JSON(resp)
}

func (h *HttpStockHandler) UpdateFoodStock(c *fiber.Ctx) error {
	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(400, "invalid user_id")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	var input dto.FoodStockUpdate
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "invalid input")
	}
	if err := h.validate.Struct(input); err != nil {
		return err
	}

	resp, err := h.stockService.UpdateFoodStock(uint(id), uint(userId64), input)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

func (h *HttpStockHandler) GetAllPagi(c *fiber.Ctx) error {
	var input dto.ParamFoodStock
	if err := c.QueryParser(&input); err != nil {
		return fiber.NewError(400, "invalid input")
	}
	stock, err := h.stockService.GetAllPagi(input)
	if err != nil {
		return err
	}
	return c.JSON(stock)
}

func (h *HttpStockHandler) DeleteFoodStock(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	if err := h.stockService.DeleteFoodStock(uint(id)); err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": "delete food stock success"})
}
