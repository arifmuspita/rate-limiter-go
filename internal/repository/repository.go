package repository

import (
	"context"
	"rate-limiter-go/internal/domain"
)

type RateLimiterRepository interface {
	Get(ctx context.Context, clientID string) (*domain.RateLimit, bool, error)
	Save(ctx context.Context, rateLimit *domain.RateLimit) error
	Delete(ctx context.Context, clientID string) error
	CreateDefault(ctx context.Context, clientID string) *domain.RateLimit
}
