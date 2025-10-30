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

	models := []any{
		&catalog.Color{},
		&catalog.Category{},
		&catalog.TShirt{},
		&users.User{},
	}
	base_mdls := []gin.HandlerFunc{
		core.CorsMiddleware(s),
		core.CSRFMiddleware(s),
		auth.AuthenticationMiddleware(s),
	}

	base_group := s.Eng.Group("api/v1/")
	auth.RegisterGroupRoutes(
		base_group,
		append(
			[]gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			base_mdls...,
		),
		s,
	)
	core.RegisterGroupRoutes(base_group, base_mdls, s)
	s.Eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := s.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
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
