package dto

import "time"

type HealthInput struct {
	PigID  uint   `json:"pig_id" `
	Date   string `json:"date" validate:"required"`
	Type   string `json:"type" validate:"required,oneof=วัคซีน สุขภาพ"`
	Detail string `json:"detail" validate:"required"`
	Note   string `json:"note"`
}

type HealthUpdate struct {
	PigID  *uint   `json:"pig_id,omitempty"`
	Date   *string `json:"date,omitempty"`
	Type   *string `json:"type,omitempty" validate:"omitempty,oneof=วัคซีน สุขภาพ"`
	Detail *string `json:"detail,omitempty"`
	Note   *string `json:"note,omitempty"`
}

type HealthResponse struct {
	ID          uint      `json:"id"`
	PigID       uint      `json:"pig_id"`
	PigCodeName string    `json:"pig_code_name"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Detail      string    `json:"detail"`
	Note        string    `json:"note"`
	CreatedName string    `json:"created_name"`
	CreatedRole string    `json:"created_role"`
	UpdatedName string    `json:"updated_name"`
	UpdatedRole string    `json:"updated_role"`
}
