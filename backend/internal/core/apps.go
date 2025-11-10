package core

import (
	"github.com/gin-gonic/gin"
)

func RegisterApp(g *gin.RouterGroup, s *Server) {
	prefix := "core"
	group := g.Group(prefix)
	routes := GetRoutes(s)
	RegisterValidators(s)
	s.RegisterRoutes(group, routes)
}
