package memory

import (
	"rate-limiter-go/internal/domain"
	"sync"
	"time"
)

type RateLimiterRepository struct {
	mu    sync.RWMutex
	store map[string]*domain.RateLimit
}

func NewRateLimiterRepository() *RateLimiterRepository {
	return &RateLimiterRepository{
		store: make(map[string]*domain.RateLimit),
	}
}

func (r *RateLimiterRepository) Get(clientID string) (*domain.RateLimit, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rateLimit, exists := r.store[clientID]
	return rateLimit, exists
}

func (r *RateLimiterRepository) Save(rateLimit *domain.RateLimit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.store[rateLimit.ClientID] = rateLimit
	return nil
}

func (r *RateLimiterRepository) Add(clientID string) *domain.RateLimit {
	return &domain.RateLimit{
		ClientID:      clientID,
		RequestCount:  0,
		CycleStart:    time.Now(),
		CycleDuration: 1,
		MaxRequests:   20,
	}
}
