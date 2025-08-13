package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"rate-limiter-go/internal/domain"
	"rate-limiter-go/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRateLimiterRepository struct {
	client               *redis.Client
	defaultMaxRequests   int
	defaultCycleDuration int
}

func NewRateLimiterRedisRepository(client *redis.Client, defaultMaxRequests, defaultCycleDuration int) repository.RateLimiterRepository {
	return &redisRateLimiterRepository{
		client:               client,
		defaultMaxRequests:   defaultMaxRequests,
		defaultCycleDuration: defaultCycleDuration,
	}
}

func (r *redisRateLimiterRepository) Get(ctx context.Context, clientID string) (*domain.RateLimit, bool, error) {
	key := fmt.Sprintf("rate_limit:%s", clientID)

	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("failed to get redis: %w", err)
	}

	var rateLimit domain.RateLimit
	if err := json.Unmarshal([]byte(data), &rateLimit); err != nil {
		return nil, false, fmt.Errorf("failed to unmarshal: %w", err)
	}

	return &rateLimit, true, nil
}

func (r *redisRateLimiterRepository) Save(ctx context.Context, rateLimit *domain.RateLimit) error {
	key := fmt.Sprintf("rate_limit:%s", rateLimit.ClientID)

	data, err := json.Marshal(rateLimit)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	expiration := time.Duration(rateLimit.CycleDuration*2) * time.Minute
	if err := r.client.Set(ctx, key, data, expiration).Err(); err != nil {
		return fmt.Errorf("failed save to redis: %w", err)
	}

	return nil
}

func (r *redisRateLimiterRepository) Delete(ctx context.Context, clientID string) error {
	key := fmt.Sprintf("rate_limit:%s", clientID)

	if err := r.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete from redis: %w", err)
	}

	return nil
}

func (r *redisRateLimiterRepository) CreateDefault(ctx context.Context, clientID string) *domain.RateLimit {
	return &domain.RateLimit{
		ClientID:      clientID,
		RequestCount:  0,
		CycleStart:    time.Now(),
		CycleDuration: r.defaultCycleDuration,
		MaxRequests:   r.defaultMaxRequests,
	}
}
