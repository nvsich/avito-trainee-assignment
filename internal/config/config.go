package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

// TODO: hide vulnerable vals in .env

type Config struct {
	HTTP `yaml:"http"`
	JWT  `yaml:"jwt"`
	Log  `yaml:"log"`
	PG   `yaml:"postgres"`
}

type HTTP struct {
	Port string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
}

type JWT struct {
	SignKey  string        `env-required:"true" yaml:"signKey" env:"JWT_SIGN_KEY"`
	TokenTTL time.Duration `env-required:"true" yaml:"token_ttl" env:"JWT_TOKEN_TTL"`
}

type Log struct {
	Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
}

type PG struct {
	URL string `env-required:"true" yaml:"url" env:"PG_URL"`
}

func Load(configPath string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	err = cleanenv.UpdateEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("error updating env: %w", err)
	}

	return cfg, nil
}
