package catalog

import "app/internal/core"

type ColorManager struct {
	server *core.Server
}

func NewColorManager(server *core.Server) *ColorManager {
	return &ColorManager{server: server}
}

func (m *ColorManager) GetAll() ([]Color, error) {
	var colors []Color
	if err := m.server.DB.Find(&colors).Error; err != nil {
		return nil, err
	}
	return colors, nil
}

type CategoryManager struct {
	server *core.Server
}

func NewCategoryManager(server *core.Server) *CategoryManager {
	return &CategoryManager{server: server}
}

func (m *CategoryManager) GetAll() ([]Category, error) {
	var categories []Category
	if err := m.server.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}
