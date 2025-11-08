package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID    uint `json:"user_id"`
	SessionID uint `json:"session_id"`
	jwt.RegisteredClaims
}

type ResetPayload struct {
	UserID uint  `json:"user_id"`
	Exp    int64 `json:"exp"`
}

type RegisterPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Exp      int64  `json:"exp"`
	Password string `json:"password"`
}
