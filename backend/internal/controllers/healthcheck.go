package controllers

import (
	"app/internal/core"

	"github.com/gin-gonic/gin"
)

type HealthcheckController struct {
	server *core.Server
}

func NewHealthcheckController(server *core.Server) *HealthcheckController {
	return &HealthcheckController{server: server}
}

func (ctrl *HealthcheckController) Get(ctx *gin.Context) {
	user := ctx.Keys["user"]
	ctx.JSON(200, gin.H{"status": "ok", "user": user})
}
