package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type ResetPayload struct {
	UserID uint  `json:"user_id"`
	Exp    int64 `json:"exp"`
}
