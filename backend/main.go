package main

import (
	"log"

	"app/config"
	"app/internal/auth"
	"app/internal/catalog"
	"app/internal/core"
	"app/internal/users"

	_ "app/docs" // <- обязательно, чтобы swag зарегистрировал JSON

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
		auth.AuthenticationMiddleware(s),
		core.CSRFMiddleware(s),
	}

	s.RegisterMiddlewares(base_mdls)

	base_group := s.Eng.Group("api/v1/")

	auth.RegisterApp(base_group, s)
	core.RegisterApp(base_group, s)
	catalog.RegisterApp(base_group, s)

	s.Eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return s
}

func main() {
	cfg := config.Load()
	s := CreateApp(cfg)
	log.Println(s.RoutesMap)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
