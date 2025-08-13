package memory

import (
	"context"
	"rate-limiter-go/internal/domain"
	"rate-limiter-go/internal/repository"
	"sync"
	"time"
)

type memoryRateLimiterRepository struct {
	mu                   sync.RWMutex
	store                map[string]*domain.RateLimit
	defaultMaxRequest    int
	defaultCycleDuration int
}

func NewRateLimiterMemoryRepository(defaultMaxRequests, defaultCycleDuration int) repository.RateLimiterRepository {
	return &memoryRateLimiterRepository{
		store:                make(map[string]*domain.RateLimit),
		defaultMaxRequest:    defaultMaxRequests,
		defaultCycleDuration: defaultCycleDuration,
	}
}

func (r *memoryRateLimiterRepository) Get(ctx context.Context, clientID string) (*domain.RateLimit, bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rateLimit, exists := r.store[clientID]

	if !exists {
		return nil, false, nil
	}

	copy := *rateLimit
	return &copy, true, nil
}

func (r *memoryRateLimiterRepository) Save(ctx context.Context, rateLimit *domain.RateLimit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	copy := *rateLimit
	r.store[rateLimit.ClientID] = &copy
	return nil
}

func (r *memoryRateLimiterRepository) CreateDefault(ctx context.Context, clientID string) *domain.RateLimit {
	return &domain.RateLimit{
		ClientID:      clientID,
		RequestCount:  0,
		CycleStart:    time.Now(),
		CycleDuration: r.defaultCycleDuration,
		MaxRequests:   r.defaultMaxRequest,
	}
}

func (r *memoryRateLimiterRepository) Delete(ctx context.Context, clientID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.store, clientID)
	return nil
}
