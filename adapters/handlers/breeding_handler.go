package handlers

import (
	"backend/dto"
	"backend/usecases"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type BreedingHttpHandler struct {
	service  *usecases.BreedingService
	validate *validator.Validate
}

func NewBreedingHttpHandler(service *usecases.BreedingService, validate *validator.Validate) *BreedingHttpHandler {
	return &BreedingHttpHandler{
		service:  service,
		validate: validate}
}

func (h *BreedingHttpHandler) CreateBreeding(c *fiber.Ctx) error {
	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(400, "invalid userID")
	}
	var input dto.BreedingInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "Invalid input format")
	}
	if err := h.validate.Struct(input); err != nil {
		return fiber.NewError(400, err.Error())
	}
	breeding, err := h.service.CreateBreeding(input, uint(userId64))
	if err != nil {
		return err
	}
	return c.Status(201).JSON(fiber.Map{
		"message": "Create breeding success",
		"data":    breeding,
	})
}

func (h *BreedingHttpHandler) UpdateBreeding(c *fiber.Ctx) error {
	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(400, "invalid user")
	}
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "Invalid ID")
	}
	var input dto.BreedingUpdate
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "invalid input")
	}

	if err := h.validate.Struct(input); err != nil {
		return fiber.NewError(400, err.Error())
	}

	breeding, err := h.service.UpdateBreeding(input, uint(id), uint(userId64))
	if err != nil {
		return fiber.NewError(400, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update breeding successful",
		"data":    breeding,
	})
}

func (h *BreedingHttpHandler) FindAllPagination(c *fiber.Ctx) error {
	var param dto.BreedingParam
	if err := c.QueryParser(&param); err != nil {
		return fiber.NewError(400, "Invalid query parameters")
	}
	breeding, err := h.service.GetAllBreeding(param)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Fetch breeding list successful",
		"data":    breeding,
	})
}

func (h *BreedingHttpHandler) DeleteBreeding(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	if err := h.service.DeleteBreeding(uint(id)); err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Breeding record deleted successfully",
	})
}
