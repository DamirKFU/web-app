package catalog

import "app/internal/core"

type ColorManager struct {
	server *core.Server
}

type CategoryManager struct {
	server *core.Server
}

type TShirtManager struct {
	server *core.Server
}

func NewColorManager(server *core.Server) *ColorManager {
	return &ColorManager{server: server}
}

func NewCategoryManager(server *core.Server) *CategoryManager {
	return &CategoryManager{server: server}
}

func NewTShirtManager(server *core.Server) *TShirtManager {
	return &TShirtManager{server: server}
}

func (m *ColorManager) GetAll() ([]Color, error) {
	var colors []Color
	if err := m.server.DB.Find(&colors).Error; err != nil {
		return nil, err
	}
	return colors, nil
}

func (m *CategoryManager) GetAll() ([]Category, error) {
	var categories []Category
	if err := m.server.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (m *TShirtManager) GetAll() ([]Garment, error) {
	var tshirts []Garment
	if err := m.server.DB.
		Preload("Category").
		Preload("Color").
		Find(&tshirts).Error; err != nil {
		return nil, err
	}
	return tshirts, nil
}
