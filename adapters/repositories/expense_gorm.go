package repositories

import (
	"backend/models"
	"backend/usecases"

	"gorm.io/gorm"
)

type ExpenseGormRepository struct {
	db *gorm.DB
}

func NewExpenseGormRepository(db *gorm.DB) usecases.ExpenseRepository {
	return &ExpenseGormRepository{db: db}
}

func (r *ExpenseGormRepository) Create(expense *models.Expense) error {
	return r.db.Create(expense).Error
}

func (r *ExpenseGormRepository) GetAll() ([]models.Expense, error) {
	var expense []models.Expense
	err := r.db.Order("date desc").Find(&expense).Error
	return expense, err
}

func (r *ExpenseGormRepository) FindByID(id uint) (*models.Expense, error) {
	var expense models.Expense
	if err := r.db.First(&expense, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &expense, nil
}
func (r *ExpenseGormRepository) Update(id uint, updates map[string]interface{}) error {
	result := r.db.Model(&models.Expense{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

func (r *ExpenseGormRepository) Search(keyword string) ([]models.Expense, error) {
	var expense []models.Expense
	kw := "%" + keyword + "%"

	err := r.db.Where("category ILIKE ? OR note ILIKE ?", kw, kw).Find(&expense).Error

	return expense, err
}

func (r *ExpenseGormRepository) Delete(id uint) error {
	result := r.db.Delete(&models.Expense{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
