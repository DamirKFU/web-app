package users

import (
	"app/internal/core"
	"log"

	"github.com/gin-gonic/gin"
)

func RegisterApp(g *gin.RouterGroup, s *core.Server) {
	models := []any{
		&User{},
	}
	if err := s.DB.AutoMigrate(models...); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
	RegisterValidators(s)
}
