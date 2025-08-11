package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort         string
	RedisURL           string
	UseRedis           bool
	DefaultWindowSize  time.Duration
	DefaultMaxRequests int
}

func LoadConfig() *Config {
	cfg := &Config{
		ServerPort:         getEnv("SERVER_PORT", "8080"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		UseRedis:           getEnvAsBool("USE_REDIS", false),
		DefaultWindowSize:  getEnvAsDuration("DEFAULT_WINDOW_SIZE", 60*time.Second),
		DefaultMaxRequests: getEnvAsInt("DEFAULT_MAX_REQUESTS", 100),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}
