package middlewares

import (
	"app/internal/core"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(server *core.Server) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     server.Cfg.AllowOrigins,
		AllowMethods:     server.Cfg.AllowMethods,
		AllowHeaders:     server.Cfg.AllowHeaders,
		AllowCredentials: server.Cfg.AllowCredentials,
	})
}
