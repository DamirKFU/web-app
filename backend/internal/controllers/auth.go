package controllers

import (
	"app/internal/controllers/dto"
	"app/internal/controllers/service"
	"app/internal/controllers/utils"
	"app/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	server    *core.Server
	validator *validator.Validate
	service   *service.AuthService
}

func NewAuthController(server *core.Server) *AuthController {
	return &AuthController{
		server:    server,
		validator: validator.New(),
		service:   service.NewAuthService(server),
	}
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	claims, err := ctrl.service.RefreshTokens(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	newAccess, newRefresh, err := utils.GenerateTokens(
		claims.UserID,
		ctrl.server.Cfg.SecretKey,
		ctrl.server.Cfg.JWT.AccessExpiresIn,
		ctrl.server.Cfg.JWT.RefreshExpiresIn,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate access token"})
		return
	}

	_ = ctrl.server.RedisServer.RDB0.Del(c.Request.Context(), refreshToken).Err() // удаляем старый

	err = ctrl.server.RedisServer.RDB0.Set(
		c.Request.Context(),
		newRefresh,
		claims.UserID,
		ctrl.server.Cfg.JWT.RefreshExpiresIn,
	).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save new refresh token"})
		return
	}

	c.SetCookie(
		ctrl.server.Cfg.JWT.AccessCookie,
		newAccess,
		int(ctrl.server.Cfg.JWT.AccessExpiresIn),
		"/",
		"",
		false,
		true,
	)

	c.SetCookie(
		ctrl.server.Cfg.JWT.RefreshCookie,
		newRefresh,
		int(ctrl.server.Cfg.JWT.RefreshExpiresIn),
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{"status": "access token refreshed"})
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var body dto.RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := ctrl.validator.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ParseValidationError(err)})
		return
	}

	if user, err := ctrl.service.Register(body.Username, body.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		if err := ctrl.server.DB.Create(user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "user registered"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var body dto.LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := ctrl.validator.Struct(body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ParseValidationError(err)})
		return
	}

	userId, err := ctrl.service.Login(body.Username, body.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(
		userId, ctrl.server.Cfg.SecretKey,
		ctrl.server.Cfg.JWT.AccessExpiresIn,
		ctrl.server.Cfg.JWT.RefreshExpiresIn,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save refresh token"})
		return
	}

	err = ctrl.server.RedisServer.RDB0.Set(
		c.Request.Context(),
		refreshToken,
		userId,
		ctrl.server.Cfg.JWT.RefreshExpiresIn,
	).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate tokens"})
		return
	}

	c.SetCookie(ctrl.server.Cfg.JWT.AccessCookie, accessToken, int(ctrl.server.Cfg.JWT.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(ctrl.server.Cfg.JWT.RefreshCookie, refreshToken, int(ctrl.server.Cfg.JWT.RefreshExpiresIn), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "logged in"})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie)
	if err != nil && refreshToken != "" {
		_ = ctrl.server.RedisServer.RDB0.Del(c.Request.Context(), refreshToken).Err()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete token"})
	}
	c.SetCookie(ctrl.server.Cfg.JWT.AccessCookie, "", -1, "/", "", false, true)
	c.SetCookie(ctrl.server.Cfg.JWT.RefreshCookie, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}
