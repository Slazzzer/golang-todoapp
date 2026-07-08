package core_http_middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

var publicRoutes = map[string]struct{}{
	"POST /auth/register": {},
}

func JWTAuth(tokenManager *core_auth.TokenManager) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			routeKey := r.Method + " " + r.URL.Path
			if _, ok := publicRoutes[routeKey]; ok {
				next.ServeHTTP(rw, r)
				return
			}

			authHeader := r.Header.Get("Authorization")
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeUnauthorized(rw)
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			if token == "" {
				writeUnauthorized(rw)
				return
			}

			userID, err := tokenManager.ParseUserID(token)
			if err != nil {
				writeUnauthorized(rw)
				return
			}

			ctx := core_auth.ContextWithUserID(r.Context(), userID)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func writeUnauthorized(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(rw).Encode(map[string]string{
		"message": "authentication required",
		"code":    "unauthorized",
	})
}
