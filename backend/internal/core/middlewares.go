package core

import (
	"log"
	"net/http"

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
		user, ok := c.Get("user")
		if !ok {
			c.AbortWithStatus(http.StatusInternalServerError)
			log.Println("[WARN] Check if the user exists in the context; if not, either skip the handler (abort) or let it run after logging a warning.")
			return
		}

		if user == nil {
			c.Next()
			return
		}

		csrfExempt, correct := CheckCsrfExempt(c)

		if csrfExempt || (!correct && inArray(ignoreMethods, c.Request.Method)) {
			c.Next()
			return
		}

		token := tokenGetter(c)
		tokenFromCookie, err := c.Cookie(CsrfCookie)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, "CSRF token missing in Headers")
			return
		} else if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, "CSRF token missing in Cookie")
			return
		}
		if tokenFromCookie != token {
			c.AbortWithStatusJSON(http.StatusForbidden, "CSRF token mismatch")
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
