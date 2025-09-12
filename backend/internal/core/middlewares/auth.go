package middlewares

import (
	"app/internal/controller"
	"app/internal/core"
	"app/internal/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthenticationMiddleware(server *core.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer c.Next()

		tokenString, err := c.Cookie("auth_token")
		if err != nil || tokenString == "" {
			c.Set("user", nil)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(server.Cfg.SECRET_KEY), nil
		})

		if err != nil || !token.Valid {
			c.Set("user", nil)
			return
		}

		claims, ok := token.Claims.(*controller.Claims)
		if !ok {
			c.Set("user", nil)
			return
		}
		var user models.User
		if err := server.DB.First(&user, claims.UserID).Error; err != nil {
			c.Set("user", nil)
			return
		}
		c.Set("user", &user)
	}
}

func AuthRequiredMiddleware(server *core.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		el, ok := c.Keys["user"]
		if !ok {
			c.JSON(http.StatusInternalServerError, struct{}{})
			c.Abort()
			log.Println("[WARN] Check if the user exists in the context; if not, either skip the handler (abort) or let it run after logging a warning.")
			return
		}
		if el == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not authenticated"})
			c.Abort()
			return
		}
		c.Next()
	}
}
