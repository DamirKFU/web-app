package catalog

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

func RegisterGroupRoutes(g *gin.RouterGroup, mdls []gin.HandlerFunc, s *core.Server) *gin.RouterGroup {
	prefix := "catalog"
	var group *gin.RouterGroup
	if g == nil {
		group = s.Eng.Group(prefix)
	} else {
		group = g.Group(prefix)
	}

	au := NewCatalogController(s)
	routes := []core.Route{
		{
			Method:                "GET",
			Path:                  "/colors",
			HandlerFuncs:          []gin.HandlerFunc{au.GetColors},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "colors",
		},
		{
			Method:                "GET",
			Path:                  "/categories",
			HandlerFuncs:          []gin.HandlerFunc{au.GetCategories},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "colors",
		},
	}
	s.RegisterRoutes(group, routes, mdls)
	return group
}
