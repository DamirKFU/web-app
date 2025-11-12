package core

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dvwright/xss-mw"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
			Fail(c, http.StatusForbidden, "server error", nil)
			log.Println("[WARN] Check if the user exists in the context; if not, either skip the handler (abort) or let it run after logging a warning.")
			c.Abort()
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
			Fail(c, http.StatusForbidden, "CSRF token missing in Headers", nil)
			c.Abort()
			return
		} else if err != nil {
			Fail(c, http.StatusForbidden, "CSRF token missing in Cookie", nil)
			c.Abort()
			return
		}
		if tokenFromCookie != token {
			Fail(c, http.StatusForbidden, "CSRF token mismatch", nil)
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

func XSSMiddleware(server *Server) gin.HandlerFunc {
	xssMdlwr := &xss.XssMw{}
	return xssMdlwr.RemoveXss()
}

func LimitRequestBodySizeMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			server.Cfg.LimitBodySize,
		)
		c.Next()
	}
}

func TransactionMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := server.DB.Begin()
		if tx.Error != nil {
			Fail(c, http.StatusInternalServerError, "Failed to start transaction", nil)
			c.Abort()
			return
		}

		c.Set("db", tx)

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			}
		}()

		c.Next()

		if forceRollback, exists := c.Get("force_transaction_rollback"); exists {
			if rollback, ok := forceRollback.(bool); ok && rollback {
				tx.Rollback()
				return
			}
		}

		status := c.Writer.Status()
		if status >= http.StatusBadRequest {
			tx.Rollback()
		} else {
			if err := tx.Commit().Error; err != nil {
				tx.Rollback()
				Fail(c, http.StatusInternalServerError, "Failed transaction commit", nil)
			}
		}
	}
}

func RateLimiterMiddleware(server *Server, key string, limit int, slidingWindow time.Duration) gin.HandlerFunc {
	redisClient := server.RedisServer.RDB0
	slidingWindow = 60 * time.Second

	return func(c *gin.Context) {
		ctx := c.Request.Context()
		now := time.Now().UnixNano()
		userCntKey := fmt.Sprint(c.ClientIP(), ":", key)

		redisClient.ZRemRangeByScore(
			ctx,
			userCntKey,
			"0",
			fmt.Sprint(now-(slidingWindow.Nanoseconds())),
		).Result()

		reqs, _ := redisClient.ZRange(ctx, userCntKey, 0, -1).Result()

		if len(reqs) >= limit {
			Fail(c, http.StatusTooManyRequests, "too many request", nil)
			c.Abort()
			return
		}

		c.Next()
		redisClient.ZAddNX(
			ctx,
			userCntKey,
			redis.Z{
				Score:  float64(now),
				Member: float64(now),
			},
		)
		redisClient.Expire(
			ctx,
			userCntKey,
			slidingWindow,
		)
	}

}
