package middlewares

import (
	"app/internal/core"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

func inArray(arr []string, value string) bool {
	inarr := slices.Contains(arr, value)

	return inarr
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

	CsrfCookie := server.Cfg.CSRF.Cookie

	return func(c *gin.Context) {
		if inArray(ignoreMethods, c.Request.Method) {
			c.Next()
			return
		}

		if inArray(server.Cfg.CSRF.ExcludePaths, c.Request.URL.Path) {
			c.Next()
			return
		}

		token := tokenGetter(c)
		tokenFromCookie, err := c.Cookie(CsrfCookie)
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
