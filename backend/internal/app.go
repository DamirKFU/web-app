package internal

import (
	"app/config"
	_ "app/docs"
	"app/internal/auth"
	"app/internal/catalog"
	"app/internal/core"
	"app/internal/users"
	"encoding/json"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateApp(cfg config.Config) *core.Server {
	if gin.Mode() != gin.TestMode && !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	if cfg.Debug {
		data, _ := json.MarshalIndent(cfg, "", "  ")
		log.Println(string(data))
	}
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

	core.RegisterApp(base_group, s)
	users.RegisterApp(base_group, s)
	auth.RegisterApp(base_group, s)
	catalog.RegisterApp(base_group, s)

	s.Eng.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return s
}
