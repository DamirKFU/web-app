package catalog

import (
	"app/internal/core"
	"net/http"
)

type CatalogService struct {
	server *core.Server
}

func NewCatalogService(server *core.Server) *CatalogService {
	return &CatalogService{server: server}
}

func (s *CatalogService) GetColors() ([]Color, error) {
	var colors []Color
	if err := s.server.DB.Find(&colors).Error; err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return colors, nil
}

func (s *CatalogService) GetCategories() ([]Category, error) {
	var categories []Category
	if err := s.server.DB.Find(&categories).Error; err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return categories, nil
}
