package core_http_middleware

import (
	"encoding/json"
	"net/http"

	core_ratelimit "github.com/Slazzzer/golang-todoapp/internal/core/ratelimit"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
)

func RateLimit(limiter *core_ratelimit.Limiter) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if !limiter.Allow(core_http_request.ClientIP(r)) {
				writeTooManyRequests(rw)
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}

func writeTooManyRequests(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusTooManyRequests)
	_ = json.NewEncoder(rw).Encode(map[string]string{
		"message": "too many requests",
		"code":    "too_many_requests",
	})
}
