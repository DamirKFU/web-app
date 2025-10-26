package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"io"

	"github.com/dchest/uniuri"
)

func GenerateCSRFToken(secret string) string {
	randomPart := uniuri.NewLen(16)

	h := sha1.New()
	io.WriteString(h, randomPart+"-"+secret)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash
}
