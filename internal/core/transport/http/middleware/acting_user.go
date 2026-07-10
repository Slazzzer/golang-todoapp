package core_http_middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Slazzzer/golang-todoapp/internal/core/actinguser"
	"github.com/Slazzzer/golang-todoapp/internal/core/adminauth"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

const (
	ActingUserHeader   = "X-User-ID"
	AdminSessionHeader = "X-Admin-Session"
)

// ActingUser извлекает контекст пользователя из X-User-ID или сессию администратора из X-Admin-Session.
func ActingUser(adminCfg adminauth.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path == "" {
				path = "/"
			}

			if isPublicAPIRoute(r.Method, path) {
				ctx := r.Context()
				if adminCfg.ValidateSession(strings.TrimSpace(r.Header.Get(AdminSessionHeader))) {
					ctx = actinguser.WithPrincipal(ctx, actinguser.Principal{IsAdmin: true})
				} else if userID, ok := parseActingUserID(r.Header.Get(ActingUserHeader)); ok {
					ctx = actinguser.WithContext(ctx, userID)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if adminCfg.ValidateSession(strings.TrimSpace(r.Header.Get(AdminSessionHeader))) {
				next.ServeHTTP(w, r.WithContext(
					actinguser.WithPrincipal(r.Context(), actinguser.Principal{IsAdmin: true}),
				))
				return
			}

			userID, ok := parseActingUserID(r.Header.Get(ActingUserHeader))
			if !ok {
				ctx := r.Context()
				log := core_logger.FromContext(ctx)
				responseHandler := core_http_response.NewHTTPResponseHandler(log, w)
				responseHandler.ErrorResponse(
					fmt.Errorf("missing credentials: %w", core_errors.ErrUnauthorized),
					"acting user or admin session is required",
				)
				return
			}

			next.ServeHTTP(w, r.WithContext(actinguser.WithContext(r.Context(), userID)))
		})
	}
}

func isPublicAPIRoute(method, path string) bool {
	switch {
	case method == http.MethodPost && path == "/auth/login":
		return true
	case method == http.MethodGet && path == "/users":
		return true
	case method == http.MethodPost && path == "/users":
		return true
	default:
		return false
	}
}

func parseActingUserID(raw string) (int, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}

	userID, err := strconv.Atoi(raw)
	if err != nil || userID <= 0 {
		return 0, false
	}

	return userID, true
}
