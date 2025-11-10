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
	err := c.ShouldBindJSON(&body)
	if core.HandleValidationError(c, err) {
		return
	}

	if core.HandleServiceError(c, ctrl.service.Register(body.Username, body.Password, body.Email)) {
		return
	}
	core.Success(c, gin.H{"message": "data is valid"})
}

func (ctrl *AuthController) RegisterConfirm(c *gin.Context) {
	var body RegisterConfirmRequest
	err := c.ShouldBindJSON(&body)
	if core.HandleValidationError(c, err) {
		return
	}

	if core.HandleServiceError(c, ctrl.service.RegisterConfirm(body.Token)) {
		return
	}
	core.SuccessWithStatus(c, http.StatusCreated, gin.H{"message": "user registered"})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var body LoginRequest
	err := c.ShouldBindJSON(&body)
	if core.HandleValidationError(c, err) {
		return
	}

	accessToken, refreshToken, err := ctrl.service.Login(c, body.Username, body.Password)
	if core.HandleServiceError(c, err) {
		return
	}

	token, err := core.SetCsrfToken(
		c,
		ctrl.server.Cfg.SecretKey,
		ctrl.server.Cfg.CSRF.Cookie,
		int(ctrl.server.Cfg.JWT.RefreshExpiresIn),
	)
	if err != nil {
		core.Fail(c, http.StatusInternalServerError, err.Error(), nil)
	}

	c.SetCookie(ctrl.server.Cfg.JWT.AccessCookie, accessToken, int(ctrl.server.Cfg.JWT.AccessExpiresIn), "/", "", false, true)
	c.SetCookie(ctrl.server.Cfg.JWT.RefreshCookie, refreshToken, int(ctrl.server.Cfg.JWT.RefreshExpiresIn), "/", "", false, true)
	core.Success(c, gin.H{
		"message":                   "logged in",
		ctrl.server.Cfg.CSRF.Cookie: token,
	})
}

func (ctrl *AuthController) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie)
	if err != nil || refreshToken == "" {
		core.Fail(c, http.StatusUnauthorized, "refresh token not found", nil)
		return
	}

	newAccess, newRefresh, err := ctrl.service.RefreshTokens(c, refreshToken)
	if core.HandleServiceError(c, err) {
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

	core.SetCsrfToken(
		c,
		ctrl.server.Cfg.SecretKey,
		ctrl.server.Cfg.CSRF.Cookie,
		int(ctrl.server.Cfg.JWT.RefreshExpiresIn),
	)

	core.Success(c, gin.H{"message": "access token refreshed"})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	if oldRefreshToken, err := c.Cookie(ctrl.server.Cfg.JWT.RefreshCookie); err == nil {
		ctrl.service.Logout(c, oldRefreshToken)
	}
	c.SetCookie(ctrl.server.Cfg.JWT.AccessCookie, "", -1, "/", "", false, true)
	c.SetCookie(ctrl.server.Cfg.JWT.RefreshCookie, "", -1, "/", "", false, true)

	core.Success(c, gin.H{"message": "logged out"})
}

func (ctrl *AuthController) ForgotPassword(c *gin.Context) {
	var body ForgotPasswordRequest
	err := c.ShouldBindJSON(&body)
	if core.HandleValidationError(c, err) {
		return
	}

	if core.HandleServiceError(c, ctrl.service.ForgotPassword(c, body.Email)) {
		return
	}

	core.Success(c, gin.H{"status": "ok"})
}

func (ctrl *AuthController) ResetPassword(c *gin.Context) {
	var body ResetPasswordRequest
	err := c.ShouldBindJSON(&body)
	if core.HandleValidationError(c, err) {
		return
	}

	if core.HandleServiceError(c, ctrl.service.ResetPassword(c, body.Token, body.Password)) {
		return
	}
	core.Success(c, gin.H{"status": "password_reset"})
}
