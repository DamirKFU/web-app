package core

import (
	"log"
	"net/http"

	"app/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	Cfg       config.Config
	DB        *gorm.DB
	Eng       *gin.Engine
	RoutesMap map[string]Route
}

type Route struct {
	Method      string
	Path        string
	HandlerFunc gin.HandlerFunc
	NameSpace   string
}

type MiddlewareGroup struct {
	Prefix     string
	Middleware []gin.HandlerFunc
	Routes     []Route
}

func NewServer(cfg config.Config) *Server {
	return &Server{Cfg: cfg, DB: NewDB(), Eng: gin.New(), RoutesMap: make(map[string]Route)}
}

func (s *Server) RegisterRoutes(rs []Route) {
	for _, route := range rs {
		s.Eng.Handle(route.Method, route.Path, route.HandlerFunc)
		if _, exists := s.RoutesMap[route.NameSpace]; exists {
			log.Printf("[WARN] Route namespace '%s' is repeated. Previous handler will be overwritten by the new one.\n", route.NameSpace)
		}
		s.RoutesMap[route.NameSpace] = route
	}
}

func (s *Server) RegisterMiddlewares(middlewares []gin.HandlerFunc) {
	for _, middleware := range middlewares {
		s.Eng.Use(middleware)
	}
}

func (s *Server) RegisterMiddlewareGroups(mgs []MiddlewareGroup) {
	for _, mg := range mgs {
		group := s.Eng.Group(mg.Prefix)
		group.Use(mg.Middleware...)
		for _, route := range mg.Routes {
			group.Handle(route.Method, route.Path, route.HandlerFunc)
		}
	}
}

func (s *Server) Start() error {
	s.Eng.Use(gin.Logger(), gin.Recovery())
	srv := &http.Server{
		Addr:         s.Cfg.HTTP.Addr,
		Handler:      s.Eng,
		ReadTimeout:  s.Cfg.HTTP.ReadTimeout,
		WriteTimeout: s.Cfg.HTTP.WriteTimeout,
	}

	return srv.ListenAndServe()
}
