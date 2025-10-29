package core

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"

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

func ParseValidationError(err error) string {
	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range errs {
			return e.Field() + ": invalid"
		}
	}
	return err.Error()
}

func HandleError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}
	return false
}
