package configs

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configurations.
type Config struct {
	AppEnv   string `mapstructure:"app_env"`
	AppName  string `mapstructure:"app_name"`
	HTTPPort string `mapstructure:"http_port"`

	DBHost     string `mapstructure:"db_host"`
	DBPort     int    `mapstructure:"db_port"`
	DBUser     string `mapstructure:"db_user"`
	DBPassword string `mapstructure:"db_password"`
	DBName     string `mapstructure:"db_name"`

	RedisAddr     string `mapstructure:"redis_addr"`
	RedisPassword string `mapstructure:"redis_password"`

	JWTSecret string `mapstructure:"jwt_secret"`

	FrontendURL string `mapstructure:"frontend_url"`
	APIBaseURL  string `mapstructure:"api_base_url"`

	ElasticsearchAddr string `mapstructure:"elasticsearch_addr"`

	GitHubClientID     string `mapstructure:"github_client_id"`
	GitHubClientSecret string `mapstructure:"github_client_secret"`
	GoogleClientID     string `mapstructure:"google_client_id"`
	GoogleClientSecret string `mapstructure:"google_client_secret"`
}

// Load reads configuration with the following precedence:
// 1. Environment variables (highest)
// 2. .env file in backend directory
// 3. Sensible defaults (lowest)
func Load() (*Config, error) {
	v := viper.New()

	// Defaults
	v.SetDefault("app_env", "development")
	v.SetDefault("app_name", "api-service")
	v.SetDefault("http_port", "8080")
	v.SetDefault("db_host", "127.0.0.1")
	v.SetDefault("db_port", 5432)
	v.SetDefault("db_user", "blog")
	v.SetDefault("db_password", "blog123")
	v.SetDefault("db_name", "blog")
	v.SetDefault("redis_addr", "127.0.0.1:6379")
	v.SetDefault("redis_password", "")
	v.SetDefault("jwt_secret", "super-secret-key-change-in-production")
	v.SetDefault("frontend_url", "http://localhost:3000")
	v.SetDefault("api_base_url", "http://localhost:8080/api/v1")
	v.SetDefault("elasticsearch_addr", "http://localhost:9200")

	// Enable .env file support
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")
	v.AddConfigPath("./backend")

	// Read .env file (ignore error if not found)
	_ = v.ReadInConfig()

	// Bind environment variables with APP_ prefix
	v.SetEnvPrefix("")
	v.AutomaticEnv()
	// Replace dots with underscores for nested access, but here we use flat keys
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Explicitly bind critical env vars that may differ in naming convention
	_ = v.BindEnv("github_client_id", "GITHUB_CLIENT_ID")
	_ = v.BindEnv("github_client_secret", "GITHUB_CLIENT_SECRET")
	_ = v.BindEnv("google_client_id", "GOOGLE_CLIENT_ID")
	_ = v.BindEnv("google_client_secret", "GOOGLE_CLIENT_SECRET")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	// Normalize
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080"
	}

	return &cfg, nil
}
