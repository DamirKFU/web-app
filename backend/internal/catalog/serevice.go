package catalog

import (
	"app/internal/core"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CatalogService struct {
	server          *core.Server
	colorManager    *ColorManager
	categoryManager *CategoryManager
	tshirtManager   *TShirtManager
}

func NewCatalogService(server *core.Server) *CatalogService {
	return &CatalogService{
		server:          server,
		colorManager:    NewColorManager(server),
		categoryManager: NewCategoryManager(server),
		tshirtManager:   NewTShirtManager(server),
	}
}

func (s *CatalogService) GetColors(c *gin.Context) ([]Color, error) {
	colors, err := s.colorManager.GetAll(c)
	if err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return colors, nil
}

func (s *CatalogService) GetCategories(c *gin.Context) ([]Category, error) {
	categories, err := s.categoryManager.GetAll(c)
	if err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return categories, nil
}

func (s *CatalogService) GetGarments(c *gin.Context) ([]Garment, error) {
	tshirts, err := s.tshirtManager.GetAll(c)
	if err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return tshirts, nil
}
