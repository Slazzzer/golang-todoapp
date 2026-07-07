package core_http_middleware

import "net/http"

func LimitBody(maxBytes int64) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(rw, r.Body, maxBytes)
			next.ServeHTTP(rw, r)
		})
	}
}
