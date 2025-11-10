package catalog

import (
	"app/internal/core"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(s *core.Server) {
	validators := map[string]validator.Func{
		"hexcolor":  ValidateColorHex,
		"size_enum": ValidateSizeGarment,
	}

	aliases := map[string]string{
		"CatalogColor":      "required,hexcolor",
		"CatalogImage":      "required,url",
		"CatalogSizeEnum":   "required,size_enum",
		"GarmentColorID":    "required,id",
		"GarmentCategoryID": "required,id",
	}

	s.RegisterValidators(validators, aliases)
}

func ValidateColorHex(fl validator.FieldLevel) bool {
	var colorRegex = regexp.MustCompile(`^#([a-f0-9]{6})$`)

	return colorRegex.MatchString(fl.Field().String())
}

func ValidateSizeGarment(fl validator.FieldLevel) bool {
	size := fl.Field().String()
	switch size {
	case SizeS, SizeM, SizeL, SizeXL:
		return true
	default:
		return false
	}
}
