package core

import (
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterGroupRoutes(g *gin.RouterGroup, mdls []gin.HandlerFunc, s *Server) *gin.RouterGroup {
	prefix := "core"
	var group *gin.RouterGroup
	if g == nil {
		group = s.Eng.Group(prefix)
	} else {
		group = g.Group(prefix)
	}

	cr := NewCoreController(s)
	routes := []Route{
		{
			Method:                "GET",
			Path:                  "/healf",
			HandlerFuncs:          []gin.HandlerFunc{cr.Get},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "1",
		},
	}
	for _, route := range routes {
		allHandlers := append(route.DecoratorHandlerFuncs, mdls...)
		allHandlers = append(allHandlers, route.HandlerFuncs...)

		group.Handle(route.Method, route.Path, allHandlers...)
		if _, exists := s.RoutesMap[route.NameSpace]; exists {
			log.Printf("[WARN] Route namespace '%s' is repeated. Previous handler will be overwritten by the new one.\n", route.NameSpace)
		}
		s.RoutesMap[route.NameSpace] = route
	}
	return group
}
