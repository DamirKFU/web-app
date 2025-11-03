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
)

type AuthService struct {
	server         *core.Server
	userManager    *users.UserManager
	sessionManager *SessionManager
}

func NewAuthService(server *core.Server) *AuthService {
	return &AuthService{
		server:         server,
		userManager:    users.NewUserManager(server),
		sessionManager: NewSessionManager(server),
	}
}

func (service *AuthService) Register(username, password, email string) error {
	if existing, err := service.userManager.GetByUsername(username); err == nil && existing != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"username": "already exists"},
		}
	}

	if existing, err := service.userManager.GetByEmail(email); err == nil && existing != nil {
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

	if err := service.userManager.Create(&user); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (service *AuthService) Login(c *gin.Context, username, password, oldRefreshToken string) (string, string, error) {
	if oldRefreshToken != "" {
		service.server.RedisServer.RDB0.Del(c.Request.Context(), oldRefreshToken)
	}

	user, err := service.userManager.GetByUsername(username)
	if err != nil || !user.CheckPassword(password, service.server.Cfg.SecretKey) {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	}

	session, err := service.sessionManager.Create(user.ID)
	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: "could not create session",
		}
	}

	accessToken, refreshToken, err := GenerateTokens(
		user.ID,
		session.ID,
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

	claims, err := ParseToken(oldRefreshToken, service.server.Cfg.SecretKey)
	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid refresh token",
		}
	}

	newAccess, newRefresh, err := GenerateTokens(
		claims.UserID,
		claims.SessionID,
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

	service.server.RedisServer.RDB0.Set(
		c.Request.Context(),
		newRefresh,
		claims.UserID,
		service.server.Cfg.JWT.RefreshExpiresIn,
	)

	return newAccess, newRefresh, nil
}

func (service *AuthService) Logout(c *gin.Context, token string) {
	if token == "" {
		return
	}

	service.server.RedisServer.RDB0.Del(c.Request.Context(), token)

	claims, err := ParseToken(token, service.server.Cfg.SecretKey)
	if err == nil {
		service.sessionManager.Delete(claims.SessionID)
	}
}

func (service *AuthService) ForgotPassword(c *gin.Context, email string) error {
	user, err := service.userManager.GetByEmail(strings.ToLower(email))
	if err != nil {
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
		Username: user.Username,
		ResetURL: resetURL,
	}

	body, err := core.RenderTextTemplate("./templates/reset_password.txt", data)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	go func() {
		if err := service.server.Email.SendMail(
			subject,
			[]byte(body),
			[]string{user.Email},
		); err != nil {
			log.Printf("[Email Error] failed to send mail to %s: %v", user.Email, err)
		} else {
			log.Printf("[DEBUG] sent password reset email to %s", user.Email)
		}
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

	user, err := service.userManager.GetByID(payload.UserID)
	if err != nil {
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

	if err := service.userManager.Save(user); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
