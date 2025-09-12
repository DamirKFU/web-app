package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type HTTP struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type Config struct {
	HTTP       HTTP `mapstructure:"http"`
	SECRET_KEY string
}

func Load() Config {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	var cfg Config
	_ = v.Unmarshal(&cfg)

	// значения по умолчанию
	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8080"
	}
	if cfg.HTTP.ReadTimeout == 0 {
		cfg.HTTP.ReadTimeout = 5 * time.Second
	}
	if cfg.HTTP.WriteTimeout == 0 {
		cfg.HTTP.WriteTimeout = 10 * time.Second
	}

	if cfg.SECRET_KEY == "" {
		log.Fatalf("SECRET KEY IS EMTY")
	}

	return cfg
}
