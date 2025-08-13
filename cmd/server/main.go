package main

import (
	"context"
	"log"
	"rate-limiter-go/internal/handler"
	"rate-limiter-go/internal/repository"
	"rate-limiter-go/internal/repository/memory"
	redisRepo "rate-limiter-go/internal/repository/redis"
	"rate-limiter-go/internal/usecase"
	"rate-limiter-go/pkg/config"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {

	cfg := config.LoadConfig()

	var repo repository.RateLimiterRepository

	if cfg.UseRedis {

		repo = initRedisRepository(cfg)
		log.Println("Using Redis for rate limiting")

	} else {

		repo = initMemoryRepository(cfg)
		log.Println("Using memory for rate limiting")
	}

	useCase := usecase.NewRateLimiterUseCase(repo)
	rateLimiterHandler := handler.NewRateLimiterHandler(useCase)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK"})
	})

	api := e.Group("/api/v1")
	api.GET("/rate-limit/:clientID", rateLimiterHandler.CheckRateLimit)
	api.PUT("/rate-limit/:clientID", rateLimiterHandler.ConfigureRateLimit)

	protected := api.Group("/protected")
	protected.Use(handler.RateLimiterMiddleware(useCase))
	protected.GET("/data", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"message": "Protected data",
		})
	})

	log.Println("Server starting on :1234")
	if err := e.Start(":1234"); err != nil {
		log.Fatal(err)
	}
}

func initMemoryRepository(cfg *config.Config) repository.RateLimiterRepository {
	return memory.NewRateLimiterMemoryRepository(
		cfg.DefaultMaxRequests,
		int(cfg.DefaultCycleDuration),
	)
}

func initRedisRepository(cfg *config.Config) repository.RateLimiterRepository {

	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	if cfg.RedisPassword != "" {
		opt.Password = cfg.RedisPassword
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Successfully connected to Redis")

	return redisRepo.NewRateLimiterRedisRepository(
		client,
		cfg.DefaultMaxRequests,
		int(cfg.DefaultCycleDuration),
	)
}
