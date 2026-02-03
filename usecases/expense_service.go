package usecases

import (
	"backend/dto"
	"backend/models"
	"errors"
	"fmt"
	"time"
)

type ExpenseService struct {
	expenseRepo ExpenseRepository
}

func NewExpenseService(expenseRepo ExpenseRepository) *ExpenseService {
	return &ExpenseService{expenseRepo: expenseRepo}
}

func (s *ExpenseService) CreateExpense(input dto.ExpenseInput) (*models.Expense, error) {

	parsedDate, err := time.Parse("2006-01-02", input.Date)
	if err != nil {
		return nil, errors.New("invalid date, expect YYYY-MM-DD")
	}

	if parsedDate.After(time.Now()) {
		return nil, errors.New("time cannot be in the future")
	}

	if input.Amount < 0 {
		return nil, fmt.Errorf("amount must more than 0")
	}
	expense := models.Expense{
		Date:     parsedDate,
		Category: input.Category,
		Amount:   input.Amount,
		Note:     input.Note,
	}
	err = s.expenseRepo.Create(&expense)
	if err != nil {
		return nil, fmt.Errorf("failed to create expense: %v", err)
	}
	return &expense, nil

}

func (s *ExpenseService) GetAllExpense() ([]models.Expense, error) {
	return s.expenseRepo.GetAll()
}

func (s *ExpenseService) UpdateExpense(id uint, input dto.ExpenseUpdate) (*models.Expense, error) {
	updates := make(map[string]interface{})

	if input.Date != nil {
		parsedDate, err := time.Parse("2006-01-02", *input.Date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format expect YYYY-MM-DD")
		}
		if parsedDate.After(time.Now()) {
			return nil, fmt.Errorf("time cannot be in future")
		}
		updates["date"] = parsedDate
	}
	if input.Category != nil {
		updates["category"] = input.Category
	}
	if input.Amount != nil {
		updates["amount"] = input.Amount
	}
	if input.Note != nil {
		updates["note"] = input.Note
	}

	if err := s.expenseRepo.Update(id, updates); err != nil {
		return nil, err
	}
	return s.expenseRepo.FindByID(id)

}

func (s *ExpenseService) GetFindByID(id uint) (*models.Expense, error) {
	return s.expenseRepo.FindByID(id)
}

func (s *ExpenseService) SearchExpense(keyword string) ([]models.Expense, error) {
	return s.expenseRepo.Search(keyword)
}

func (s *ExpenseService) DeleteExpense(id uint) error {
	return s.expenseRepo.Delete(id)
}
