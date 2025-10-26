package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=42"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=42"`
	Password string `json:"password" validate:"required"`
}
