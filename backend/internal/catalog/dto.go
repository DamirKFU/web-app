package catalog

import "app/internal/core"

type ColorResponse struct {
	core.AbstractNameModelResponce
	Color string `json:"color"`
}

type CategoryResponse struct {
	core.AbstractNameModelResponce
}

type GarmentResponse struct {
	ID       uint              `json:"id"`
	Size     string            `json:"size"`
	Image    string            `json:"image"`
	Count    uint              `json:"count"`
	Category *CategoryResponse `json:"category"`
	Color    *ColorResponse    `json:"color"`
}
