package auth

import (
	"app/internal/core"
	"app/internal/users"
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

func (service *AuthService) RegisterConfirm(c *gin.Context, token string) error {
	payload, err := core.VerifyPayloadToken[RegisterPayload](token, service.server.Cfg.SecretKey)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"token": "invalid or malformed"},
		}
	}

	if time.Now().Unix() > payload.Exp {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"token": "expired"},
		}
	}

	if existing, err := service.userManager.GetByUsername(c, payload.Username); err == nil && existing != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "username already exists",
		}
	}

	if existing, err := service.userManager.GetByEmail(c, payload.Email); err == nil && existing != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "email already exists",
		}
	}

	user := users.User{Username: payload.Username, Email: payload.Email}
	if err := user.SetPassword(payload.Password, service.server.Cfg.SecretKey); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if err := service.userManager.Create(c, &user); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (service *AuthService) Register(c *gin.Context, username, password, email string) error {
	if existing, err := service.userManager.GetByUsername(c, username); err == nil && existing != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"username": "already exists"},
		}
	}

	if existing, err := service.userManager.GetByEmail(c, email); err == nil && existing != nil {
		return &core.ServiceError{
			Code:    http.StatusBadRequest,
			Message: "validation failed",
			Fields:  map[string]string{"email": "already exists"},
		}
	}
	payload := RegisterPayload{
		Username: username,
		Email:    email,
		Password: password,
		Exp:      time.Now().Add(service.server.Cfg.PayloadTokenExpiresIn).Unix(),
	}

	token, err := core.GeneratePayloadToken(payload, service.server.Cfg.SecretKey)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	user := &users.User{
		Username: username,
		Email:    email,
	}

	err = SendRegistrationEmail(
		service.server.Email,
		user,
		service.server.Cfg.FrontURL,
		token,
	)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}

func (service *AuthService) Login(c *gin.Context, username, password string) (string, string, error) {
	user, err := service.userManager.GetByUsername(c, username)
	if err != nil || !user.CheckPassword(password, service.server.Cfg.SecretKey) {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid credentials",
		}
	}

	session, err := service.sessionManager.Create(c, user.ID)
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

	return accessToken, refreshToken, nil
}

func (service *AuthService) RefreshTokens(c *gin.Context, oldRefreshToken string) (string, string, error) {
	claims, err := ParseToken(oldRefreshToken, service.server.Cfg.SecretKey)
	if err != nil {
		return "", "", &core.ServiceError{
			Code:    http.StatusUnauthorized,
			Message: "invalid refresh token",
		}
	}

	session, err := service.sessionManager.GetByID(c, claims.SessionID)
	if err != nil || session == nil {
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

	return newAccess, newRefresh, nil
}

func (service *AuthService) Logout(c *gin.Context, token string) {
	if token == "" {
		return
	}

	claims, err := ParseToken(token, service.server.Cfg.SecretKey)
	if err == nil {
		service.sessionManager.Delete(c, claims.SessionID)
	}
}

func (service *AuthService) ForgotPassword(c *gin.Context, email string) error {
	user, err := service.userManager.GetByEmail(c, strings.ToLower(email))
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

	err = SendRegistrationEmail(
		service.server.Email,
		user,
		service.server.Cfg.FrontURL,
		token,
	)
	if err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
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

	user, err := service.userManager.GetByID(c, payload.UserID)
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

	if err := service.userManager.Save(c, user); err != nil {
		return &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return nil
}
