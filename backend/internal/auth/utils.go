package auth

import (
	"app/internal/core"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateTokens(userID uint, secretKey string, accessExpiresIn time.Duration, refreshExpiresIn time.Duration) (string, string, error) {
	accessClaims := &core.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessExpiresIn)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &core.Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiresIn)),
		},
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(secretKey))

	return accessToken, refreshToken, err
}
