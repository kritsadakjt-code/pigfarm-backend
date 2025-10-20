package dto

type ProfileResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`
	Role     string `json:"role"`
}

type UpdateProfileRequest struct {
	FullName *string `json:"full_name,omitempty" `
	Phone    *string `json:"phone,omitempty"`
}
