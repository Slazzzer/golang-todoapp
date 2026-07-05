package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/middleware"
)

type APIVersion string

var (
	APIVersionV1 = APIVersion("v1")
	APIVersionV2 = APIVersion("v2")
	APIVersionV3 = APIVersion("v3")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion APIVersion
	middleware []core_http_middleware.Middleware
}

func NewAPIVersionRouter(apiVersion APIVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

// Use добавляет middleware на все маршруты этой версии API.
// Вызывать до RegisterRoutes.
func (r *APIVersionRouter) Use(middleware ...core_http_middleware.Middleware) *APIVersionRouter {
	r.middleware = append(r.middleware, middleware...)
	return r
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		handler := core_http_middleware.ChainMiddleware(route.Handler, route.Middleware...)

		r.Handle(pattern, handler)
	}
}
