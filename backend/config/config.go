package config

import (
	"encoding/json"
	"log"
	"time"

	"github.com/spf13/viper"
)

type HTTP struct {
	Addr         string        `mapstructure:"addr"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
}

type JWT struct {
	AccessCookie     string        `mapstructure:"access_cookie"`
	RefreshCookie    string        `mapstructure:"refresh_cookie"`
	AccessExpiresIn  time.Duration `mapstructure:"access_expires_in"`
	RefreshExpiresIn time.Duration `mapstructure:"refresh_expires_in"`
}

type CORS struct {
	AllowOrigins     []string `mapstructure:"allow_origins"`
	AllowMethods     []string `mapstructure:"allow_methods"`
	AllowHeaders     []string `mapstructure:"allow_headers"`
	AllowCredentials bool     `mapstructure:"allow_credentials"`
}

type REDIS struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
}

type CSRF struct {
	Cookie       string   `mapstructure:"csrf_cookie"`
	ExcludePaths []string `mapstructure:"exclude_paths"`
}

type Config struct {
	HTTP      HTTP     `mapstructure:"http"`
	DB        Database `mapstructure:"database"`
	JWT       JWT      `mapstructure:"jwt"`
	CORS      CORS     `mapstructure:"cors"`
	CSRF      CSRF     `mapstructure:"csrf"`
	REDIS     REDIS    `mapstructure:"redis"`
	SecretKey string   `mapstructure:"secret_key"`
}

func Load() Config {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	v.SetDefault("http.addr", ":8080")
	v.SetDefault("http.read_timeout", 5*time.Second)
	v.SetDefault("http.write_timeout", 10*time.Second)

	v.SetDefault("jwt.access_cookie", "access_token")
	v.SetDefault("jwt.refresh_cookie", "refresh_token")
	v.SetDefault("jwt.access_expires_in", 15*time.Minute)
	v.SetDefault("jwt.refresh_expires_in", 7*24*time.Hour)

	v.SetDefault("redis.addr", "localhost:6379")

	v.SetDefault("cookies.csrf_cookie", "csrf_token")

	if err := v.ReadInConfig(); err != nil {
		log.Printf("warning: could not read config file: %v", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	if cfg.SecretKey == "" {
		log.Fatalf("JWT SECRET KEY IS EMPTY")
	}

	data, _ := json.MarshalIndent(cfg, "", "  ")
	log.Println(string(data))
	return cfg
}
