package config

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	AppEnv  string
	AppPort string
	AppName string

	// JWT
	JWTSecret     string
	JWTAccessTTL  time.Duration
	JWTRefreshTTL time.Duration

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Redis
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// LavinMQ
	AMQPUrl string

	// CORS
	CORSAllowedOrigins []string
}

// Load reads environment variables and returns a Config
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, reading from environment")
	}

	accessTTL, _ := time.ParseDuration(getEnv("JWT_ACCESS_TTL", "15m"))
	refreshTTL, _ := time.ParseDuration(getEnv("JWT_REFRESH_TTL", "168h"))

	corsOrigins := strings.Split(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"), ",")

	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		AppName:            getEnv("APP_NAME", "CWCU"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me"),
		JWTAccessTTL:       accessTTL,
		JWTRefreshTTL:      refreshTTL,
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "cwcu_user"),
		DBPassword:         getEnv("DB_PASSWORD", "cwcu_password"),
		DBName:             getEnv("DB_NAME", "cwcu_db"),
		DBSSLMode:          getEnv("DB_SSLMODE", "disable"),
		RedisHost:          getEnv("REDIS_HOST", "localhost"),
		RedisPort:          getEnv("REDIS_PORT", "6379"),
		RedisPassword:      getEnv("REDIS_PASSWORD", ""),
		AMQPUrl:            getEnv("AMQP_URL", "amqp://guest:guest@localhost:5672/"),
		CORSAllowedOrigins: corsOrigins,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
