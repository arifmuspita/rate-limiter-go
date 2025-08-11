package main

import (
	"log"
	"rate-limiter-go/pkg/config"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.UseRedis {
		log.Println("using redis")
	} else {
		log.Println("Using in-memory storage for rate limiting")
	}
}
