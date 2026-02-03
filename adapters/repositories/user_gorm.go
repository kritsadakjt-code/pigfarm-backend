package repositories

import (
	"backend/dto"
	"backend/models"
	"backend/usecases"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserGormRepository struct {
	db *gorm.DB
}

func NewUserGormRepository(db *gorm.DB) usecases.UserRepository {
	return &UserGormRepository{db: db}
}

func (r *UserGormRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserGormRepository) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserGormRepository) Save(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *UserGormRepository) FindByResetToken(token string, expiry time.Time) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, "reset_token = ? AND reset_token_expiry > ?", token, expiry).Error; err != nil {
		return nil, err
	}
	fmt.Println(user)
	return user, nil
}

func (r *UserGormRepository) FindByID(id uint) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserGormRepository) FindByEmailVerifyToken(token string, expiry time.Time) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, "email_verification_token = ? AND email_verification_expiry > ?", token, expiry).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserGormRepository) FindByEmailVerifyRegisterToken(token string) (*models.User, error) {
	user := &models.User{}
	if err := r.db.First(user, "email_verification_token = ? AND email_verified_register IS NULL", token).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserGormRepository) ExistEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *UserGormRepository) GetAll() ([]models.User, error) {
	var user []models.User
	if err := r.db.Order("created_at DESC").Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserGormRepository) FindAllPagination(param dto.UserParam) ([]models.User, int64, error) {
	var user []models.User
	var total int64

	// เริ่ม query
	db := r.db.Model(&models.User{})

	// ถ้ามี search ค้นหาจาก name email role
	if param.Search != "" {
		keyword := "%" + param.Search + "%"
		db = db.Where("full_name ILIKE ? OR email ILIKE ? OR role ILIKE ?", keyword, keyword, keyword)
	}
	// นับจํานวนทั้งหมด
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// จํานวนต่อหน้า
	offset := (param.Page - 1) * param.Limit
	if err := db.Offset(offset).Limit(param.Limit).Order("id DESC").Find(&user).Error; err != nil {
		return nil, 0, err
	}

	return user, total, nil
}

func (r *UserGormRepository) Delete(id uint) error {
	if err := r.db.Unscoped().Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserGormRepository) UpdateStatus(id uint, status string) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("record not found")
	}
	return nil
}
