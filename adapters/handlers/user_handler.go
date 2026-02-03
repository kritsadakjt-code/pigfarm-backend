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

type HttpUserHandlers struct {
	service  *usecases.UserService
	validate *validator.Validate
}

func NewHttpUserHandlers(service *usecases.UserService, validate *validator.Validate) *HttpUserHandlers {
	return &HttpUserHandlers{
		service:  service,
		validate: validate}
}

func (h *HttpUserHandlers) CreateUser(c *fiber.Ctx) error {
	var input dto.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid Input"})
	}

	user, err := h.service.CreateUser(input)
	if err != nil {
		if errors.Is(err, usecases.ErrEmailAlready) {
			return c.Status(400).JSON(fiber.Map{"error": "email already exists"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "server failed"})
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Registration successful",
		"status":  user.Status,
	})
}

func (h *HttpUserHandlers) Login(c *fiber.Ctx) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON("invalid input")
	}
	token, user, err := h.service.LoginUser(input.Email, input.Password)
	if err != nil {
		switch err {
		case usecases.ErrInvalidCredentials:
			return c.Status(401).JSON(fiber.Map{"error": usecases.ErrInvalidCredentials})
		case usecases.ErrEmailNotVerify, usecases.ErrAccountPending, usecases.ErrAccountReject:
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		default:
			return c.Status(500).JSON(fiber.Map{"error": "server failed"})
		}

	}
	return c.JSON(fiber.Map{
		"token":    token,
		"fullname": user.FullName,
		"email":    user.Email,
		"role":     user.Role,
	})
}

func (h *HttpUserHandlers) ForgetPassword(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.service.ForgetPassword(input.Email); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to sent email"})
	}
	return c.Status(200).JSON(fiber.Map{"message": "หากมีบัญชีที่ใช้อีเมลนี้อยู่ ระบบได้ส่งลิงก์สำหรับรีเซ็ตรหัสผ่านไปให้แล้ว"})
}

func (h *HttpUserHandlers) ResetPassword(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Missing token"})
	}
	var input struct {
		NewPassword string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}
	if err := h.service.ResetPassword(token, input.NewPassword); err != nil {
		if errors.Is(err, errors.New("token invalid or has expiry")) {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		// return c.Status(500).JSON(fiber.Map{"error": "Failed to reset password"})
	}
	return c.JSON(fiber.Map{"message": "Password reset successfully"})

}

func (h *HttpUserHandlers) ChangePassword(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	id := uint(user_id.(float64))
	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}
	if err := h.service.ChangePassword(id, input.CurrentPassword, input.NewPassword); err != nil {
		if errors.Is(err, usecases.ErrCurrentPasswordWrong) {

			return c.Status(400).JSON(fiber.Map{"error": "รหัสผ่านปัจจุบันไม่ถูกต้อง"})
		}
		if errors.Is(err, usecases.ErrUserNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "ไม่พบข้อมูลผู้ใช้"})
		}

		return c.Status(500).JSON(fiber.Map{"error": "เกิดข้อผิดพลาดภายในระบบ"})
	}

	return c.JSON(fiber.Map{"message": "เปลี่ยนรหัสผ่านสำเร็จ"})
}

func (h *HttpUserHandlers) RequestChangeEmail(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	id := uint(user_id.(float64))

	var input struct {
		NewEmail string `json:"new_email"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := h.service.RequestChangeEmail(id, input.NewEmail); err != nil {
		if errors.Is(err, usecases.ErrEmailAlready) {
			return c.Status(400).JSON(fiber.Map{"error": "email is already"})
		}
		if errors.Is(err, usecases.ErrUserNotFound) {
			return c.Status(404).JSON(fiber.Map{"error": "not found user"})
		}

		return c.Status(500).JSON(fiber.Map{"error": "failed to request change email"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "ระบบได้ส่งลิงก์สำหรับยืนยันไปที่อีเมลใหม่ของคุณแล้ว",
	})
}

func (h *HttpUserHandlers) VerifyEmailForChange(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid token"})
	}
	if err := h.service.VerifyEmailForChange(token); err != nil {
		if errors.Is(err, usecases.ErrInvalidOrExpiry) {
			return c.Status(400).JSON(fiber.Map{"error": "invalid or expiry token"})
		}
		return c.Status(500).JSON(fiber.Map{"error": usecases.ErrFailed})
	}
	return c.JSON(fiber.Map{"message": "email updated successfully"})
}

func (h *HttpUserHandlers) VerifyEmailForRegister(c *fiber.Ctx) error {
	token := c.Params("token")
	if token == "" {
		return c.Status(400).JSON(fiber.Map{"error": "invalid token"})
	}
	if err := h.service.EmailVerifiedRegister(token); err != nil {
		if errors.Is(err, usecases.ErrInvalidOrExpiry) {
			return c.Status(400).JSON(fiber.Map{"error": "invalid or expiry token"})
		}
		return c.Status(500).JSON(fiber.Map{"error": usecases.ErrFailed})
	}
	return c.Status(200).JSON(fiber.Map{"message": "Email verified successfully. Please wait for admin approval."})
}

func (h *HttpUserHandlers) ResendVerifyEmail(c *fiber.Ctx) error {
	var input struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	if err := h.service.ResendVerifyEmail(input.Email); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "resend email success"})

}

func (h *HttpUserHandlers) CreateUserByAdmin(c *fiber.Ctx) error {
	var input dto.RegisterRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	user, err := h.service.CreateByAdmin(input)
	if err != nil {

		switch {
		case errors.Is(err, usecases.ErrEmailAlready):
			// return c.Status(400).JSON(fiber.Map{"error": "email already"})
			return fiber.NewError(400, "อีเมลนี้ถูกใช้งานเเล้ว")
		// error บางอย่างไม่ควรให้ผู้ใช้รู้
		// case  errors.Is(err, usecases.ErrPasswordHashFailed) :
		// 	return c.Status(500).JSON(fiber.Map{"error": usecases.ErrPasswordHashFailed})
		// case errors.Is(err, usecases.ErrFailedToFindEmail):
		// 	return c.Status(500).JSON(fiber.Map{"error": usecases.ErrFailedToFindEmail})
		default:
			// return c.Status(500).JSON(fiber.Map{"error": "Internal Server Error"})
			return err
		}

	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Create user successful",
		"data": dto.UserResponse{
			ID:        fmt.Sprintf("%d", user.ID),
			FullName:  user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *HttpUserHandlers) UpdateUserByAdmin(c *fiber.Ctx) error {
	id := c.Params("id")

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var req dto.UpdateUserReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(400, "invalid input")
	}
	// if err := h.validate.Struct(req); err != nil {
	// 	return fiber.NewError(400, err.Error())
	// }
	user, err := h.service.UpdateUserByAdmin(uint(user_id), req)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrUserNotFound):
			return fiber.NewError(fiber.StatusNotFound, "ไม่พบข้อมูลผู้ใช้")
		case errors.Is(err, usecases.ErrEmailAlready):
			return fiber.NewError(fiber.StatusBadRequest, "อีเมลนี้ถูกใช้งานโดยผู้ใช้อื่นแล้ว")
		default:
			return err // ส่งให้ Global Error Handler จัดการ 500
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Update user successful",
		"data": dto.UserResponse{
			ID:        fmt.Sprintf("%d", user.ID),
			FullName:  user.FullName,
			Email:     user.Email,
			Role:      user.Role,
			Status:    user.Status,
			Phone:     user.Phone,
			CreatedAt: user.CreatedAt,
		},
	})
}

func (h *HttpUserHandlers) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return err
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

func (h *HttpUserHandlers) FindAllPagination(c *fiber.Ctx) error {
	var param dto.UserParam
	if err := c.QueryParser(&param); err != nil {
		return fiber.NewError(400, "invalid query param")
	}

	users, err := h.service.FindAllPagination(param)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Get users successful",
		"data":    users,
	})
}

func (h *HttpUserHandlers) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	user, err := h.service.FindByID(uint(user_id))
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			return fiber.NewError(404, "Not found user")
		}
		return err
	}
	return c.JSON(user)
}

func (h *HttpUserHandlers) DeleteUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID")
	}
	if err := h.service.DeleteUser(uint(id)); err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

func (h *HttpUserHandlers) ApproveUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	if err := h.service.ApproveUser(uint(id)); err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Approved successfully",
	})
}

func (h *HttpUserHandlers) RejectUser(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fiber.NewError(400, "invalid id")
	}
	if err := h.service.RejectUser(uint(id)); err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			return fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User Rejected successfully",
	})
}
