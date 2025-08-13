package integration

import (
	"context"
	"rate-limiter-go/internal/handler"
	"rate-limiter-go/internal/repository"
	"rate-limiter-go/internal/repository/memory"
	redisRepo "rate-limiter-go/internal/repository/redis"
	"rate-limiter-go/internal/usecase"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type TestApp struct {
	Echo       *echo.Echo
	Repository repository.RateLimiterRepository
	UseCase    usecase.RateLimiterUseCase
	Handler    *handler.RateLimiterHandler
}

func SetupTestApp(t *testing.T, useRedis bool) *TestApp {
	var repo repository.RateLimiterRepository

	if useRedis {

		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			DB:   15,
		})

		ctx := context.Background()
		if err := client.Ping(ctx).Err(); err != nil {
			t.Skip("Redis not available, skipping integration test")
		}

		client.FlushDB(ctx)

		repo = redisRepo.NewRateLimiterRedisRepository(client, 100, 1)
	} else {
		repo = memory.NewRateLimiterMemoryRepository(100, 1)
	}

	uc := usecase.NewRateLimiterUseCase(repo)
	h := handler.NewRateLimiterHandler(uc)

	e := echo.New()

	api := e.Group("/api/v1")
	api.GET("/rate-limit/:clientID", h.CheckRateLimit)
	api.PUT("/rate-limit/:clientID", h.ConfigureRateLimit)

	protected := api.Group("/protected")
	protected.Use(handler.RateLimiterMiddleware(uc))
	protected.GET("/data", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "OK"})
	})

	return &TestApp{
		Echo:       e,
		Repository: repo,
		UseCase:    uc,
		Handler:    h,
	}
}
