package handler

import (
	"net/http"
	"rate-limiter-go/internal/usecase"
	"rate-limiter-go/pkg/response"

	"github.com/labstack/echo/v4"
)

type RateLimiterHandler struct {
	useCase usecase.RateLimiterUseCase
}

func NewRateLimiterHandler(useCase usecase.RateLimiterUseCase) *RateLimiterHandler {
	return &RateLimiterHandler{
		useCase: useCase,
	}
}

func (h *RateLimiterHandler) CheckRateLimit(c echo.Context) error {

	ctx := c.Request().Context()

	clientID := c.Param("clientID")
	if clientID == "" {
		return response.Error(c, http.StatusBadRequest, "Client ID required")
	}

	allowed, remaining, resetTime := h.useCase.CheckRateLimit(ctx, clientID)

	return response.Success(c, map[string]interface{}{
		"allowed":   allowed,
		"remaining": remaining,
		"reset":     resetTime,
	})
}

func (h *RateLimiterHandler) ConfigureRateLimit(c echo.Context) error {

	ctx := c.Request().Context()

	clientID := c.Param("clientID")
	if clientID == "" {
		return response.Error(c, http.StatusBadRequest, "Client ID required")
	}

	var req struct {
		MaxRequests   int `json:"max_requests"`
		CycleDuration int `json:"cycle_duration"`
	}

	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "Invalid request body")
	}

	if req.MaxRequests <= 0 || req.CycleDuration <= 0 {
		return response.Error(c, http.StatusBadRequest, "Invalid configuration values")
	}

	err := h.useCase.ConfigureRateLimit(ctx, clientID, req.MaxRequests, req.CycleDuration)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "Failed to configure rate limit")
	}

	return response.Success(c, map[string]string{
		"message": "Save setting successfully",
	})
}
