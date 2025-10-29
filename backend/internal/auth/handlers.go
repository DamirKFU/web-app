package auth

import (
	"app/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	server  *core.Server
	service *AuthService
}

func NewAuthController(server *core.Server) *AuthController {
	return &AuthController{
		server:  server,
		service: NewAuthService(server),
	}
}

func (ctrl *AuthController) Register(c *gin.Context) {
	var body RegisterRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": core.ParseValidationError(err)})
		return
	}

	err := ctrl.service.Register(body.Username, body.Password)
	if core.HandleError(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "user registered"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var body LoginRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": core.ParseValidationError(err)})
		return
	}

	refreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie)
	if err != nil && refreshToken != "" {
		ctrl.server.RedisServer.RDB0.Del(c.Request.Context(), refreshToken).Err()
	}

	accessToken, refreshToken, err := ctrl.service.Login(c, body.Username, body.Password)
	if core.HandleError(c, err) {
		return
	}

	c.SetCookie(ctrl.server.Cfg.JWT.AccessCookie, accessToken, int(ctrl.server.Cfg.JWT.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(ctrl.server.Cfg.JWT.RefreshCookie, refreshToken, int(ctrl.server.Cfg.JWT.RefreshExpiresIn), "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"status": "logged in"})
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	newAccess, newRefresh, err := ctrl.service.RefreshTokens(c, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (ctrl *AuthController) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie)
	if err != nil && refreshToken != "" {
		ctrl.server.RedisServer.RDB0.Del(c.Request.Context(), refreshToken).Err()
	}
	c.SetCookie(ctrl.server.Cfg.JWT.AccessCookie, "", -1, "/", "", false, true)
	c.SetCookie(ctrl.server.Cfg.JWT.RefreshCookie, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}
