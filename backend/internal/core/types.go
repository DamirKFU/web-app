package core

import (
	"app/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type RedisServer struct {
	RDB0 *redis.Client
}

type Server struct {
	Cfg         config.Config
	DB          *gorm.DB
	Eng         *gin.Engine
	RoutesMap   map[string]Route
	RedisServer *RedisServer
	Email       *EmailSMTPEngine
	Mdls        []gin.HandlerFunc
	Models      []any
}

type Route struct {
	Method                string
	Path                  string
	HandlerFuncs          []gin.HandlerFunc
	NameSpace             string
	DecoratorHandlerFuncs []gin.HandlerFunc
	FullPath              string
}

type MiddlewareGroup struct {
	Prefix     string
	Middleware []gin.HandlerFunc
	Routes     []Route
}

type APIResponse struct {
	Success bool      `json:"success"`
	Data    any       `json:"data,omitempty"`
	Error   *APIError `json:"error,omitempty"`
}

type APIError struct {
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

type EmailSMTPEngine struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}
