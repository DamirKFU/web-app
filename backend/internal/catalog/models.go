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
	Color string `gorm:"unique;not null;size:7"`
}
