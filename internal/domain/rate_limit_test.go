package domain

import (
	"testing"
	"time"
)

func TestRateLimit_IsAllowed(t *testing.T) {
	t.Run("allow when under limit", func(t *testing.T) {
		rateLimit := &RateLimit{
			ClientID:      "test-client",
			RequestCount:  5,
			MaxRequests:   10,
			CycleDuration: 1,
			CycleStart:    time.Now(),
		}

		allowed := rateLimit.IsAllowed()
		if !allowed {
			t.Error("Expected to allow request, but got blocked")
		}
	})

	t.Run("block when at limit", func(t *testing.T) {
		rateLimit := &RateLimit{
			ClientID:      "test-client",
			RequestCount:  10,
			MaxRequests:   10,
			CycleDuration: 1,
			CycleStart:    time.Now(),
		}

		allowed := rateLimit.IsAllowed()
		if allowed {
			t.Error("Expected to block request, but got allowed")
		}
	})

	t.Run("reset after time passed", func(t *testing.T) {
		rateLimit := &RateLimit{
			ClientID:      "test-client",
			RequestCount:  10,
			MaxRequests:   10,
			CycleDuration: 1,
			CycleStart:    time.Now().Add(-2 * time.Minute),
		}

		allowed := rateLimit.IsAllowed()
		if !allowed {
			t.Error("Expected to allow after reset, but got blocked")
		}

		if rateLimit.RequestCount != 0 {
			t.Errorf("Expected count to be 0 after reset, got %d", rateLimit.RequestCount)
		}
	})
}

func TestRateLimit_Increment(t *testing.T) {
	rateLimit := &RateLimit{
		RequestCount: 5,
		MaxRequests:  10,
	}

	before := rateLimit.RequestCount

	rateLimit.Increment()

	after := rateLimit.RequestCount

	if after != before+1 {
		t.Errorf("Expected count to be %d, got %d", before+1, after)
	}
}

func TestRateLimit_GetRemainingRequests(t *testing.T) {
	t.Run("normal remaining", func(t *testing.T) {
		rateLimit := &RateLimit{
			RequestCount: 3,
			MaxRequests:  10,
		}

		remaining := rateLimit.GetRemainingRequests()
		expected := 7

		if remaining != expected {
			t.Errorf("Expected %d remaining, got %d", expected, remaining)
		}
	})

	t.Run("over limit returns zero", func(t *testing.T) {
		rateLimit := &RateLimit{
			RequestCount: 15,
			MaxRequests:  10,
		}

		remaining := rateLimit.GetRemainingRequests()

		if remaining != 0 {
			t.Errorf("Expected 0 remaining when over limit, got %d", remaining)
		}
	})
}

func TestRateLimit_GetResetTime(t *testing.T) {
	cycleStart := time.Now()
	rateLimit := &RateLimit{
		CycleDuration: 5,
		CycleStart:    cycleStart,
	}

	resetTime := rateLimit.GetResetTime()

	expectedReset := cycleStart.Add(5 * time.Minute)

	if resetTime.Unix() != expectedReset.Unix() {
		t.Errorf("Reset time not correct. Expected %v, got %v",
			expectedReset, resetTime)
	}
}
