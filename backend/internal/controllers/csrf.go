package controllers

import (
	"app/internal/controllers/utils"
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

type CsrfController struct {
	server *core.Server
}

func NewCsrfController(server *core.Server) *AuthController {
	return &AuthController{server: server}
}

func (ctrl *CsrfController) GetToken(c *gin.Context) string {
	csrfCockie := ctrl.server.Cfg.CSRF.Cookie
	cookieToken, err := c.Cookie(ctrl.server.Cfg.CSRF.Cookie)
	var token string

	if err != nil || cookieToken == "" {
		token = utils.GenerateCSRFToken(ctrl.server.Cfg.SecretKey)

		c.SetCookie(csrfCockie,
			token,
			0,
			"/",
			"",
			false,
			true,
		)
	} else {
		token = cookieToken
	}

	return token
}
