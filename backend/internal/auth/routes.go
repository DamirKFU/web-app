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
			DecoratorHandlerFuncs: []gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			NameSpace:             "3",
		},
		{
			Method:                "POST",
			Path:                  "/register",
			HandlerFunc:           au.Register,
			DecoratorHandlerFuncs: []gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			NameSpace:             "4",
		},
		{
			Method:                "POST",
			Path:                  "/logout",
			HandlerFunc:           au.Logout,
			DecoratorHandlerFuncs: []gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			NameSpace:             "5",
		},
		{
			Method:                "POST",
			Path:                  "/refresh",
			HandlerFunc:           au.RefreshToken,
			DecoratorHandlerFuncs: []gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			NameSpace:             "6",
		},
		{
			Method:                "POST",
			Path:                  "/forgot",
			HandlerFunc:           au.ForgotPassword,
			DecoratorHandlerFuncs: []gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			NameSpace:             "7",
		},
		{
			Method:                "POST",
			Path:                  "/reset",
			HandlerFunc:           au.ResetPassword,
			DecoratorHandlerFuncs: []gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			NameSpace:             "8",
		},
	}
	return routes

}
