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
		ServerPort:         getEnv("SERVER_PORT", "1234"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		UseRedis:           getEnvAsBool("USE_REDIS", false),
		DefaultWindowSize:  getEnvAsDuration("DEFAULT_CYCLE_DURATION", 60*time.Second),
		DefaultMaxRequests: getEnvAsInt("DEFAULT_MAX_REQUESTS", 100),
	}

	return cfg
}

func getEnv(key, value string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return value
}

func getEnvAsBool(key string, value bool) bool {
	if v := os.Getenv(key); v != "" {
		if boolVal, err := strconv.ParseBool(v); err == nil {
			return boolVal
		}
	}
	return value
}

func getEnvAsDuration(key string, value time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if duration, err := time.ParseDuration(v); err == nil {
			return duration
		}
	}
	return value
}

func getEnvAsInt(key string, value int) int {
	if v := os.Getenv(key); v != "" {
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal
		}
	}
	return value
}
