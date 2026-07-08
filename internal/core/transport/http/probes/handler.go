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

// ProbeStatusResponse статус пробы (health/ready).
type ProbeStatusResponse struct {
	Status string `json:"status" example:"ok"` // Статус: ok, ready или not ready
}

// Health проверка живости процесса.
//
// @Summary      Проверка живости
// @Description  Возвращает ok, если HTTP-сервер запущен. Не требует авторизации.
// @Tags         probes
// @Produce      json
// @Success      200 {object} ProbeStatusResponse "Сервер работает"
// @Router       /health [get]
func (h *Handler) health(rw http.ResponseWriter, _ *http.Request) {
	writeJSON(rw, http.StatusOK, map[string]string{"status": "ok"})
}

// Ready проверка готовности (включая PostgreSQL).
//
// @Summary      Проверка готовности
// @Description  Возвращает ready, если PostgreSQL доступен. Не требует авторизации.
// @Tags         probes
// @Produce      json
// @Success      200 {object} ProbeStatusResponse "Приложение готово принимать запросы"
// @Failure      503 {object} ProbeStatusResponse "База данных недоступна"
// @Router       /ready [get]
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
