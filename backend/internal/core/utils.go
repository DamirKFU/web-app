package core

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"slices"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GenerateCSRFToken(secret string) (string, error) {
	randomPart := uniuri.NewLen(16)

	h := sha1.New()
	_, err := io.WriteString(h, randomPart+"-"+secret)
	if err != nil {
		return "", err
	}
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash, nil
}

func ParseValidationError(err error) map[string]string {
	fieldErrors := make(map[string]string)
	for _, fe := range err.(validator.ValidationErrors) {
		fieldErrors[fe.Field()] = fe.Error()
	}
	return fieldErrors
}

func HandleServiceError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	if se, ok := err.(*ServiceError); ok {
		Fail(c, se.Code, se.Message, se.Fields)
	} else {
		Fail(c, http.StatusInternalServerError, "internal server error", nil)
	}
	return true
}

func JSONResponse(c *gin.Context, code int, success bool, data any, apiErr *APIError) {
	c.JSON(code, APIResponse{
		Success: success,
		Data:    data,
		Error:   apiErr,
	})
}

func Fail(c *gin.Context, code int, message string, fields map[string]string) {
	JSONResponse(c, code, false, nil, &APIError{
		Message: message,
		Fields:  fields,
	})
}

func Success(c *gin.Context, data any) {
	JSONResponse(c, http.StatusOK, true, data, nil)
}

func inArray(arr []string, value string) bool {
	inarr := slices.Contains(arr, value)

	return inarr
}

func CheckCsrfExempt(c *gin.Context) (exempt bool, valid bool) {
	value, exists := c.Get("csrf_exempt")
	if !exists {
		return false, false
	}

	csrfExempt, ok := value.(bool)
	if !ok {
		return false, false
	}

	return csrfExempt, true
}
