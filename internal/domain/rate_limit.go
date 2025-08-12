package domain

import "time"

type RateLimit struct {
	ClientID      string
	RequestCount  int
	MaxRequests   int
	CycleDuration int
	CycleStart    time.Time
}

func (r *RateLimit) IsAllowed() bool {
	now := time.Now()
	duration := time.Duration(r.CycleDuration) * time.Minute

	if now.Sub(r.CycleStart) >= duration {
		r.CycleStart = now
		r.RequestCount = 0
	}

	return r.RequestCount < r.MaxRequests
}

func (r *RateLimit) Increment() {
	r.RequestCount++
}

func (r *RateLimit) GetRemainingRequests() int {
	remaining := r.MaxRequests - r.RequestCount
	if remaining < 0 {
		return 0
	}
	return remaining
}

func (r *RateLimit) GetResetTime() time.Time {
	duration := time.Duration(r.CycleDuration) * time.Minute
	return r.CycleStart.Add(duration)
}
