package auth

import "app/internal/core"

func RegisterValidators(s *core.Server) {

	aliases := map[string]string{
		"SessionsUserID": "required,id",
	}

	s.RegisterValidators(nil, aliases)

}
