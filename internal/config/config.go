package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP
	JWT
	Log
	PG
}

type HTTP struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func (h HTTP) Address() string {
	return net.JoinHostPort(h.Host, h.Port)
}

type JWT struct {
	SignKey  string
	TokenTTL time.Duration
}

type Log struct {
	Level string
}

type PG struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	MaxPoolSize int
}

func (pg PG) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		pg.Host, pg.Port, pg.Database, pg.User, pg.Password,
	)
}

func MustLoad(envPath string) *Config {
	err := godotenv.Load(envPath)
	if err != nil && !os.IsNotExist(err) {
		panic(fmt.Errorf("failed to load .env file: %w", err))
	}

	cfg := &Config{}
	cfg.HTTP, err = loadHTTPConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load http config: %w", err))
	}
	cfg.JWT, err = loadJWTConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load jwt config: %w", err))
	}
	cfg.Log, err = loadLogConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load log config: %w", err))
	}
	cfg.PG, err = loadPGConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load pg config: %w", err))
	}

	return cfg
}

func loadHTTPConfig() (HTTP, error) {
	host, err := getEnv("HTTP_HOST")
	if err != nil {
		return HTTP{}, fmt.Errorf("missing HTTP_HOST: %w", err)
	}
	port, err := getEnv("HTTP_PORT")
	if err != nil {
		return HTTP{}, fmt.Errorf("missing HTTP_PORT: %w", err)
	}
	readTimeout, err := parseDuration("HTTP_READ_TIMEOUT")
	if err != nil {
		return HTTP{}, fmt.Errorf("invalid or missing HTTP_READ_TIMEOUT: %w", err)
	}
	writeTimeout, err := parseDuration("HTTP_WRITE_TIMEOUT")
	if err != nil {
		return HTTP{}, fmt.Errorf("invalid or missing HTTP_WRITE_TIMEOUT: %w", err)
	}
	idleTimeout, err := parseDuration("HTTP_IDLE_TIMEOUT")
	if err != nil {
		return HTTP{}, fmt.Errorf("invalid or missing HTTP_IDLE_TIMEOUT: %w", err)
	}

	return HTTP{
		Host:         host,
		Port:         port,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}, nil
}

func loadJWTConfig() (JWT, error) {
	signKey, err := getEnv("JWT_SIGN_KEY")
	if err != nil {
		return JWT{}, fmt.Errorf("missing JWT_SIGN_KEY: %w", err)
	}
	tokenTTL, err := parseDuration("JWT_TOKEN_TTL")
	if err != nil {
		return JWT{}, fmt.Errorf("invalid or missing JWT_TOKEN_TTL: %w", err)
	}

	return JWT{
		SignKey:  signKey,
		TokenTTL: tokenTTL,
	}, nil
}

func loadLogConfig() (Log, error) {
	level, err := getEnv("LOGGER_LEVEL")
	if err != nil {
		return Log{}, fmt.Errorf("missing LOGGER_LEVEL: %w", err)
	}

	return Log{
		Level: level,
	}, nil
}

func loadPGConfig() (PG, error) {
	host, err := getEnv("POSTGRES_HOST")
	if err != nil {
		return PG{}, fmt.Errorf("missing POSTGRES_HOST: %w", err)
	}
	port, err := getEnv("POSTGRES_PORT")
	if err != nil {
		return PG{}, fmt.Errorf("missing POSTGRES_PORT: %w", err)
	}
	user, err := getEnv("POSTGRES_USER")
	if err != nil {
		return PG{}, fmt.Errorf("missing POSTGRES_USER: %w", err)
	}
	password, err := getEnv("POSTGRES_PASSWORD")
	if err != nil {
		return PG{}, fmt.Errorf("missing POSTGRES_PASSWORD: %w", err)
	}
	database, err := getEnv("POSTGRES_DB")
	if err != nil {
		return PG{}, fmt.Errorf("missing POSTGRES_DB: %w", err)
	}
	maxPoolSizeStr, err := getEnv("POSTGRES_MAX_POOL_SIZE")
	if err != nil {
		return PG{}, fmt.Errorf("missing POSTGRES_MAX_POOL_SIZE: %w", err)
	}

	maxPoolSize, err := strconv.Atoi(maxPoolSizeStr)
	if err != nil {
		return PG{}, fmt.Errorf("invalid POSTGRES_MAX_POOL_SIZE: %w", err)
	}

	return PG{
		Host:        host,
		Port:        port,
		User:        user,
		Password:    password,
		Database:    database,
		MaxPoolSize: maxPoolSize,
	}, nil
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return value, nil
}

func parseDuration(key string) (time.Duration, error) {
	value, err := getEnv(key)
	if err != nil {
		return 0, err
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return 0, fmt.Errorf("invalid duration format for %s: %w", key, err)
	}
	return duration, nil
}
