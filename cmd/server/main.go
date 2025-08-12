package main

import (
	"log"
	"rate-limiter-go/internal/handler"
	"rate-limiter-go/internal/repository/memory"
	"rate-limiter-go/internal/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	repo := memory.NewRateLimiterRepository()
	useCase := usecase.NewRateLimiterUseCase(repo)
	rateLimiterHandler := handler.NewRateLimiterHandler(useCase)

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
