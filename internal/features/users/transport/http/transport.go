package users_transport_http

import (
	"context"
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
)

type UsersHTTPHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user domain.User) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		userID int) (domain.User, error)

	DeleteUser(
		ctx context.Context,
		userID int) error

	PatchUser(
		ctx context.Context,
		userID int,
		userPatch domain.UserPatch) (domain.User, error)
}

func NewUsersHTTPHandler(usersService UsersService) *UsersHTTPHandler {
	return &UsersHTTPHandler{
		usersService: usersService,
	}
}

func (h *UsersHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodPost, "/users", h.CreateUser),
		core_http_server.NewRoute(http.MethodGet, "/users", h.GetUsers),
		core_http_server.NewRoute(http.MethodGet, "/users/{id}", h.GetUser),
		core_http_server.NewRoute(http.MethodDelete, "/users/{id}", h.DeleteUser),
		core_http_server.NewRoute(http.MethodPatch, "/users/{id}", h.PatchUser),
	}
}
