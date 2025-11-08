package users

import (
	"app/internal/core"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func RegisterValidators(s *core.Server) {
	aliases := map[string]string{
		"UsersUsername":       "required,min=3,max=42,username_msx",
		"UsersPassword":       "required,max=128",
		"DtoUsersPassword":    "required,min=6,max=56",
		"UsersEmail":          "required,email,max=128",
		"RepeatUsersPassword": "required,eqfield=Password",
	}

	s.RegisterValidators(map[string]validator.Func{
		"username_msx": ValidateUsernameMSX,
	}, aliases)
}

func ValidateUsernameMSX(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9._-]+$`, username)
	return matched
}
