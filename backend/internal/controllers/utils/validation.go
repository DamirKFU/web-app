package utils

import "github.com/go-playground/validator/v10"

func ParseValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			return e.Field() + ": " + e.Tag()
		}
	}
	return err.Error()
}
