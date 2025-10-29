package auth

import (
	"app/internal/core"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterGroupRoutes(g *gin.RouterGroup, mdls []gin.HandlerFunc, s *core.Server) *gin.RouterGroup {
	prefix := "auth"
	var group *gin.RouterGroup
	if g == nil {
		group = s.Eng.Group(prefix)
	} else {
		group = g.Group(prefix)
	}

	au := NewAuthController(s)
	routes := []core.Route{
		{
			Method:                "POST",
			Path:                  "/login",
			HandlerFuncs:          []gin.HandlerFunc{au.Login},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "3",
		},
		{
			Method:                "POST",
			Path:                  "/register",
			HandlerFuncs:          []gin.HandlerFunc{au.Register},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "4",
		},
		{
			Method:                "POST",
			Path:                  "/logout",
			HandlerFuncs:          []gin.HandlerFunc{au.Logout},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "5",
		},
		{
			Method:                "POST",
			Path:                  "/refresh",
			HandlerFuncs:          []gin.HandlerFunc{au.RefreshToken},
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "6",
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
