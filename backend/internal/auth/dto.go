package auth

import "app/internal/core"

type RegisterRequest struct {
	Username       string `json:"username" binding:"UsersUsername"`
	Email          string `json:"email" binding:"UsersEmail"`
	Password       string `json:"password" binding:"DtoUsersPassword"`
	RepeatPassword string `json:"repeat_password" binding:"RepeatUsersPassword"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"UsersEmail"`
}

type ResetPasswordRequest struct {
	core.AbstractTokenRequest
	Password       string `json:"password" binding:"DtoUsersPassword"`
	RepeatPassword string `json:"repeat_password" binding:"RepeatUsersPassword"`
}

type RegisterConfirmRequest struct {
	core.AbstractTokenRequest
}
