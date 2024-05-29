package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

type MockStorage struct {
	blocked  map[string]time.Time
	limiters map[string]*rate.Limiter
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		blocked:  make(map[string]time.Time),
		limiters: make(map[string]*rate.Limiter),
	}
}

func (ms *MockStorage) IsBlocked(key string) bool {
	if unblockTime, ok := ms.blocked[key]; ok {
		if time.Now().Before(unblockTime) {
			return true
		}
		delete(ms.blocked, key)
	}
	return false
}

func (ms *MockStorage) Block(key string, duration time.Duration) {
	ms.blocked[key] = time.Now().Add(duration)
}

func (ms *MockStorage) GetLimiter(key string) *rate.Limiter {
	return ms.limiters[key]
}

func (ms *MockStorage) SetLimiter(key string, limiter *rate.Limiter) {
	ms.limiters[key] = limiter
}

func TestRateLimiter(t *testing.T) {
	tokenLimits := map[string]int{"abc123": 10}
	storage := NewMockStorage()
	rateLimiter := NewRateLimiter(5, 2*time.Second, storage, tokenLimits)

	handler := rateLimiter.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Test without token
	for i := 0; i < 6; i++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if i < 5 {
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}

	// Test with valid token
	req.Header.Set("API_KEY", "abc123")
	for i := 0; i < 11; i++ {
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if i < 10 {
			assert.Equal(t, http.StatusOK, rr.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, rr.Code)
		}
	}

	// Test with invalid token
	req.Header.Set("API_KEY", "invalid")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code)
}
