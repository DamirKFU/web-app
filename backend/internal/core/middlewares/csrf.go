package middlewares

import (
	"app/internal/core"
	"crypto/sha1"
	"encoding/base64"
	"io"
	"net/http"
	"slices"

	"github.com/dchest/uniuri"
	"github.com/gin-gonic/gin"
)

const (
	csrfCockie = "csrf_token"
)

func inArray(arr []string, value string) bool {
	inarr := slices.Contains(arr, value)

	return inarr
}

func generateCSRFToken(secret string) string {
	randomPart := uniuri.NewLen(16)

	h := sha1.New()
	io.WriteString(h, randomPart+"-"+secret)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))

	return hash
}

func GetToken(c *gin.Context, secret string) string {
	cookieToken, err := c.Cookie(csrfCockie)
	var token string

	if err != nil || cookieToken == "" {
		token = generateCSRFToken(secret)

		c.SetCookie(csrfCockie,
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

	return token
}

func CSRFMiddleware(server *core.Server) gin.HandlerFunc {
	ignoreMethods := []string{"GET", "HEAD", "OPTIONS"}

	tokenGetter := func(c *gin.Context) string {
		r := c.Request

		if t := r.FormValue("_csrf"); len(t) > 0 {
			return t
		} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
			return t
		} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
			return t
		} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
			return t
		}
		return ""
	}

	return func(c *gin.Context) {
		if inArray(ignoreMethods, c.Request.Method) {
			c.Next()
			return
		}

		token := tokenGetter(c)
		tokenFromCookie, err := c.Cookie(csrfCockie)
		if token == "" {
			c.JSON(http.StatusForbidden, "CSRF token missing in Headers")
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusForbidden, "CSRF token missing in Cookie")
			c.Abort()
			return
		}
		if tokenFromCookie != token {
			c.JSON(http.StatusForbidden, "CSRF token mismatch")
			c.Abort()
			return
		}

		c.Next()
	}
}
