package auth

import (
	"app/internal/core"
	"app/internal/users"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware(server *core.Server) gin.HandlerFunc {
	userManager := users.NewUserManager(server)
	sessionManager := NewSessionManager(server)

	return func(c *gin.Context) {
		defer c.Next()

		tokenString, err := c.Cookie(server.Cfg.JWT.AccessCookie)
		if err != nil || tokenString == "" {
			c.Set("user", nil)
			return
		}

		claims, err := ParseToken(tokenString, server.Cfg.SecretKey)
		if err != nil {
			c.Set("user", nil)
			return
		}

		session, err := sessionManager.GetByID(claims.SessionID)
		if err != nil || session == nil {
			c.Set("user", nil)
			return
		}

		user, err := userManager.GetByID(claims.UserID)
		if err != nil {
			c.Set("user", nil)
			return
		}

		c.Set("user", user)
	}
}

func AuthRequiredMiddleware(server *core.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		el, ok := c.Get("user")
		if !ok {
			core.Fail(c, http.StatusInternalServerError, "server error", nil)
			log.Println("[ERROR] User not found in context; middleware must be applied after AuthenticationMiddleware")
			c.Abort()
			return
		}
		if el == nil {
			core.Fail(c, http.StatusUnauthorized, "user is not authenticated", nil)
			c.Abort()
			return
		}
		c.Next()
	}
}
