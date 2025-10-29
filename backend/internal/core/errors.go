package core

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

type ValidationError struct {
	Message string `json:"message"`
}

type FieldValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e ValidationError) Error() string {
	return e.Message
}

func (e FieldValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) *ValidationError {
	return &ValidationError{Message: message}
}

func NewFieldValidationError(field string, err error) *FieldValidationError {
	return &FieldValidationError{
		Field:   field,
		Message: err.Error(),
	}
}
