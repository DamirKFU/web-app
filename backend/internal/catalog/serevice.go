package catalog

import (
	"app/internal/core"
	"net/http"
)

type CatalogService struct {
	server          *core.Server
	colorManager    *ColorManager
	categoryManager *CategoryManager
}

func NewCatalogService(server *core.Server) *CatalogService {
	return &CatalogService{
		server:          server,
		colorManager:    NewColorManager(server),
		categoryManager: NewCategoryManager(server),
	}
}

func (s *CatalogService) GetColors() ([]Color, error) {
	colors, err := s.colorManager.GetAll()
	if err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return colors, nil
}

func (s *CatalogService) GetCategories() ([]Category, error) {
	categories, err := s.categoryManager.GetAll()
	if err != nil {
		return nil, &core.ServiceError{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return categories, nil
}
