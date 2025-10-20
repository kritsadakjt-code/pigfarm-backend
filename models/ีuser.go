package models

import (
	"time"

	"gorm.io/gorm"
)

// type User struct {
// 	mgm.DefaultModel        `bson:",inline"`
// 	FullName                string     `json:"full_name" bson:"full_name"`
// 	Email                   string     `json:"email" bson:"email"`
// 	Password                string     `json:"password" bson:"password"`
// 	Phone                   string     `json:"phone" bson:"phone"`
// 	Role                    string     `json:"role" bson:"role"`
// 	Status                  string     `json:"status" bson:"status"`
// 	ResetToken              string     `bson:"reset_token,omitempty"`
// 	ResetTokenExpiry        time.Time  `bson:"reset_token_expiry,omitempty"`
// 	PendingEmail            string     `json:"-" bson:"pending_email,omitempty"`
// 	EmailVerificationToken  string     `json:"-" bson:"email_verification_token,omitempty"`
// 	EmailVerificationExpiry *time.Time `json:"-" bson:"email_verification_expiry,omitempty"`
// }

type User struct {
	gorm.Model
	FullName                string     `gorm:"size:100;not null" json:"full_name"`
	Email                   string     `gorm:"size:100;unique;not null" json:"email"`
	Password                string     `gorm:"size:255;not null" json:"-"`
	Phone                   string     `gorm:"size:20" json:"phone"`
	Role                    string     `gorm:"size:50;" json:"role"`   // owner, employee
	Status                  string     `gorm:"size:50;" json:"status"` // pending, active, inactive
	EmailVerifiedRegister   *time.Time `json:"email_verified_register"`
	ResetToken              string     `gorm:"size:255" json:"-"`
	ResetTokenExpiry        *time.Time `json:"-"`
	PendingEmail            string     `gorm:"size:100" json:"-"`
	EmailVerificationToken  string     `gorm:"size:255" json:"-"`
	EmailVerificationExpiry *time.Time `json:"-"`
}
