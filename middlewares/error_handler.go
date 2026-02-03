package middlewares

import (
	"backend/usecases"
	"errors"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default = 500 Internal Server Error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// 2. เช็คว่าเป็น Fiber Error หรือไม่ (เช่น 404 Route Not Found)
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		message = e.Message
	}
	// 3. เช็ค Domain Error ของเรา (Centralized Error Mapping)

	switch {
	// กลุ่ม 404 Not Found
	case errors.Is(err, usecases.ErrBreedingNotFound),
		errors.Is(err, usecases.ErrFatherNotFound),
		errors.Is(err, usecases.ErrMotherNotFound),
		errors.Is(err, usecases.ErrFoodNotFound),
		errors.Is(err, usecases.ErrNoValidPigs),
		errors.Is(err, usecases.ErrFoodStockNotFound),
		errors.Is(err, usecases.ErrUserNotFound):
		code = fiber.StatusNotFound // 404
		message = err.Error()       // ใช้ข้อความ "user not found" จากต้นฉบับ

	// กลุ่ม 400 Bad Request (Validation / Business Rules)
	case errors.Is(err, usecases.ErrEmailAlready),
		errors.Is(err, usecases.ErrInvalidCredentials),
		errors.Is(err, usecases.ErrCurrentPasswordWrong),
		errors.Is(err, usecases.ErrInvalidWeight),
		errors.Is(err, usecases.ErrMaleAsMother),
		errors.Is(err, usecases.ErrFemaleAsFather),
		errors.Is(err, usecases.ErrMaleInvalidStatus),
		errors.Is(err, usecases.ErrPigCodeAlreadyExists),
		errors.Is(err, usecases.ErrFailedUpdate),
		errors.Is(err, usecases.ErrInvalidDate),
		errors.Is(err, usecases.ErrPigNotFound),
		errors.Is(err, usecases.ErrPigCodeAlreadyExists),
		errors.Is(err, usecases.ErrCantFuture),

		errors.Is(err, usecases.ErrSamePig),
		errors.Is(err, usecases.ErrDuplicateBreeding),
		errors.Is(err, usecases.ErrInvalidMotherBreeder),
		errors.Is(err, usecases.ErrPigNotReady),
		errors.Is(err, usecases.ErrInvalidFatherBreeder),
		errors.Is(err, usecases.ErrIsUsedInBreeding),

		errors.Is(err, usecases.ErrNotEnoughFood),
		errors.Is(err, usecases.ErrFoodNotZero),
		errors.Is(err, usecases.ErrFeedingNotFound),

		errors.Is(err, usecases.ErrInvalidStockAmount),
		errors.Is(err, usecases.ErrNotEnoughFoodStock),
		errors.Is(err, usecases.ErrFoodStockAlreadyExists),
		errors.Is(err, usecases.ErrFoodStockUsed),

		errors.Is(err, usecases.ErrApproved):
		// เพิ่ม Error อื่นๆ ต่อท้ายตรงนี้ได้เลย

		code = fiber.StatusBadRequest // 400
		message = err.Error()

	// กลุ่ม 401 Unauthorized
	case errors.Is(err, usecases.ErrInvalidOrExpiry):
		code = fiber.StatusUnauthorized // 401
		message = err.Error()

	// กลุ่ม 403 Forbidden
	case errors.Is(err, usecases.ErrAccountPending),
		errors.Is(err, usecases.ErrAccountReject),
		errors.Is(err, usecases.ErrEmailNotVerify):
		code = fiber.StatusForbidden // 403
		message = err.Error()
	}

	// 4. ส่ง Response
	if os.Getenv("APP_ENV") == "dev" {
		// ถ้าเป็น Dev ให้ส่ง error ตัวเต็มไปด้วย เพื่อ debug ง่าย
		return c.Status(code).JSON(fiber.Map{
			"status":  "error",
			"message": message,
			"debug":   err.Error(),
		})
	} else {
		// ถ้าเป็น Production ต้องปิดบัง 500 ไว้
		if code == fiber.StatusInternalServerError {
			message = "เกิดข้อผิดพลาดภายในระบบ"
		}
		return c.Status(code).JSON(fiber.Map{
			"status":  "error",
			"message": message,
		})
	}
}
