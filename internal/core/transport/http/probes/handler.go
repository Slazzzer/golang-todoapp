package core_http_probes

import (
	"encoding/json"
	"net/http"

	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
)

type Handler struct {
	pool core_postgres_pool.Pool
}

func NewHandler(pool core_postgres_pool.Pool) *Handler {
	return &Handler{pool: pool}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("GET /ready", h.ready)
}

func (h *Handler) health(rw http.ResponseWriter, _ *http.Request) {
	writeJSON(rw, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) ready(rw http.ResponseWriter, r *http.Request) {
	if err := h.pool.Ping(r.Context()); err != nil {
		writeJSON(rw, http.StatusServiceUnavailable, map[string]string{
			"status": "not ready",
		})
		return
	}

	writeJSON(rw, http.StatusOK, map[string]string{"status": "ready"})
}

func writeJSON(rw http.ResponseWriter, statusCode int, body any) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	_ = json.NewEncoder(rw).Encode(body)
}
