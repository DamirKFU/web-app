package core

import (
	"app/config"
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

func NewRedisServer(cfg config.Config) *RedisServer {
	rdb0 := redis.NewClient(&redis.Options{
		Addr:     cfg.REDIS.Addr,
		Password: cfg.REDIS.Password,
		DB:       0,
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
	})

	if _, err := rdb0.Ping(context.Background()).Result(); err != nil {
		log.Panicf("failed to connect to redis: %v", err)
	}
	return &RedisServer{RDB0: rdb0}

}
