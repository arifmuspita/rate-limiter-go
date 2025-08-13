package usecase

import (
	"context"
	"rate-limiter-go/internal/repository"
	"sync"
)

type RateLimiterUseCase interface {
	CheckRateLimit(ctx context.Context, clientID string) (allowed bool, remaining int, resetTime int64)
	ConfigureRateLimit(ctx context.Context, clientID string, maxRequests int, windowSizeMin int) error
}

type rateLimiterUseCase struct {
	repo repository.RateLimiterRepository
	mu   sync.Mutex
}

func NewRateLimiterUseCase(repo repository.RateLimiterRepository) RateLimiterUseCase {
	return &rateLimiterUseCase{
		repo: repo,
	}
}

func (uc *rateLimiterUseCase) CheckRateLimit(ctx context.Context, clientID string) (allowed bool, remaining int, resetTime int64) {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	rateLimit, exists, err := uc.repo.Get(ctx, clientID)
	if err != nil {
		rateLimit = uc.repo.CreateDefault(ctx, clientID)
	}
	if !exists {
		rateLimit = uc.repo.CreateDefault(ctx, clientID)
		uc.repo.Save(ctx, rateLimit)
	}

	allowed = rateLimit.IsAllowed()

	if allowed {
		rateLimit.Increment()
		uc.repo.Save(ctx, rateLimit)
	}

	remaining = rateLimit.GetRemainingRequests()
	resetTime = rateLimit.GetResetTime().Unix()

	return allowed, remaining, resetTime
}

func (uc *rateLimiterUseCase) ConfigureRateLimit(ctx context.Context, clientID string, maxRequests int, cycleDurationMin int) error {
	uc.mu.Lock()
	defer uc.mu.Unlock()

	rateLimit, exists, err := uc.repo.Get(ctx, clientID)
	if err != nil {
		return err
	}
	if !exists {
		rateLimit = uc.repo.CreateDefault(ctx, clientID)
	}

	rateLimit.MaxRequests = maxRequests
	rateLimit.CycleDuration = cycleDurationMin

	return uc.repo.Save(ctx, rateLimit)
}
