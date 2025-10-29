package core

import (
	"app/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type RedisServer struct {
	RDB0 *redis.Client
}

type Server struct {
	Cfg         config.Config
	DB          *gorm.DB
	Eng         *gin.Engine
	RoutesMap   map[string]Route
	RedisServer *RedisServer
}

type Route struct {
	Method                string
	Path                  string
	HandlerFuncs          []gin.HandlerFunc
	NameSpace             string
	DecoratorHandlerFuncs []gin.HandlerFunc
}

type MiddlewareGroup struct {
	Prefix     string
	Middleware []gin.HandlerFunc
	Routes     []Route
}
