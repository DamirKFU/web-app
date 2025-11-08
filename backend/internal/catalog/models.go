package catalog

import (
	"app/internal/core"

	"gorm.io/gorm"
)

type Color struct {
	core.AbstractNameModel
	Color    string    `gorm:"type:varchar(7);unique;not null" binding:"CatalogColor"`
	Garments []Garment `gorm:"foreignKey:ColorID;constraint:OnDelete:CASCADE"`
}

type Category struct {
	core.AbstractNameModel
	Garments []Garment `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE"`
}

type Garment struct {
	core.AbstractModel

	CategoryID uint      `gorm:"not null" binding:"GarmentCategoryID"`
	Category   *Category `gorm:"constraint:OnDelete:CASCADE"`

	ColorID uint   `gorm:"not null" binding:"GarmentColorID"`
	Color   *Color `gorm:"constraint:OnDelete:CASCADE"`

	Size string `gorm:"type:varchar(3);not null" binding:"CatalogSizeEnum"`

	Image string `gorm:"type:text;not null" binding:"CatalogImage"`

	Count uint `gorm:"default:0"`
}

func (c *Color) BeforeSave(tx *gorm.DB) error {
	return core.ValidateStruct(c)
}

func (c *Category) BeforeSave(tx *gorm.DB) error {
	return core.ValidateStruct(c)
}

func (g *Garment) BeforeSave(tx *gorm.DB) error {
	return core.ValidateStruct(g)
}
