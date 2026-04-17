package configs

import (
	"os"
	"strconv"
)

// Config holds all application configurations.
type Config struct {
	AppEnv   string
	AppName  string
	HTTPPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	RedisAddr string

	JWTSecret string

	FrontendURL string
	APIBaseURL  string

	GitHubClientID     string
	GitHubClientSecret string
	GoogleClientID     string
	GoogleClientSecret string
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func Load() (*Config, error) {
	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppName:            getEnv("APP_NAME", "user-service"),
		HTTPPort:           getEnv("HTTP_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "blog"),
		DBPassword:         getEnv("DB_PASSWORD", "blog123"),
		DBName:             getEnv("DB_NAME", "blog"),
		RedisAddr:          getEnv("REDIS_ADDR", "localhost:6379"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:3000"),
		APIBaseURL:         getEnv("API_BASE_URL", "http://localhost:8080/api/v1"),
		GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
	}, nil
}

func mustAtoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
