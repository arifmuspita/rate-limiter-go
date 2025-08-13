package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort           string
	UseRedis             bool
	RedisURL             string
	RedisPassword        string
	DefaultCycleDuration int
	DefaultMaxRequests   int
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		ServerPort:           getEnv("SERVER_PORT", "8080"),
		UseRedis:             getEnvAsBool("USE_REDIS", false),
		RedisURL:             getEnv("REDIS_URL", "redis://localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		DefaultCycleDuration: getEnvAsInt("DEFAULT_CYCLE_DURATION", 1),
		DefaultMaxRequests:   getEnvAsInt("DEFAULT_MAX_REQUESTS", 100),
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

func getEnvAsInt(key string, value int) int {
	if v := os.Getenv(key); v != "" {
		if intVal, err := strconv.Atoi(v); err == nil {
			return intVal
		}
	}
	return value
}
