package auth_transport_http

import (
	"context"
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_http_middleware "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/middleware"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService        AuthService
	registerMiddleware []core_http_middleware.Middleware
}

type AuthService interface {
	Register(ctx context.Context, user domain.User) (domain.User, string, error)
	Me(ctx context.Context, userID int) (domain.User, error)
}

func NewAuthHTTPHandler(
	authService AuthService,
	registerMiddleware ...core_http_middleware.Middleware,
) *AuthHTTPHandler {
	return &AuthHTTPHandler{
		authService:        authService,
		registerMiddleware: registerMiddleware,
	}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(
			http.MethodPost,
			"/auth/register",
			h.Register,
			h.registerMiddleware...,
		),
		core_http_server.NewRoute(http.MethodGet, "/auth/me", h.Me),
	}
}
