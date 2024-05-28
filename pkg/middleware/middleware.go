package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter  *rate.Limiter
	storage  Storage
	blockDur time.Duration
}

type Storage interface {
	IsBlocked(key string) bool
	Block(key string, duration time.Duration)
}

func NewRateLimiter(maxRequests int, blockDur time.Duration, storage Storage) *RateLimiter {
	return &RateLimiter{
		limiter:  rate.NewLimiter(rate.Limit(maxRequests), maxRequests),
		storage:  storage,
		blockDur: blockDur,
	}
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		if rl.storage.IsBlocked(ip) {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		if !rl.limiter.Allow() {
			rl.storage.Block(ip, rl.blockDur)
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
