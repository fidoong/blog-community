package configs

// Config holds all application configurations.
type Config struct {
	AppEnv   string `mapstructure:"APP_ENV"`
	AppName  string `mapstructure:"APP_NAME"`
	HTTPPort string `mapstructure:"HTTP_PORT"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	RedisAddr string `mapstructure:"REDIS_ADDR"`

	JWTSecret string `mapstructure:"JWT_SECRET"`

	FrontendURL string `mapstructure:"FRONTEND_URL"`
	APIBaseURL  string `mapstructure:"API_BASE_URL"`

	GitHubClientID     string `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
}

func Load() (*Config, error) {
	// TODO: integrate viper for real env loading
	return &Config{
		AppEnv:      "development",
		AppName:     "user-service",
		HTTPPort:    "8080",
		DBHost:      "localhost",
		DBPort:      "5432",
		DBUser:      "blog",
		DBPassword:  "blog123",
		DBName:      "blog",
		RedisAddr:   "localhost:6379",
		JWTSecret:   "super-secret-key-change-in-production",
		FrontendURL: "http://localhost:3000",
		APIBaseURL:  "http://localhost:8080/api/v1",
	}, nil
}
