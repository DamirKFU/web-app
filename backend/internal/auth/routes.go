package auth

import (
	"app/internal/core"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRoutes(s *core.Server) []core.Route {
	au := NewAuthController(s)
	register_rate_limit := core.RateLimiterMiddleware(s, "general", 20, 60*time.Second)
	csrf_exempt := core.CsrfExemptMiddleware(s)
	routes := []core.Route{
		{
			Method:                "POST",
			Path:                  "/login",
			HandlerFuncs:          []gin.HandlerFunc{au.Login},
			DecoratorHandlerFuncs: []gin.HandlerFunc{csrf_exempt},
			NameSpace:             "login",
		},
		{
			Method:                "POST",
			Path:                  "/register-confirm",
			HandlerFuncs:          []gin.HandlerFunc{au.RegisterConfirm},
			DecoratorHandlerFuncs: []gin.HandlerFunc{csrf_exempt},
			NameSpace:             "register_confirm",
		},
		{
			Method:                "POST",
			Path:                  "/register",
			HandlerFuncs:          []gin.HandlerFunc{au.Register},
			DecoratorHandlerFuncs: []gin.HandlerFunc{register_rate_limit, csrf_exempt},
			NameSpace:             "register",
		},
		{
			Method:                "POST",
			Path:                  "/logout",
			HandlerFuncs:          []gin.HandlerFunc{au.Logout},
			DecoratorHandlerFuncs: []gin.HandlerFunc{csrf_exempt},
			NameSpace:             "logout",
		},
		{
			Method:                "POST",
			Path:                  "/refresh",
			HandlerFuncs:          []gin.HandlerFunc{au.RefreshToken},
			DecoratorHandlerFuncs: []gin.HandlerFunc{csrf_exempt},
			NameSpace:             "refresh",
		},
		{
			Method:                "POST",
			Path:                  "/forgot",
			HandlerFuncs:          []gin.HandlerFunc{au.ForgotPassword},
			DecoratorHandlerFuncs: []gin.HandlerFunc{csrf_exempt},
			NameSpace:             "forgot",
		},
		{
			Method:                "POST",
			Path:                  "/reset",
			HandlerFuncs:          []gin.HandlerFunc{au.ResetPassword},
			DecoratorHandlerFuncs: []gin.HandlerFunc{csrf_exempt},
			NameSpace:             "reset",
		},
	}
	return routes

}
