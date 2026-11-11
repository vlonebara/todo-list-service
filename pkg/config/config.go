package config

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/vlonebara/todo-list/pkg/errutil"
)

// Config содержит настройки, которые читаются один раз при запуске.
type Config struct {
	Host          string
	Port          string
	DBFile        string
	Password      string
	JWTSecret     string
	SecureCookies bool
}

// Load загружает .env и возвращает готовую конфигурацию приложения.
func Load() (Config, error) {
	if err := loadEnvFile(".env"); errutil.CheckError(err) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	cfg := Config{
		Host:      envOrDefault("TODO_HOST", "127.0.0.1"),
		Port:      envOrDefault("TODO_PORT", "7540"),
		DBFile:    envOrDefault("TODO_DBFILE", "scheduler.db"),
		Password:  os.Getenv("TODO_PASSWORD"),
		JWTSecret: os.Getenv("TODO_JWT_SECRET"),
	}
	port, err := strconv.Atoi(cfg.Port)
	if errutil.CheckError(err) || port < 1 || port > 65535 {
		return Config{}, fmt.Errorf("port must be between 1 and 65535")
	}
	local := cfg.Host == "localhost" || net.ParseIP(cfg.Host).IsLoopback()
	if cfg.Password == "" && !local {
		return Config{}, fmt.Errorf("password is required when listening outside localhost")
	}
	cfg.SecureCookies, err = strconv.ParseBool(envOrDefault("TODO_COOKIE_SECURE", strconv.FormatBool(!local)))
	if errutil.CheckError(err) {
		return Config{}, fmt.Errorf("cookie secure setting must be true or false")
	}
	return cfg, nil
}

// envOrDefault возвращает значение переменной или указанное значение по умолчанию.
func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
