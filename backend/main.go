package main

import (
	"log"

	"app/config"
	"app/internal/auth"
	"app/internal/catalog"
	"app/internal/core"
	"app/internal/users"

	"github.com/gin-gonic/gin"
)

func CreateApp(cfg config.Config) *core.Server {

	s := core.NewServer(cfg)

	catalog.RegisterValidators(s)
	users.RegisterValidators(s)

	models := []interface{}{
		&catalog.Color{},
		&users.User{},
	}
	base_mdls := []gin.HandlerFunc{
		core.CorsMiddleware(s),
		core.CSRFMiddleware(s),
		auth.AuthenticationMiddleware(s),
	}

	base_group := s.Eng.Group("")
	auth.RegisterGroupRoutes(
		base_group,
		append(
			[]gin.HandlerFunc{core.CsrfExemptMiddleware(s)},
			base_mdls...,
		),
		s,
	)
	core.RegisterGroupRoutes(base_group, base_mdls, s)

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
