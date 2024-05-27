package middleware

import "golang.org/x/time/rate"

var limit = rate.NewLimiter(1, 4)
