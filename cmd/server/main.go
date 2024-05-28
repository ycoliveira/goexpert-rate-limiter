package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ycoliveira/goexpert-rate-limiter/configs"
	middlewareRateLimiter "github.com/ycoliveira/goexpert-rate-limiter/pkg/middleware"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	maxRequests, err := strconv.Atoi(config.RateLimiterMaxRequests)
	if err != nil {
		panic(err)
	}
	blockDur := time.Duration(config.BlockTimeSeconds) * time.Second

	storage := middlewareRateLimiter.NewRedisStorage()
	rateLimiter := middlewareRateLimiter.NewRateLimiter(maxRequests, blockDur, storage)
	r.Use(rateLimiter.RateLimit)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})

	fmt.Println("Starting web server on port", config.WebServerPort)
	http.ListenAndServe(config.WebServerPort, r)
}
