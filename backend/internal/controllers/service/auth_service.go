package service

import (
	"app/internal/controllers/types"
	"app/internal/core"
	"app/internal/models"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	server *core.Server
}

func NewAuthService(server *core.Server) *AuthService {
	return &AuthService{server: server}
}

func (service *AuthService) Register(username, password string) (*models.User, error) {
	var existing models.User
	if err := service.server.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	user := models.User{Username: username}
	if err := user.SetPassword(password, service.server.Cfg.SecretKey); err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *AuthService) Login(username, password string) (uint, error) {
	var user models.User
	if err := service.server.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return 0, errors.New("invalid credentials")
	}

	if !user.CheckPassword(password, service.server.Cfg.SecretKey) {
		return 0, errors.New("invalid credentials")
	}

	return user.ID, nil
}

func (service *AuthService) RefreshTokens(c *gin.Context, oldRefreshToken string) (claims *types.Claims, err error) {
	_, err = service.server.RedisServer.RDB0.Get(c, oldRefreshToken).Result()

	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	token, err := jwt.ParseWithClaims(oldRefreshToken, &types.Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(service.server.Cfg.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(*types.Claims)
	if !ok {
		return nil, errors.New("invalid claims")
	}
	return claims, nil
}
