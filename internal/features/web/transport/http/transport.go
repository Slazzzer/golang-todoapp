package web_transport_http

import (
	"net/http"

	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
)

type WebHTTPHandler struct {
	webService WebService
}

type WebService interface {
	GetMainPage() ([]byte, error)
	GetAsset(name string) ([]byte, string, error)
}

func NewWebHTTPHandler(
	webService WebService,
) *WebHTTPHandler {
	return &WebHTTPHandler{
		webService: webService,
	}
}

func (h *WebHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		// Корень регистрируем без метода (совместимо с подпутями /api/v1, /swagger).
		{Path: "/", Handler: h.GetMainPage},
		core_http_server.NewRoute(http.MethodGet, "/styles.css", h.GetAsset),
		core_http_server.NewRoute(http.MethodGet, "/app.js", h.GetAsset),
	}
}
