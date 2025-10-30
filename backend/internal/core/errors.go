package core

type ServiceError struct {
	Code    int
	Message string
	Fields  map[string]string
}

func (e *ServiceError) Error() string {
	return e.Message
}
