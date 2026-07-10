package core_http_middleware

import (
	"net/http"
	"time"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

func CORS(cfg CORSConfig) Middleware {
	checker := cfg.checker()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && checker.isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
			}

			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Admin-Session")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}

			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)

			l := log.With(
				core_logger.String("request_id", requestID),
				core_logger.String("url", r.URL.String()),
			)

			ctx := core_logger.ContextWithLogger(r.Context(), l)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHTTPResponseHandler(log, w)

			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "During handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)

			before := time.Now()

			log.Debug(
				">>> incoming HTTP request",
				core_logger.String("http_method", r.Method),
				core_logger.Time("time", before.UTC()),
			)

			defer func() {
				log.Debug(
					"<<< done HTTP request",
					core_logger.Int("status_code", rw.GetStatusCode()),
					core_logger.Duration("latency", time.Since(before)),
				)
			}()

			next.ServeHTTP(rw, r)
		})
	}
}

// Dummy — демонстрационный middleware: логирует вход/выход в цепочку до/после handler.
func Dummy() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := core_logger.FromContext(r.Context())

			log.Debug("[Dummy] before handler")

			next.ServeHTTP(w, r)

			log.Debug("[Dummy] after handler")
		})
	}
}
