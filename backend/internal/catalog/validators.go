package catalog

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidateColorHex(fl validator.FieldLevel) bool {
	var colorRegex = regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`)

	return colorRegex.MatchString(fl.Field().String())
}
