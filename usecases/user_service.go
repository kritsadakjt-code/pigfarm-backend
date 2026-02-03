package usecases

import (
	"backend/dto"
	"backend/models"
	"backend/utils"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(input dto.RegisterRequest) (*models.User, error) {
	email := strings.ToLower(input.Email)
	existingEmail, _ := s.userRepo.FindByEmail(email)
	if existingEmail != nil {
		return nil, ErrEmailAlready
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordHashFailed
	}

	verifyToken := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(30 * time.Minute)

	user := models.User{
		FullName:                input.FullName,
		Email:                   email,
		Password:                string(hashed),
		Phone:                   input.Phone,
		Role:                    input.Role,
		Status:                  "pending_verification_email", // สถานะรอการยืนยัน
		EmailVerificationToken:  verifyToken,
		EmailVerificationExpiry: &expiry,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return nil, errors.New("failed to register user")
	}

	// อาจเเยกเป็น emailservice อีกทีก็ได้
	frontendURL := os.Getenv("FRONTEND_URL_DEV")
	verifyURL := fmt.Sprintf("%s/verify-email/register/%s", frontendURL, verifyToken)

	// ส่งเเบบ async goroutine เพื่อไม่ให้ user รอนาน
	go func() {
		if err := utils.SendEmailVerification(user.Email, verifyURL); err != nil {
			fmt.Println("error sending verification email:", err)
		}
	}()
	return &user, nil
}

func (s *UserService) LoginUser(email string, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}

	switch user.Status {
	case "pending_verification_email":
		return "", nil, ErrEmailNotVerify
	case "pending":
		return "", nil, ErrAccountPending
	case "rejected":
		return "", nil, ErrAccountReject
		// ถ้าเป็น "approved" จะผ่านไปทำต่อ
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil

}

func (s *UserService) ForgetPassword(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return ErrInvalidCredentials
	}
	token := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(30 * time.Minute)
	user.ResetToken = token
	user.ResetTokenExpiry = &expiry

	if err := s.userRepo.Save(user); err != nil {
		return err
	}
	frontendURL := os.Getenv("FRONTEND_URL_DEV")
	resetURL := fmt.Sprintf("%s/reset-password/%s", frontendURL, token)

	// ส่ง link reset to email
	go func() {
		if err := utils.SentResetPasswordFromEmail(user.Email, resetURL); err != nil {
			fmt.Println("Error sending reset email:", err)
		}
	}()

	return nil
}

func (s *UserService) ResetPassword(token string, newPassword string) error {
	user, err := s.userRepo.FindByResetToken(token, time.Now())
	if err != nil {
		return errors.New("token invalid or expired")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashed)
	user.ResetToken = ""
	user.ResetTokenExpiry = nil

	if err := s.userRepo.Save(user); err != nil {
		return errors.New("failed to update password")
	}
	return nil
}

func (s *UserService) ChangePassword(id uint, currentPassword string, newPassword string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword)); err != nil {
		return ErrCurrentPasswordWrong
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return ErrPasswordHashFailed
	}
	user.Password = string(hashed)

	if err := s.userRepo.Save(user); err != nil {
		return errors.New("failed to update password")
	}
	return nil

}

func (s *UserService) RequestChangeEmail(user_id uint, email string) error {
	user, err := s.userRepo.FindByID(user_id)
	if err != nil {
		return ErrUserNotFound
	}
	exist, err := s.userRepo.ExistEmail(email)
	if err != nil {

		return errors.New("failed to find email")
	}
	if exist {

		return ErrEmailAlready
	}

	email_token := utils.GenerateRandomToken(32)
	expiry := time.Now().Add(30 * time.Minute)

	user.PendingEmail = email
	user.EmailVerificationToken = email_token
	user.EmailVerificationExpiry = &expiry

	if err := s.userRepo.Save(user); err != nil {
		return ErrFailed
	}

	frontendURL := os.Getenv("FRONTEND_URL_DEV")
	verifyURL := fmt.Sprintf("%s/verify-email/change/%s", frontendURL, email_token)

	go func() {
		if err := utils.SendEmailVerification(email, verifyURL); err != nil {
			fmt.Println("Error sending reset email:", err)
		}
	}()

	return nil
}

func (s *UserService) VerifyEmailForChange(token string) error {
	user, err := s.userRepo.FindByEmailVerifyToken(token, time.Now())
	if err != nil {
		return ErrInvalidOrExpiry
	}
	user.Email = user.PendingEmail
	user.PendingEmail = ""
	user.EmailVerificationToken = ""
	user.EmailVerificationExpiry = nil
	if err := s.userRepo.Save(user); err != nil {
		return ErrFailed
	}

	return nil
}

func (s *UserService) EmailVerifiedRegister(token string) error {
	user, err := s.userRepo.FindByEmailVerifyRegisterToken(token)
	if err != nil {
		return ErrInvalidOrExpiry
	}
	if user.EmailVerificationExpiry == nil || user.EmailVerificationExpiry.Before(time.Now()) {
		return ErrInvalidOrExpiry
	}
	now := time.Now()
	user.EmailVerifiedRegister = &now
	user.Status = "pending"

	user.EmailVerificationToken = ""
	user.EmailVerificationExpiry = nil

	if err := s.userRepo.Save(user); err != nil {
		return ErrFailed
	}
	return nil
}

func (s *UserService) ResendVerifyEmail(email string) error {
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return ErrUserNotFound
	}
	switch user.Status {
	case "approved":
		return ErrApproved
	case "pending":
		return ErrAccountPending
	case "rejected":
		return ErrAccountReject
	}
	newToken := utils.GenerateRandomToken(32)
	newExpiry := time.Now().Add(24 * time.Hour)

	user.EmailVerificationToken = newToken
	user.EmailVerificationExpiry = &newExpiry

	if err := s.userRepo.Save(user); err != nil {
		return ErrFailed
	}
	frontendURL := os.Getenv("FRONTEND_URL_DEV")
	verifyURL := fmt.Sprintf("%s/verify-email/register/%s", frontendURL, newToken)

	go func() {
		if err := utils.SendEmailVerification(user.Email, verifyURL); err != nil {
			fmt.Println("error sending verification email:", err)
		}
	}()
	return nil

}

func (s *UserService) CreateByAdmin(input dto.RegisterRequest) (*models.User, error) {
	email := strings.ToLower(input.Email)
	exist, err := s.userRepo.ExistEmail(email)
	if err != nil {
		return nil, ErrFailedToFindEmail
	}
	if exist {
		return nil, ErrEmailAlready
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, ErrPasswordHashFailed
	}
	user := &models.User{
		FullName: input.FullName,
		Email:    email,
		Password: string(hashed),
		Phone:    input.Phone,
		Role:     input.Role,
		Status:   input.Status,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, ErrFailed
	}
	return user, nil
}

func (s *UserService) UpdateUserByAdmin(id uint, req dto.UpdateUserReq) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if req.Email != nil && *req.Email != user.Email {
		existUser, err := s.userRepo.FindByEmail(*req.Email)
		if err == nil && existUser.ID != user.ID {
			return nil, ErrEmailAlready
		}
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.Phone != nil {
		user.Phone = *req.Phone
	}
	if req.Role != nil {
		user.Role = *req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}
	if err := s.userRepo.Save(user); err != nil {
		return nil, ErrFailed
	}
	return user, nil
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

func (s *UserService) FindAllPagination(param dto.UserParam) (*dto.PaginationResponse, error) {
	// set defualt ถ้าไม่ได้ส่งค่ามา
	if param.Page <= 0 {
		param.Page = 1
	}
	if param.Limit <= 0 {
		param.Limit = 10
	}
	users, total, err := s.userRepo.FindAllPagination(param)
	if err != nil {
		return nil, err
	}

	// เตรียม response
	var resp []dto.UserResponse
	for _, u := range users {
		resp = append(resp, dto.UserResponse{
			ID:        fmt.Sprintf("%d", u.ID),
			FullName:  u.FullName,
			Email:     u.Email,
			Role:      u.Role,
			Status:    u.Status,
			Phone:     u.Phone,
			CreatedAt: u.CreatedAt,
		})
	}

	// คํานวณหน้าสุดท้าย
	lastPage := int(math.Ceil(float64(total) / float64(param.Limit)))

	return &dto.PaginationResponse{
		Data:     resp,
		Total:    total,
		Page:     param.Page,
		LastPage: lastPage,
		Limit:    param.Limit,
	}, nil
}

func (s *UserService) FindByID(id uint) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) DeleteUser(id uint) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	if err := s.userRepo.Delete(id); err != nil {
		return ErrFailed
	}
	return nil
}

func (s *UserService) ApproveUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	return s.userRepo.UpdateStatus(user.ID, "approved")

}
func (s *UserService) RejectUser(id uint) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return ErrUserNotFound
	}
	return s.userRepo.UpdateStatus(user.ID, "rejected")

}
