package core

import (
	"github.com/gin-gonic/gin"
)

type CoreController struct {
	server *Server
}

func NewCoreController(server *Server) *CoreController {
	return &CoreController{server: server}
}

func (ctrl *CoreController) Get(c *gin.Context) {
	user := c.Keys["user"]
	c.JSON(200, gin.H{"status": "ok", "user": user})
}
