package catalog

import (
	"app/internal/core"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(s *core.Server) {
	validators := map[string]validator.Func{
		"hexcolor": ValidateColorHex,
		"size_enum": func(fl validator.FieldLevel) bool {
			size := fl.Field().String()
			switch size {
			case SizeS, SizeM, SizeL, SizeXL:
				return true
			default:
				return false
			}
		},
		"foreign_key_id": func(fl validator.FieldLevel) bool {
			id, ok := fl.Field().Interface().(uint)
			return ok && id > 0
		},
	}

	aliases := map[string]string{
		"CatalogColor":      "required,hexcolor",
		"CatalogName":       "required,min=2,max=64",
		"CatalogImage":      "required,url",
		"CatalogSizeEnum":   "required,size_enum",
		"GarmentColorID":    "required,foreign_key_id",
		"GarmentCategoryID": "required,foreign_key_id",
	}

	s.RegisterValidators(validators, aliases)
}

func ValidateColorHex(fl validator.FieldLevel) bool {
	var colorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`)

	return colorRegex.MatchString(fl.Field().String())
}
