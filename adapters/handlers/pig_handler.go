package handlers

import (
	"backend/dto"
	"backend/usecases"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type HttpPigHandler struct {
	service  *usecases.PigService
	validate *validator.Validate
}

func NewHttpPigHandler(service *usecases.PigService, validate *validator.Validate) *HttpPigHandler {
	return &HttpPigHandler{
		service:  service,
		validate: validate}
}

func (h *HttpPigHandler) CreatePig(c *fiber.Ctx) error {
	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(400, "invalid user_id")
	}

	var input dto.PigInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input format")
	}
	if err := h.validate.Struct(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	resp, err := h.service.CreatePig(input, uint(userId64))

	if err != nil {
		switch {
		// กลุ่ม Error 400 (Bad Request) - ข้อมูลนำเข้าไม่ถูกต้อง
		case errors.Is(err, usecases.ErrInvalidWeight),
			errors.Is(err, usecases.ErrMaleAsMother),
			errors.Is(err, usecases.ErrFemaleAsFather),
			errors.Is(err, usecases.ErrMaleInvalidStatus),
			errors.Is(err, usecases.ErrPigTypeInvalidStatus),
			errors.Is(err, usecases.ErrCantFuture),
			errors.Is(err, usecases.ErrInvalidDate),
			errors.Is(err, usecases.ErrPigCodeAlreadyExists):
			// ส่งข้อความ Error จาก Service ออกไปตรงๆ ได้เลย (เพราะเราเขียนภาษาไทยไว้แล้ว)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())

		// กรณีอื่นๆ (เช่น Database ล่ม) -> ส่ง 500
		default:
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pig created successfully",
		"data":    resp,
	})
}
func (h *HttpPigHandler) UpdatePig(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "Invalid ID")
	}

	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(400, "invalid user_id")
	}

	var input dto.PigUpdate
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input format")
	}
	if err := h.validate.Struct(input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	pig, err := h.service.UpdatePig(uint(id), input, uint(userId64))
	if err != nil {
		return err
		// switch {
		// // กลุ่ม Error 400 (Bad Request) - ข้อมูลนำเข้าไม่ถูกต้อง
		// case errors.Is(err, usecases.ErrInvalidWeight),
		// 	errors.Is(err, usecases.ErrMaleAsMother),
		// 	errors.Is(err, usecases.ErrFemaleAsFather),
		// 	errors.Is(err, usecases.ErrMaleInvalidStatus),
		// 	errors.Is(err, usecases.ErrPigTypeInvalidStatus),
		// 	errors.Is(err, usecases.ErrCantFuture),
		// 	errors.Is(err, usecases.ErrFailedUpdate),
		// 	errors.Is(err, usecases.ErrInvalidDate),
		// 	errors.Is(err, usecases.ErrPigNotFound),
		// 	errors.Is(err, usecases.ErrPigCodeAlreadyExists):
		// 	// ส่งข้อความ Error จาก Service ออกไปตรงๆ ได้เลย (เพราะเราเขียนภาษาไทยไว้แล้ว)
		// 	return fiber.NewError(fiber.StatusBadRequest, err.Error())

		// // กรณีอื่นๆ (เช่น Database ล่ม) -> ส่ง 500
		// default:
		// 	return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		// }
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pig created successfully",
		"data":    pig,
	})
}

func (h *HttpPigHandler) FindAllPagination(c *fiber.Ctx) error {
	var param dto.PigParam
	if err := c.QueryParser(&param); err != nil {
		return fiber.NewError(400, "invalid query param")
	}
	result, err := h.service.FindAllPagination(param)
	if err != nil {
		return fiber.NewError(500, err.Error())
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "get pigs success",
		"data":    result,
	})
}

func (h *HttpPigHandler) GetPigByID(c *fiber.Ctx) error {
	id := c.Params("id")
	pig_id, err := strconv.Atoi(id)
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	pig, err := h.service.GetPigByID(uint(pig_id))
	if err != nil {
		return err
	}
	return c.JSON(pig)
}

func (h *HttpPigHandler) DeletePig(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	if err := h.service.DeletePig(uint(id)); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "delete success",
		"id":      id,
	})

}
