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
// @Success 200 {object} core.APIResponse{data=map[string]string} "Successfully logged in"
// @Failure 400 {object} core.APIResponse{error=core.APIError} "Bad request"
// @Failure 401 {object} core.APIResponse{error=core.APIError} "Unauthorized"
// @Router /api/v1/auth/login [post]
func LoginHandler(c *gin.Context) {}

// RegisterHandler godoc
// @Summary Register user
// @Tags auth
// @Accept json
// @Produce json
// @Param input body auth.RegisterRequest true "Register credentials"
// @Success 200 {object} core.APIResponse{data=map[string]string} "Successfully logged in"
// @Failure 400 {object} core.APIResponse{error=core.APIError} "Bad request"
// @Router /api/v1/auth/register [post]
func RegisterHandler(c *gin.Context) {}

// ColorHandler godoc
// @Summary Ger colors
// @Tags catalog
// @Produce json
// @Success 200 {object} core.APIResponse{data=[]catalog.ColorResponse} "Successfully get colors"
// @Failure 400 {object} core.APIResponse{error=core.APIError} "Bad request"
// @Router /api/v1/catalog/colors [get]
func ColorsHandler(c *gin.Context) {}

// CategoriesHandler godoc
// @Summary Ger categories
// @Tags catalog
// @Produce json
// @Success 200 {object} core.APIResponse{data=[]catalog.CategoryResponse} "Successfully get categories"
// @Failure 400 {object} core.APIResponse{error=core.APIError} "Bad request"
// @Router /api/v1/catalog/categories [get]
func CategoriesHandler(c *gin.Context) {}
