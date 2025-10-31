package catalog

import "app/internal/core"

type ColorResponse struct {
	core.AbstractNameModelResponce
	Color string `json:"color"`
}

type CategoryResponse struct {
	core.AbstractNameModelResponce
}
