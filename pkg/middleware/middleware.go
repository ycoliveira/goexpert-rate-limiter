package middleware

import (
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter     *rate.Limiter
	storage     Storage
	blockDur    time.Duration
	tokenLimits map[string]int
}

type Storage interface {
	IsBlocked(key string) bool
	Block(key string, duration time.Duration)
	GetLimiter(key string) *rate.Limiter
	SetLimiter(key string, limiter *rate.Limiter)
}

func NewRateLimiter(maxRequests int, blockDur time.Duration, storage Storage, tokenLimits map[string]int) *RateLimiter {
	return &RateLimiter{
		limiter:     rate.NewLimiter(rate.Limit(maxRequests), maxRequests),
		storage:     storage,
		blockDur:    blockDur,
		tokenLimits: tokenLimits,
	}
}

func getIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip == "" {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return ip
}

func (rl *RateLimiter) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		token := r.Header.Get("API_KEY")
		var limiter *rate.Limiter
		var key string

		log.Printf("Received token: %s", token)

		if token != "" {
			key = token
			if limit, ok := rl.tokenLimits[token]; ok {
				limiter = rl.storage.GetLimiter(key)
				if limiter == nil {
					limiter = rate.NewLimiter(rate.Limit(limit), limit)
					rl.storage.SetLimiter(key, limiter)
				}
			} else {
				http.Error(w, "invalid token", http.StatusForbidden)
				return
			}
		} else {
			key = ip
			limiter = rl.limiter
		}

		if rl.storage.IsBlocked(key) {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		if !limiter.Allow() {
			rl.storage.Block(key, rl.blockDur)
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
