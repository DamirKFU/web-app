package core

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CoreController struct {
	server *Server
}

func NewCoreController(server *Server) *CoreController {
	return &CoreController{server: server}
}

func (ctrl *CoreController) GetToken(c *gin.Context) {
	csrfCookieName := ctrl.server.Cfg.CSRF.Cookie
	cookieToken, err := c.Cookie(csrfCookieName)
	var token string

	if err != nil || cookieToken == "" {
		token, err = GenerateCSRFToken(ctrl.server.Cfg.SecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to generate CSRF token",
			})
			return
		}

		c.SetCookie(
			csrfCookieName,
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

	c.JSON(http.StatusOK, gin.H{
		"csrf_token": token,
	})
}

func (ctrl *CoreController) Get(c *gin.Context) {
	user := c.Keys["user"]
	c.JSON(200, gin.H{"status": "ok", "user": user})
}
