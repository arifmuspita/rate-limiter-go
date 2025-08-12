package handler

import (
	"net/http"
	"rate-limiter-go/internal/usecase"
	"rate-limiter-go/pkg/response"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RateLimiterMiddleware(useCase usecase.RateLimiterUseCase) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			clientID := c.Request().Header.Get("X-Client-ID")
			if clientID == "" {
				clientID = c.RealIP()
			}

			allowed, remaining, resetTime := useCase.CheckRateLimit(clientID)

			c.Response().Header().Set("X-RateLimit-Limit", "100")
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

			if !allowed {
				return response.Error(c, http.StatusTooManyRequests, "Rate limit exceeded")
			}

			return next(c)
		}
	}
}
