package catalog

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

func RegisterApp(g *gin.RouterGroup, s *core.Server) {
	models := []any{
		&Color{},
		&Category{},
		&Garment{},
	}
	prefix := "catalog/"
	group := g.Group(prefix)
	routes := GetRoutes(s)
	RegisterValidators(s)
	s.RegisterRoutes(group, routes)
	s.RegisterModels(models...)
}
