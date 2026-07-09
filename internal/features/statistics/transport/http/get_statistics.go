package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

// GetStatisticsResponse статистика по задачам пользователя.
type GetStatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created" example:"10"`              // Создано задач за период
	TasksCompleted             int      `json:"tasks_completed" example:"7"`             // Завершено задач за период
	TaskCompletedRate          *float64 `json:"task_completed_rate" example:"70"`     // Процент завершённых (0–100) или null
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"24h30m15s"` // Среднее время выполнения или null
}

// GetStatistics возвращает статистику по задачам.
//
// @Summary      Статистика
// @Description  Агрегированная статистика по задачам. Параметр user_id опционален: без него — по всем пользователям.
// @Tags         statistics
// @Produce      json
// @Param        user_id query int false "Фильтр по пользователю"
// @Param        from query string false "Начало периода (YYYY-MM-DD)"
// @Param        to query string false "Конец периода (YYYY-MM-DD, не включая этот день)"
// @Success      200 {object} GetStatisticsResponse "Статистика"
// @Failure      400 {object} map[string]string "Невалидный период или параметры"
// @Router       /statistics [get]
func (h *StatisticsHTTPHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID/from/to query params")
		return
	}

	statistics, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get statistics")
		return
	}

	response := toDTOFromDomain(statistics)
	responseHandler.JSONResponse(response, http.StatusOK)

}

func toDTOFromDomain(statistics domain.Statistics) GetStatisticsResponse {
	var avgTime *string
	if statistics.TasksAverageCompletionTime != nil {
		duration := statistics.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return GetStatisticsResponse{
		TasksCreated:               statistics.TasksCreated,
		TasksCompleted:             statistics.TasksCompleted,
		TaskCompletedRate:          statistics.TaskCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {

	const (
		userIDQueryParam  = "user_id"
		fromQueryParamKey = "from"
		toQueryParamKey   = "to"
	)

	userID, err := core_http_request.GetIntQueryParam(r, userIDQueryParam)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user ID from query params: %w", err)
	}

	from, err := core_http_request.GetDateFromQueryParam(r, fromQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get from time from query params: %w", err)
	}

	to, err := core_http_request.GetDateToQueryParam(r, toQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get to time from query params: %w", err)
	}

	return userID, from, to, nil

}
