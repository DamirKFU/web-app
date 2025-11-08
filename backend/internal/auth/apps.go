package auth

import (
	"app/internal/core"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterApp(g *gin.RouterGroup, s *core.Server) {
	models := []any{
		&Session{},
	}
	prefix := "auth"
	group := g.Group(prefix)
	routes := GetRoutes(s)
	s.RegisterRoutes(group, routes)
	if err := s.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}
