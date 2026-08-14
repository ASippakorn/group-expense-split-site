package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	APIAddr           string
	DatabaseURL       string
	WebOrigin         string
	SessionCookieName string
	SessionTTL        time.Duration
	PasswordPepper    string
}

func Load() Config {
	ttlHours := intEnv("SESSION_TTL_HOURS", 168)

	return Config{
		APIAddr:           stringEnv("API_ADDR", ":8080"),
		DatabaseURL:       stringEnv("DATABASE_URL", "postgres://splitr:splitr@localhost:5432/splitr?sslmode=disable"),
		WebOrigin:         stringEnv("WEB_ORIGIN", "http://localhost:5173"),
		SessionCookieName: stringEnv("SESSION_COOKIE_NAME", "splitr_session"),
		SessionTTL:        time.Duration(ttlHours) * time.Hour,
		PasswordPepper:    stringEnv("PASSWORD_PEPPER", "dev-pepper"),
	}
}

func stringEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func intEnv(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
