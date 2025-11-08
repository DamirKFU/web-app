package core

import (
	"log"
	"net/http"

	"app/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func NewServer(cfg config.Config) *Server {
	return &Server{
		Cfg:         cfg,
		DB:          NewDB(cfg),
		Eng:         gin.Default(),
		RoutesMap:   make(map[string]Route),
		RedisServer: NewRedisServer(cfg),
		Mdls:        make([]gin.HandlerFunc, 0),
		Email:       NewEmailSMTPEngine(cfg),
	}
}

func (s *Server) RegisterValidators(
	customValidators map[string]validator.Func,
	aliases map[string]string,
) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		for name, fn := range customValidators {
			if err := v.RegisterValidation(name, fn); err != nil {
				panic("failed to register validator: " + name)
			}
		}

		for alias, rule := range aliases {
			v.RegisterAlias(alias, rule)
		}
	}
}

func (s *Server) RegisterRoutes(group *gin.RouterGroup, routes []Route) {
	for _, route := range routes {
		orderHandlers := append(route.DecoratorHandlerFuncs, s.Mdls...)
		orderHandlers = append(orderHandlers, route.HandlerFuncs...)
		group.Handle(route.Method, route.Path, orderHandlers...)
		if _, exists := s.RoutesMap[route.NameSpace]; exists {
			log.Printf("[WARN] Route namespace '%s' is repeated. Previous handler will be overwritten by the new one.\n", route.NameSpace)
		}
		s.RoutesMap[route.NameSpace] = route
	}
}

func (s *Server) RegisterMiddlewares(mdls []gin.HandlerFunc) {
	s.Mdls = append(s.Mdls, mdls...)
}

func (s *Server) Start() error {
	srv := &http.Server{
		Addr:         s.Cfg.HTTP.Addr,
		Handler:      s.Eng,
		ReadTimeout:  s.Cfg.HTTP.ReadTimeout,
		WriteTimeout: s.Cfg.HTTP.WriteTimeout,
	}

	return srv.ListenAndServe()
}
