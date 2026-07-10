package auth_transport_http

import (
	"net/http"

	auth_service "github.com/Slazzzer/golang-todoapp/internal/features/auth/service"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
)

type AuthHTTPHandler struct {
	authService AuthService
}

type AuthService interface {
	Login(login, password string) (auth_service.LoginResult, error)
}

func NewAuthHTTPHandler(authService AuthService) *AuthHTTPHandler {
	return &AuthHTTPHandler{authService: authService}
}

func (h *AuthHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodPost, "/auth/login", h.Login),
	}
}

type LoginRequest struct {
	Login    string `json:"login" validate:"required,min=1,max=50" example:"admin"`
	Password string `json:"password" validate:"required,min=1,max=100" example:"admin"`
}

type LoginResponse struct {
	Role    string `json:"role" example:"admin"`
	Session string `json:"session"`
}

func (h *AuthHTTPHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	result, err := h.authService.Login(request.Login, request.Password)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")
		return
	}

	responseHandler.JSONResponse(LoginResponse{
		Role:    result.Role,
		Session: result.Session,
	}, http.StatusOK)
}
