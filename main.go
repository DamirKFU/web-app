package main

import (
	"log"

	"app/config"
	"app/internal/controller"
	"app/internal/core"
	"app/internal/core/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	s := core.NewServer(cfg)

	var hc = controller.NewHealthcheckController(s)
	var us = controller.NewUsersController(s)

	routes := []core.Route{
		{Method: "GET", Path: "healf/", HandlerFunc: hc.Get, NameSpace: "fff"},
		{Method: "GET", Path: "users/", HandlerFunc: us.Create, NameSpace: "fff"},
	}
	s.RegisterRoutes(routes)

	mdlws := []gin.HandlerFunc{
		middlewares.SessionMiddleware(s),
	}
	s.RegisterMiddlewares(mdlws)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
