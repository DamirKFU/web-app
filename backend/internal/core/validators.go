package core

import (
	"github.com/go-playground/validator/v10"
)

func RegisterValidators(s *Server) {
	validators := map[string]validator.Func{
		"id": ValidateID,
	}

	aliases := map[string]string{
		"AbstractName": "required,max=150",
	}

	s.RegisterValidators(validators, aliases)

}

func ValidateID(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(uint)
	return ok && id > 0
}
