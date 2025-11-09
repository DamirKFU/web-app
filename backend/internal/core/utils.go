package core

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"slices"

	"github.com/dchest/uniuri"
	"github.com/fernet/fernet-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
		if fe.Param() != "" {
			fieldErrors[fe.Field()] = fmt.Sprintf("%v=%v", fe.ActualTag(), fe.Param())
		} else {
			fieldErrors[fe.Field()] = fmt.Sprintf("%v", fe.ActualTag())
		}
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

func SuccessWithStatus(c *gin.Context, status int, data any) {
	JSONResponse(c, status, true, data, nil)
}

func inArray(arr []string, value string) bool {
	inarr := slices.Contains(arr, value)

	return inarr
}

func SetCsrfToken(c *gin.Context, secretKey, csrfCookieName string) (string, error) {
	cookieToken, err := c.Cookie(csrfCookieName)
	var token string

	if err != nil || cookieToken == "" {
		token, err = GenerateCSRFToken(secretKey)
		if err != nil {
			return "", errors.New("failed to generate CSRF token")
		}

		c.SetCookie(
			csrfCookieName,
			token,
			0,
			"/",
			"",
			false,
			true,
		)
	} else {
		token = cookieToken
	}

	return token, nil
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

func GeneratePayloadToken(payload any, secretKey string) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256([]byte(secretKey))
	keyString := base64.URLEncoding.EncodeToString(hash[:32])[:43]
	key, err := fernet.DecodeKey(keyString + "=")
	if err != nil {
		return "", err
	}

	token, err := fernet.EncryptAndSign(data, key)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func VerifyPayloadToken[T any](token string, secretKey string) (*T, error) {
	hash := sha256.Sum256([]byte(secretKey))
	keyString := base64.URLEncoding.EncodeToString(hash[:32])[:43]
	key, err := fernet.DecodeKey(keyString + "=")
	if err != nil {
		return nil, err
	}

	msg := fernet.VerifyAndDecrypt([]byte(token), 0, []*fernet.Key{key})

	if msg == nil {
		return nil, fmt.Errorf("invalid or expired token")
	}

	var payload T
	if err := json.Unmarshal(msg, &payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func RenderTextTemplate(filename string, data any) (string, error) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func ValidateStruct(obj any) error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.Struct(obj); err != nil {
			return err
		}
	}
	return errors.New("validator not initialized")
}
