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
			Path:                  "colors/",
			HandlerFuncs:          []gin.HandlerFunc{au.GetColors},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "colors",
		},
		{
			Method:                "GET",
			Path:                  "categories/",
			HandlerFuncs:          []gin.HandlerFunc{au.GetCategories},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "categories",
		},
		{
			Method:                "GET",
			Path:                  "garments/",
			HandlerFuncs:          []gin.HandlerFunc{au.GetGarments},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "garments",
		},
	}

	return routes

}
