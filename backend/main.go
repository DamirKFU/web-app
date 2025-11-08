package main

import (
	"log"
	"time"

	"app/config"
	"app/internal/auth"
	"app/internal/catalog"
	"app/internal/core"
	"app/internal/users"

	_ "app/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateApp(cfg config.Config) *core.Server {

	s := core.NewServer(cfg)

	catalog.RegisterValidators(s)
	users.RegisterValidators(s)

	base_mdls := []gin.HandlerFunc{
		core.CorsMiddleware(s),
		core.LimitRequestBodySizeMiddleware(s),
		core.XSSMiddleware(s),
		auth.AuthenticationMiddleware(s),
		core.CSRFMiddleware(s),
		core.RateLimiterMiddleware(s, "general", 20, 60*time.Second),
	}

	s.RegisterMiddlewares(base_mdls)

	base_group := s.Eng.Group("api/v1/")

	users.RegisterApp(base_group, s)
	auth.RegisterApp(base_group, s)
	core.RegisterApp(base_group, s)
	catalog.RegisterApp(base_group, s)

	s.Eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return s
}

func main() {
	cfg := config.Load()
	s := CreateApp(cfg)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
