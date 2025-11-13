package users

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

func RegisterApp(g *gin.RouterGroup, s *core.Server) {
	models := []any{
		&User{},
	}
	RegisterValidators(s)
	s.RegisterModels(models...)
}
