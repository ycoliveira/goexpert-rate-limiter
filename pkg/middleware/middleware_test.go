package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	blocked map[string]time.Time
}

func NewMockStorage() *MockStorage {
	return &MockStorage{blocked: make(map[string]time.Time)}
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

func TestRateLimiter(t *testing.T) {
	storage := NewMockStorage()
	rateLimiter := NewRateLimiter(1, 2*time.Second, storage)

	handler := rateLimiter.RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	}))

	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// First request should pass
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Second request should be blocked
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// After 2 seconds, request should pass again
	time.Sleep(2 * time.Second)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}
