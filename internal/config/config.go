package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"time"
)

// TODO: hide vulnerable vals in .env

type Config struct {
	HTTP `yaml:"http"`
	JWT  `yaml:"jwt"`
	Log  `yaml:"log"`
	PG   `yaml:"pg"`
}

type HTTP struct {
	Port         string        `env-required:"true" yaml:"port" env:"HTTP_PORT"`
	ReadTimeout  time.Duration `env-required:"true" yaml:"readTimeout" env:"HTTP_READ_TIMEOUT"`
	WriteTimeout time.Duration `env-required:"true" yaml:"writeTimeout" env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout  time.Duration `env-required:"true" yaml:"idleTimeout" env:"HTTP_IDLE_TIMEOUT"`
}

type JWT struct {
	SignKey  string        `env-required:"true" yaml:"signKey" env:"JWT_SIGN_KEY"`
	TokenTTL time.Duration `env-required:"true" yaml:"tokenTTL" env:"JWT_TOKEN_TTL"`
}

type Log struct {
	Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
}

type PG struct {
	URL         string `env-required:"true" yaml:"url" env:"PG_URL"`
	MaxPoolSize int    `env-required:"true" yaml:"maxPoolSize" env:"PG_MAX_POOL_SIZE"`
}

func MustLoad(configPath string) *Config {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("error reading config file: %w", err))
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		log.Fatal(fmt.Errorf("error updating env: %w", err))
	}

	return cfg
}
