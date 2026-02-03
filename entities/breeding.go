package entities

import "time"

type Breeding struct {
	ID       string
	FatherID string
	MotherID string

	FatherCodename string
	MotherCodename string
	CreatedName    string
	CreatedRole    string
	UpdatedName    string
	UpdatedRole    string

	BreedingDate  time.Time
	ExpectedBirth time.Time
	Status        string
	Result        string
	Note          string

	CreatedBy string
	UpdatedBy string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Breeding) DaysUntilBirth() int {
	remaining := time.Until(b.ExpectedBirth)

	return int(remaining.Hours() / 24)
}

func (b *Breeding) IsPregnant() bool {
	return b.Status == "อุ้มท้อง"
}

// func (b *Breeding) IsOverdue() bool {
// 	return b.Status == "อุ้มท้อง" && time.Now().After(b.ExpectedBirth)
// }
