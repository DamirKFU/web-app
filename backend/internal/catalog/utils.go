package catalog

import (
	"app/internal/core"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(s *core.Server) {
	validators := map[string]validator.Func{
		"hexcolor": ValidateColorHex,
	}
	s.RegisterValidators(validators, nil)
}
