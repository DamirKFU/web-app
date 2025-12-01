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

type SMTP struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
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
	Cookie string `mapstructure:"csrf_cookie"`
}

type Config struct {
	HTTP                  HTTP              `mapstructure:"http"`
	DB                    Database          `mapstructure:"database"`
	JWT                   JWT               `mapstructure:"jwt"`
	CORS                  CORS              `mapstructure:"cors"`
	CSRF                  CSRF              `mapstructure:"csrf"`
	REDIS                 REDIS             `mapstructure:"redis"`
	SMTP                  SMTP              `mapstructure:"smtp"`
	SecretKey             string            `mapstructure:"secret_key"`
	PayloadTokenExpiresIn time.Duration     `mapstructure:"payload_token_expires_in"`
	FrontURL              string            `mapstructure:"front_url"`
	LimitBodySize         int64             `mapstructure:"limit_body_size"`
	TestDBName            string            `mapstructure:"test_db_name"`
	Debug                 bool              `mapstructure:"debug"`
	ValidationMessages    map[string]string `mapstructure:"validation_messages"`
}

func Load(config_path string) Config {
	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigName("config")
	v.AddConfigPath(config_path)
	v.AutomaticEnv()

	v.SetDefault("http.addr", ":8080")
	v.SetDefault("http.read_timeout", 5*time.Second)
	v.SetDefault("http.write_timeout", 10*time.Second)

	v.SetDefault("jwt.access_cookie", "access_token")
	v.SetDefault("jwt.refresh_cookie", "refresh_token")
	v.SetDefault("jwt.access_expires_in", 15*time.Minute)
	v.SetDefault("jwt.refresh_expires_in", 7*24*time.Hour)

	v.SetDefault("payload_token_expires_in", time.Hour)
	v.SetDefault("limit_body_size", 1024)
	v.SetDefault("test_db_name", "test_db")
	v.SetDefault("debug", true)

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

	if cfg.DB.Name == cfg.TestDBName {
		log.Fatalf("DB.Name must not be equal to TestDBName")
	}

	return cfg
}
