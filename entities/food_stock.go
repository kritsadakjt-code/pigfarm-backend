package entities

import "time"

type FoodStock struct {
	ID string

	FoodTypeID string
	Amount     float64
	DateTime   time.Time
	Note       string

	FoodTypeName string

	CreatedBy   string
	CreatedRole string
	UpdatedBy   string
	UpdatedRole string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (f *FoodStock) IsCriticalStock() bool {
	return f.Amount <= 5 && f.Amount > 0
}
