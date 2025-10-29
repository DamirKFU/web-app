package users

import (
	"app/internal/core"
)

func RegisterValidators(s *core.Server) {
	aliases := map[string]string{
		"username": "required,min=3,max=32,alphanum",
		"password": "required,min=6,max=64",
	}
	s.RegisterValidators(nil, aliases)
}
