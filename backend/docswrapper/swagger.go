package docswrapper

import (
	"github.com/gin-gonic/gin"
)

// LoginHandler godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth.LoginRequest true "Login credentials"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /api/v1/auth/login [post]
func LoginHandler(c *gin.Context) {}
