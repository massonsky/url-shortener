package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort    string
	DatabaseURL   string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	AppEnv        string
}

// Load загружает конфиг из .env или системных переменных
func Load() *Config {
	// Пытаемся загрузить .env (игнорируем ошибку в продакшене)
	_ = godotenv.Load()

	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		AppEnv:        getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
