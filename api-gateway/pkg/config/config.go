package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	CalendarService string
	ActionService   string
	HabitService    string
	MetricsService  string
	AuthService     string
	JWTSecret       string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("GATEWAY_PORT", "80"),
		CalendarService: getEnv("CALENDAR_SERVICE_URL", ""),
		ActionService:   getEnv("ACTION_SERVICE_URL", ""),
		HabitService:    getEnv("HABIT_SERVICE_URL", ""),
		MetricsService:  getEnv("METRICS_SERVICE_URL", ""),
		AuthService:     getEnv("AUTH_SERVICE_URL", ""),
		JWTSecret:       getEnv("JWT_SECRET", ""),
	}

	if cfg.CalendarService == "" ||
		cfg.ActionService == "" ||
		cfg.HabitService == "" ||
		cfg.AuthService == "" ||
		cfg.MetricsService == "" {
		return nil, fmt.Errorf("one or more service URLs are not set")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
