package core

import (
	"net/http"
	"slices"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(server *Server) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     server.Cfg.CORS.AllowOrigins,
		AllowMethods:     server.Cfg.CORS.AllowMethods,
		AllowHeaders:     server.Cfg.CORS.AllowHeaders,
		AllowCredentials: server.Cfg.CORS.AllowCredentials,
	})
}

func inArray(arr []string, value string) bool {
	inarr := slices.Contains(arr, value)

	return inarr
}

func CSRFMiddleware(server *Server) gin.HandlerFunc {
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
		value, ok := c.Get("csrf_exempt")

		csrf_exempt := ok && value.(bool)

		if csrf_exempt {
			c.Next()
			return
		}

		if inArray(ignoreMethods, c.Request.Method) && !ok {
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

func CsrfExemptMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("csrf_exempt", true)
		c.Next()
	}
}

func CsrfEnforceMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("csrf_exempt", false)
		c.Next()
	}
}
