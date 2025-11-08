package catalog

import (
	"app/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type CatalogController struct {
	server  *core.Server
	service *CatalogService
}

func NewCatalogController(server *core.Server) *CatalogController {
	return &CatalogController{
		server:  server,
		service: NewCatalogService(server),
	}
}

func (ctrl *CatalogController) GetColors(c *gin.Context) {
	colors, err := ctrl.service.GetColors()
	if core.HandleServiceError(c, err) {
		return
	}

	var response []ColorResponse
	if err := copier.Copy(&response, &colors); err != nil {
		core.Fail(c, http.StatusInternalServerError, "Failed to process data", nil)
		return
	}

	core.Success(c, response)
}

func (ctrl *CatalogController) GetCategories(c *gin.Context) {
	categories, err := ctrl.service.GetCategories()
	if core.HandleServiceError(c, err) {
		return
	}

	var response []CategoryResponse
	if err := copier.Copy(&response, &categories); err != nil {
		core.Fail(c, http.StatusInternalServerError, "Failed to process data", nil)
		return
	}

	core.Success(c, response)
}

func (ctrl *CatalogController) GetGarments(c *gin.Context) {
	garments, err := ctrl.service.GetGarments()
	if core.HandleServiceError(c, err) {
		return
	}

	var response []GarmentResponse
	if err := copier.Copy(&response, &garments); err != nil {
		core.Fail(c, http.StatusInternalServerError, "Failed to process data", nil)
		return
	}

	core.Success(c, response)
}
