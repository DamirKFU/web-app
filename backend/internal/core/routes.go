package core

import (
	"github.com/gin-gonic/gin"
)

func GetRoutes(s *Server) []Route {
	cr := NewCoreController(s)
	routes := []Route{
		{
			Method:                "GET",
			Path:                  "/healf",
			HandlerFunc:           cr.Get,
			DecoratorHandlerFuncs: []gin.HandlerFunc{},
			NameSpace:             "1",
		},
	}

	return routes
}
