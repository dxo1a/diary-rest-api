package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config содержит настройки auth-service.
type Config struct {
	Port       string
	DBPath     string
	JWTSecret  string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("AUTH_PORT", "8080")
	dbPath := getEnv("AUTH_DB_PATH", "./data/auth.db")
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}

	atSec, err := strconv.Atoi(getEnv("ACCESS_TTL_SEC", "900"))
	if err != nil {
		return nil, fmt.Errorf("ACCESS_TTL_SEC invalid: %v", err)
	}
	rtH, err := strconv.Atoi(getEnv("REFRESH_TTL_HOURS", "24"))
	if err != nil {
		return nil, fmt.Errorf("REFRESH_TTL_HOURS invalid: %v", err)
	}

	return &Config{
		Port:       port,
		DBPath:     dbPath,
		JWTSecret:  secret,
		AccessTTL:  time.Duration(atSec) * time.Second,
		RefreshTTL: time.Duration(rtH) * time.Hour,
	}, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
