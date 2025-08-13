package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateLimiter_Memory(t *testing.T) {
	app := SetupTestApp(t, false)

	testRateLimiterFlow(t, app)
}

func TestRateLimiter_Redis(t *testing.T) {
	app := SetupTestApp(t, true)

	testRateLimiterFlow(t, app)
}

func testRateLimiterFlow(t *testing.T, app *TestApp) {
	clientID := "test-client"

	t.Run("configure rate limit", func(t *testing.T) {
		body := map[string]int{
			"max_requests":   5,
			"cycle_duration": 1,
		}
		bodyJSON, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/rate-limit/"+clientID, bytes.NewReader(bodyJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected 200, got %d", rec.Code)
		}
	})

	t.Run("test rate limiting", func(t *testing.T) {

		for i := 0; i < 6; i++ {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/data", nil)
			req.Header.Set("X-Client-ID", clientID)
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)

			if i < 5 {

				if rec.Code != http.StatusOK {
					t.Errorf("Request %d should pass, got %d", i+1, rec.Code)
				}
			} else {

				if rec.Code != http.StatusTooManyRequests {
					t.Errorf("Request %d should be blocked, got %d", i+1, rec.Code)
				}
			}
		}
	})

	t.Run("check rate limit status", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/rate-limit/"+clientID, nil)
		rec := httptest.NewRecorder()

		app.Echo.ServeHTTP(rec, req)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)

		data := response["data"].(map[string]interface{})
		if data["allowed"].(bool) != false {
			t.Error("Should not be allowed")
		}
		if data["remaining"].(float64) != 0 {
			t.Error("Should have 0 remaining")
		}
	})
}

func TestConcurrentRequests(t *testing.T) {
	app := SetupTestApp(t, false)

	body := map[string]int{
		"max_requests":   10,
		"cycle_duration": 1,
	}
	bodyJSON, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/rate-limit/concurrent", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	app.Echo.ServeHTTP(rec, req)

	results := make(chan int, 20)

	for i := 0; i < 20; i++ {
		go func() {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/protected/data", nil)
			req.Header.Set("X-Client-ID", "concurrent")
			rec := httptest.NewRecorder()

			app.Echo.ServeHTTP(rec, req)
			results <- rec.Code
		}()
	}

	successCount := 0
	for i := 0; i < 20; i++ {
		code := <-results
		if code == http.StatusOK {
			successCount++
		}
	}

	if successCount != 10 {
		t.Errorf("Expected 10 successful requests, got %d", successCount)
	}
}
