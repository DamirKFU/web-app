package catalog

import (
	"app/internal/core"
)

const (
	SizeS  = "S"
	SizeM  = "M"
	SizeL  = "L"
	SizeXL = "XL"
)

type Color struct {
	core.AbstractNameModel
	Color   string   `gorm:"type:varchar(7);unique;not null" json:"color"`
	TShirts []TShirt `gorm:"foreignKey:ColorID;constraint:OnDelete:CASCADE" json:"tshirts"`
}

type Category struct {
	core.AbstractNameModel
	TShirts []TShirt `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"tshirts"`
}

type TShirt struct {
	core.AbstractModel

	CategoryID uint      `gorm:"not null" json:"category_id"`
	Category   *Category `gorm:"constraint:OnDelete:CASCADE" json:"category"`

	ColorID uint   `gorm:"not null" json:"color_id"`
	Color   *Color `gorm:"constraint:OnDelete:CASCADE" json:"color"`

	Size string `gorm:"type:varchar(3);not null" json:"size"`

	Image string `gorm:"type:text;not null" json:"image"`

	Count uint `gorm:"default:0" json:"count"`
}
