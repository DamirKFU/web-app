package auth

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

func GetRoutes(s *core.Server) []core.Route {
	au := NewAuthController(s)
	routes := []core.Route{
		{
			Method:                "POST",
			Path:                  "/login",
			HandlerFunc:           au.Login,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "3",
		},
		{
			Method:                "POST",
			Path:                  "/register",
			HandlerFunc:           au.Register,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "4",
		},
		{
			Method:                "POST",
			Path:                  "/logout",
			HandlerFunc:           au.Logout,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "5",
		},
		{
			Method:                "POST",
			Path:                  "/refresh",
			HandlerFunc:           au.RefreshToken,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "6",
		},
	}
	return routes

}
