package web_transport_http

import (
	"net/http"
	"path"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetAsset(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	name := path.Base(r.URL.Path)

	content, contentType, err := h.webService.GetAsset(name)
	if err != nil {
		responseHandler.ErrorResponse(err, "get_asset from service")
		return
	}

	responseHandler.AssetResponse(content, contentType)
}
