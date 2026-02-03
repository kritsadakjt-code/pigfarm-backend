package usecases

import (
	"backend/models"
)

type ExpenseRepository interface {
	Create(expense *models.Expense) error
	Search(keyword string) ([]models.Expense, error)
	GetAll() ([]models.Expense, error)
	FindByID(id uint) (*models.Expense, error)
	Update(id uint, updates map[string]interface{}) error
	Delete(id uint) error
}
