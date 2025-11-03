package catalog

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

func GetRoutes(s *core.Server) []core.Route {
	au := NewCatalogController(s)
	routes := []core.Route{
		{
			Method:                "GET",
			Path:                  "/colors",
			HandlerFunc:           au.GetColors,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "colors",
		},
		{
			Method:                "GET",
			Path:                  "/categories",
			HandlerFunc:           au.GetCategories,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "categories",
		},
	}

	return routes

}
