package middlewares

import (
	"app/internal/core"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(server *core.Server) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     server.Cfg.CORS.AllowOrigins,
		AllowMethods:     server.Cfg.CORS.AllowMethods,
		AllowHeaders:     server.Cfg.CORS.AllowHeaders,
		AllowCredentials: server.Cfg.CORS.AllowCredentials,
	})
}
