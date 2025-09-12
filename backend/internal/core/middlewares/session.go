package middlewares

import (
	"app/internal/core"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SessionMiddleware(server *core.Server) gin.HandlerFunc {
	store := cookie.NewStore([]byte(server.Cfg.SECRET_KEY))
	return sessions.Sessions("session_id", store)
}
