package middleware

import (
	"net/http"
)

func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limit.Allow() {
			message := "you have reached the maximum number of requests or actions allowed within a certain time frame"
			http.Error(w, message, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
