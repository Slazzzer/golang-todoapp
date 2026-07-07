package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
)

type HTTPResponseHandler struct {
	log core_logger.Logger
	rw  http.ResponseWriter
}

type errorBody struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

func NewHTTPResponseHandler(log core_logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{log: log, rw: rw}
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", core_logger.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(msg string, fields ...core_logger.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, core_logger.Error(err))

	h.writeErrorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, core_logger.Error(err))

	h.writeErrorResponse(http.StatusInternalServerError, nil, "internal server error")
}

func (h *HTTPResponseHandler) writeErrorResponse(statusCode int, err error, msg string) {
	body := errorBody{Message: msg}

	if statusCode >= http.StatusInternalServerError {
		body.Message = "internal server error"
		body.Code = ""
	} else if err != nil {
		body.Code = errorCode(err)
	}

	h.JSONResponse(body, statusCode)
}

func errorCode(err error) string {
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		return "invalid_argument"
	case errors.Is(err, core_errors.ErrNotFound):
		return "not_found"
	case errors.Is(err, core_errors.ErrConflict):
		return "conflict"
	default:
		return ""
	}
}
