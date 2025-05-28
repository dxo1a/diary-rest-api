package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	DBPath   string
	CacheTTL time.Duration // TTL для кэша
	LogLevel string        // уровень логирования
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	port := getEnv("ACTION_PORT", "8081")
	dbPath := getEnv("ACTION_DB_PATH", ".data/action.db")
	ttlStr := getEnv("CACHE_TTL_SEC", "60")
	logLevel := getEnv("LOG_LEVEL", "info")

	// конвертация TTL в time.Duration
	ttlSec, err := strconv.Atoi(ttlStr)
	if err != nil {
		return nil, fmt.Errorf("CACHE_TTL_SEC must be integer, got %q", ttlStr)
	}

	cfg := &Config{
		Port:     port,
		DBPath:   dbPath,
		CacheTTL: time.Duration(ttlSec) * time.Second,
		LogLevel: logLevel,
	}

	if cfg.DBPath == "" {
		return nil, fmt.Errorf("DB_PATH must be set")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
