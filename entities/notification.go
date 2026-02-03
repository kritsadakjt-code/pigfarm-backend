package entities

import "time"

// type Notification struct {
// 	gorm.Model
// 	Type    string `json:"type" gorm:"not null"` // "food_low", "birth_due", "health_alert"
// 	Title   string `json:"title" gorm:"not null"`
// 	Message string `json:"message" gorm:"not null"`
// 	IsRead  bool   `json:"is_read" gorm:"default:false"`
// }

// Pure domain entity - no framework tags
type Notification struct {
	ID        string // ให้เป็น string เพื่อให้รองรับทั้ง sql noSql กรณีต้องเปลี่ยน db
	Type      string
	Title     string
	Message   string
	IsRead    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

const (
	TypeFoodLow  = "food_low"
	TypeBirthDue = "birth_due"
)

func NewNotification(notiType, title, message string) *Notification {
	return &Notification{
		Type:    notiType,
		Title:   title,
		Message: message,
		IsRead:  false,
	}
}

func (n *Notification) MarkAsRead() {
	n.IsRead = true
}

func (n *Notification) IsUnread() bool {
	return !n.IsRead
}
