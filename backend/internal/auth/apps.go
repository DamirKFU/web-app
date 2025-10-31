package auth

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

func RegisterApp(g *gin.RouterGroup, s *core.Server) {
	prefix := "auth"
	group := g.Group(prefix)
	routes := GetRoutes(s)
	s.RegisterRoutes(group, routes)
}
