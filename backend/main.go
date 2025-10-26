package main

import (
	"log"

	"app/config"
	"app/internal/controllers"
	"app/internal/core"
	"app/internal/core/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	s := core.NewServer(cfg)

	var hc = controllers.NewHealthcheckController(s)
	var au = controllers.NewAuthController(s)

	routes := []core.Route{
		{Method: "GET", Path: "/healf", HandlerFunc: hc.Get, NameSpace: "1"},
		{Method: "POST", Path: "/login", HandlerFunc: au.Login, NameSpace: "3"},
		{Method: "POST", Path: "/register", HandlerFunc: au.Register, NameSpace: "4"},
		{Method: "POST", Path: "/logout", HandlerFunc: au.Logout, NameSpace: "5"},
	}
	mdls := []gin.HandlerFunc{
		middlewares.CorsMiddleware(s),
		middlewares.CSRFMiddleware(s),
		middlewares.AuthenticationMiddleware(s),
	}
	s.RegisterMiddlewares(mdls)
	s.RegisterRoutes(routes)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
