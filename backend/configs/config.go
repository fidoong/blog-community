package configs

import "os"

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

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	return &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		AppName:     getEnv("APP_NAME", "user-service"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		DBHost:      getEnv("DB_HOST", "192.168.139.191"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "blog"),
		DBPassword:  getEnv("DB_PASSWORD", "blog123"),
		DBName:      getEnv("DB_NAME", "blog"),
		RedisAddr:   getEnv("REDIS_ADDR", "192.168.139.191:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "super-secret-key-change-in-production"),
		FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		APIBaseURL:  getEnv("API_BASE_URL", "http://localhost:8080/api/v1"),
		GitHubClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		GitHubClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}, nil
}
