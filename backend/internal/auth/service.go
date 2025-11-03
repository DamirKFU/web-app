package auth

import (
	"app/internal/core"
	"app/internal/users"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	server *core.Server
}

func NewAuthService(server *core.Server) *AuthService {
	return &AuthService{server: server}
}

func (service *AuthService) Register(username, password, email string) error {
	var existing users.User
	if err := service.server.DB.Where("username = ?", username).First(&existing).Error; err == nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"username": "already exists"},
		}
	}
	if err := service.server.DB.Where("normalize_email = ?", users.NormalizeEmail(email)).First(&existing).Error; err == nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"email": "already exists"},
		}
	}

	user := users.User{Username: username, Email: email}
	if err := user.SetPassword(password, service.server.Cfg.SecretKey); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if err := service.server.DB.Create(&user).Error; err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (service *AuthService) Login(c *gin.Context, username, password, oldrefreshToken string) (string, string, error) {
	var user users.User

	if oldrefreshToken != "" {
		_ = service.server.RedisServer.RDB0.Del(c.Request.Context(), oldrefreshToken).Err()
	}

	if err := service.server.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	}

	if !user.CheckPassword(password, service.server.Cfg.SecretKey) {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	}

	accessToken, refreshToken, err := GenerateTokens(
		user.ID, service.server.Cfg.SecretKey,
		service.server.Cfg.JWT.AccessExpiresIn,
		service.server.Cfg.JWT.RefreshExpiresIn,
	)

	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "could not generate tokens",
		}
	}

	err = service.server.RedisServer.RDB0.Set(
		c.Request.Context(),
		refreshToken,
		user.ID,
		service.server.Cfg.JWT.RefreshExpiresIn,
	).Err()

	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "could not save refresh token",
		}
	}

	return accessToken, refreshToken, nil
}

func (service *AuthService) RefreshTokens(c *gin.Context, oldRefreshToken string) (string, string, error) {
	_, err := service.server.RedisServer.RDB0.Get(c, oldRefreshToken).Result()
	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "refresh token not found or expired",
		}
	}

	token, err := jwt.ParseWithClaims(oldRefreshToken, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(service.server.Cfg.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid refresh token",
		}
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid claims",
		}
	}

	newAccess, newRefresh, err := GenerateTokens(
		claims.UserID,
		service.server.Cfg.SecretKey,
		service.server.Cfg.JWT.AccessExpiresIn,
		service.server.Cfg.JWT.RefreshExpiresIn,
	)

	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "could not generate tokens",
		}
	}

	service.server.RedisServer.RDB0.Del(c.Request.Context(), oldRefreshToken)

	err = service.server.RedisServer.RDB0.Set(
		c.Request.Context(),
		newRefresh,
		claims.UserID,
		service.server.Cfg.JWT.RefreshExpiresIn,
	).Err()

	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "could not save new refresh token",
		}
	}

	return newAccess, newRefresh, nil
}

func (service *AuthService) Logout(c *gin.Context, token string) {
	if token != "" {
		service.server.RedisServer.RDB0.Del(c.Request.Context(), token).Err()
	}
}

func (service *AuthService) ForgotPassword(c *gin.Context, email string) error {
	var user users.User
	if err := service.server.DB.Where("LOWER(email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil
	}

	payload := ResetPayload{
		UserID: user.ID,
		Exp:    time.Now().Add(service.server.Cfg.PayloadTokenExpiresIn).Unix(),
	}

	token, err := core.GeneratePayloadToken(payload, service.server.Cfg.SecretKey)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", service.server.Cfg.FrontURL, token)
	subject := "Сброс пароля"
	data := struct {
		Username string
		ResetURL string
	}{
		Username: "Damir",
		ResetURL: resetURL,
	}

	body, err := core.RenderTextTemplate("./templates/reset_password.txt", data)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
			Fields:  nil,
		}
	}
	go func() {
		if err := service.server.Email.SendMail(
			subject,
			[]byte(body),
			[]string{user.Email},
		); err != nil {
			log.Printf("[Email Error] failed to send mail to %s: %v", user.Email, err)
		}
		log.Printf("[DEBUG] send mail to %s", user.Email)
	}()

	return nil
}

func (service *AuthService) ResetPassword(c *gin.Context, token string, newPassword string) error {
	payload, err := core.VerifyPayloadToken[ResetPayload](token, service.server.Cfg.SecretKey)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "invalid token",
			Fields:  map[string]string{"token": "invalid or malformed"},
		}
	}

	if time.Now().Unix() > payload.Exp {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "token expired",
			Fields:  map[string]string{"token": "expired"},
		}
	}

	var user users.User
	if err := service.server.DB.First(&user, payload.UserID).Error; err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if err := user.SetPassword(newPassword, service.server.Cfg.SecretKey); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if err := service.server.DB.Save(&user).Error; err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
