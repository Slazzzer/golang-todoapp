package web_transport_http

import (
	"net/http"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

func (h *WebHTTPHandler) GetMainPage(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	html, err := h.webService.GetMainPage()
	if err != nil {
		responseHandler.ErrorResponse(err, "get_main_page from service")
		return
	}

	responseHandler.HTMLResponse(html)

}
