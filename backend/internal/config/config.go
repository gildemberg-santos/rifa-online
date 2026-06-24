package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    int
	MongoURI                string
	MongoDBName             string
	RedisURI                string
	JWTSecret               string
	AbacatepayAPIKey        string
	AbacatepayWebhookSecret string
	FrontendURL             string
	LogLevel                slog.Level
	LogFormat               string
}

func Load() *Config {
	godotenv.Load()

	level := slog.LevelInfo
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	}

	format := getEnv("LOG_FORMAT", "text")

	return &Config{
		Port:                    getEnvInt("PORT", 8080),
		MongoURI:                getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		MongoDBName:             getEnv("MONGODB_DB_NAME", "rifaonline"),
		RedisURI:                getEnv("REDIS_URI", "redis://localhost:6379"),
		JWTSecret:               getEnv("JWT_SECRET", "change-me"),
		AbacatepayAPIKey:        getEnv("ABACATEPAY_API_KEY", ""),
		AbacatepayWebhookSecret: getEnv("ABACATEPAY_WEBHOOK_SECRET", ""),
		FrontendURL:             getEnv("FRONTEND_URL", "http://localhost:5173"),
		LogLevel:                level,
		LogFormat:               format,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
