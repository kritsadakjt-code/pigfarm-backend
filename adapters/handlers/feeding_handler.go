package handlers

import (
	"backend/dto"
	"backend/usecases"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type FeedingHttpHandler struct {
	service *usecases.FeedingService
}

func NewFeedingHttpHandler(service *usecases.FeedingService) *FeedingHttpHandler {
	return &FeedingHttpHandler{service: service}
}

func (h *FeedingHttpHandler) CreateFeeding(c *fiber.Ctx) error {
	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return fiber.NewError(401, "Invalid User ID")
	}

	var input dto.FeedingInput
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(400, "Invalid input format")
	}

	resp, err := h.service.CreateFeeding(input, uint(userId64))
	fmt.Println("sadasd")
	fmt.Println(resp, "testt")
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Feeding created successfully",
		"data":    resp,
	})

}

func (h *FeedingHttpHandler) GetAllFeedingPagination(c *fiber.Ctx) error {
	var param dto.ParamFeeding
	if err := c.QueryParser(&param); err != nil {
		return fiber.NewError(400, "invalid query param")
	}
	result, err := h.service.GetAllFeedingPagination(param)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Fetch feeding list successful",
		"data":    result,
	})
}

func (h *FeedingHttpHandler) GetFeedingByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}

	resp, err := h.service.GetFeedingByID(uint(id))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ดึงรายละเอียดการให้อาหารสำเร็จ",
		"data":    resp,
	})
}

func (h *FeedingHttpHandler) UpdateFeeding(c *fiber.Ctx) error {

	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง")
	}

	userIdStr := fmt.Sprintf("%v", c.Locals("user_id"))
	userId64, _ := strconv.ParseUint(userIdStr, 10, 64)

	var input dto.FeedingUpdate
	if err := c.BodyParser(&input); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid input format")
	}

	resp, err := h.service.UpdateFeeding(uint(id), input, uint(userId64))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "แก้ไขข้อมูลการให้อาหารสำเร็จ",
		"data":    resp,
	})
}

func (h *FeedingHttpHandler) DeleteFeeding(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "รูปแบบ ID ไม่ถูกต้อง")
	}

	err = h.service.DeleteFeeding(uint(id))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "ลบข้อมูลการให้อาหารและคืนยอดสต็อกสำเร็จ",
	})
}
