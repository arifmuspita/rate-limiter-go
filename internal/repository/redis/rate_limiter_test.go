package redis

import (
	"context"
	"rate-limiter-go/internal/domain"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func setupTestRedis(t *testing.T) (*redis.Client, func()) {

	mr, err := miniredis.Run()
	if err != nil {
		t.Fatal("Failed to start miniredis:", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	return client, func() {
		client.Close()
		mr.Close()
	}
}

func TestRedisGet(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewRateLimiterRedisRepository(client, 100, 1)
	ctx := context.Background()

	_, exists, err := repo.Get(ctx, "not-exist")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if exists {
		t.Error("Expected exists to be false")
	}

	rateLimit := &domain.RateLimit{
		ClientID:      "test-client",
		RequestCount:  5,
		MaxRequests:   100,
		CycleDuration: 1,
		CycleStart:    time.Now(),
	}

	err = repo.Save(ctx, rateLimit)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	got, exists, err := repo.Get(ctx, "test-client")
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if !exists {
		t.Error("Expected data to exist")
	}
	if got.ClientID != "test-client" {
		t.Error("Wrong client ID")
	}
	if got.RequestCount != 5 {
		t.Errorf("Expected count 5, got %d", got.RequestCount)
	}
}

func TestRedisSave(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewRateLimiterRedisRepository(client, 100, 1)
	ctx := context.Background()

	rateLimit := &domain.RateLimit{
		ClientID:      "save-test",
		RequestCount:  10,
		MaxRequests:   100,
		CycleDuration: 1,
		CycleStart:    time.Now(),
	}

	err := repo.Save(ctx, rateLimit)
	if err != nil {
		t.Errorf("Save failed: %v", err)
	}

	got, exists, _ := repo.Get(ctx, "save-test")
	if !exists {
		t.Error("Data should exist after save")
	}
	if got.RequestCount != 10 {
		t.Error("Wrong data saved")
	}
}

func TestRedisDelete(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewRateLimiterRedisRepository(client, 100, 1)
	ctx := context.Background()

	rateLimit := &domain.RateLimit{
		ClientID:      "delete-test",
		RequestCount:  5,
		MaxRequests:   100,
		CycleDuration: 1,
		CycleStart:    time.Now(),
	}
	repo.Save(ctx, rateLimit)

	err := repo.Delete(ctx, "delete-test")
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, exists, _ := repo.Get(ctx, "delete-test")
	if exists {
		t.Error("Data should not exist after delete")
	}
}

func TestRedisCreateDefault(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewRateLimiterRedisRepository(client, 50, 5)
	ctx := context.Background()

	defaultRL := repo.CreateDefault(ctx, "new-client")

	if defaultRL.ClientID != "new-client" {
		t.Error("Wrong client ID")
	}
	if defaultRL.RequestCount != 0 {
		t.Error("Count should be 0")
	}
	if defaultRL.MaxRequests != 50 {
		t.Error("Wrong max requests")
	}
	if defaultRL.CycleDuration != 5 {
		t.Error("Wrong cycle duration")
	}
}

func TestRedisExpiration(t *testing.T) {
	client, cleanup := setupTestRedis(t)
	defer cleanup()

	repo := NewRateLimiterRedisRepository(client, 100, 1)
	ctx := context.Background()

	rateLimit := &domain.RateLimit{
		ClientID:      "expire-test",
		RequestCount:  5,
		MaxRequests:   100,
		CycleDuration: 1,
		CycleStart:    time.Now(),
	}

	err := repo.Save(ctx, rateLimit)
	if err != nil {
		t.Fatal("Save failed:", err)
	}

	key := "rate_limit:expire-test"
	ttl, err := client.TTL(ctx, key).Result()
	if err != nil {
		t.Fatal("Failed to get TTL:", err)
	}

	if ttl.Seconds() < 60 || ttl.Seconds() > 180 {
		t.Errorf("Unexpected TTL: %v", ttl)
	}
}
