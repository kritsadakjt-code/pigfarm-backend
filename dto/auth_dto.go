package dto

import "time"

type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Status   string `json:"status"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	FullName  string    `json:"fullName"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Status    string    `json:"status"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateUserReq struct {
	FullName *string `json:"full_name"` // ใช้ *string เพื่อเช็คว่าส่งมาแก้หรือไม่ (nil = ไม่ส่ง)
	Email    *string `json:"email"`
	Phone    *string `json:"phone"`
	Role     *string `json:"role"`
	Status   *string `json:"status"`
}

// get all เเบบ pagination + search
// รับค่าจาก Query Param (เช่น ?page=1&limit=10&search=john)
type UserParam struct {
	Page   int    `query:"page"`
	Limit  int    `query:"limit"`
	Search string `query:"search"`
}

// โครงสร้าง response เเบบเเบ่งหน้า
type PaginationResponse struct {
	Data     []UserResponse `json:"data"`
	Total    int64          `json:"total"`     // จํานวนรายการทั้งหมดที่ค้นเจอ
	Page     int            `json:"page"`      // หน้าปัจจุบัน
	LastPage int            `json:"last_page"` // หน้าสุดท้าย
	Limit    int            `json:"limit"`     // จํานวนต่อหน้า
}
