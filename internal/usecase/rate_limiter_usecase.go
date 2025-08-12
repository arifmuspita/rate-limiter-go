package usecase

import "rate-limiter-go/internal/repository/memory"

type RateLimiterUseCase interface {
	CheckRateLimit(clientID string) (allowed bool, remaining int, resetTime int64)
	ConfigureRateLimit(clientID string, maxRequests int, windowSizeMin int) error
}

type rateLimiterUseCase struct {
	repo *memory.RateLimiterRepository
}

func NewRateLimiterUseCase(repo *memory.RateLimiterRepository) RateLimiterUseCase {
	return &rateLimiterUseCase{
		repo: repo,
	}
}

func (uc *rateLimiterUseCase) CheckRateLimit(clientID string) (allowed bool, remaining int, resetTime int64) {

	rateLimit, exists := uc.repo.Get(clientID)
	if !exists {
		rateLimit = uc.repo.Add(clientID)
		uc.repo.Save(rateLimit)
	}

	allowed = rateLimit.IsAllowed()

	if allowed {
		rateLimit.Increment()
		uc.repo.Save(rateLimit)
	}

	remaining = rateLimit.GetRemainingRequests()
	resetTime = rateLimit.GetResetTime().Unix()

	return allowed, remaining, resetTime
}

func (uc *rateLimiterUseCase) ConfigureRateLimit(clientID string, maxRequests int, windowSizeMin int) error {
	rateLimit, exists := uc.repo.Get(clientID)
	if !exists {
		rateLimit = uc.repo.Add(clientID)
	}

	rateLimit.MaxRequests = maxRequests
	rateLimit.CycleDuration = windowSizeMin

	return uc.repo.Save(rateLimit)
}
