package memory

import (
	"context"
	"rate-limiter-go/internal/domain"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	repo := NewRateLimiterMemoryRepository(100, 1)
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
		t.Errorf("Save failed: %v", err)
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

func TestSave(t *testing.T) {
	repo := NewRateLimiterMemoryRepository(100, 1)
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

func TestDelete(t *testing.T) {
	repo := NewRateLimiterMemoryRepository(100, 1)
	ctx := context.Background()

	rateLimit := &domain.RateLimit{
		ClientID: "delete-test",
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

func TestCreateDefault(t *testing.T) {
	repo := NewRateLimiterMemoryRepository(50, 5)
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
	if defaultRL.CycleStart.IsZero() {
		t.Error("Cycle start should not be zero")
	}
}
