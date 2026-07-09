package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse TasksPageResponse

// GetTasks возвращает список задач.
//
// @Summary      Список задач
// @Description  Возвращает страницу задач с total/limit/offset. Можно отфильтровать по user_id (по умолчанию limit=50, max=100).
// @Tags         tasks
// @Produce      json
// @Param        user_id query int false "Фильтр по автору"
// @Param        limit query int false "Размер страницы"
// @Param        offset query int false "Смещение"
// @Success      200 {object} GetTasksResponse "Страница задач"
// @Failure      400 {object} map[string]string "Невалидные параметры"
// @Router       /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDAndLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user id and limit offset query params",
		)
		return
	}

	page, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(tasksPageFromDomains(page))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDAndLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		queryParamUserID = "user_id"
		queryParamLimit  = "limit"
		queryParamOffset = "offset"
	)

	userID, err := core_http_request.GetIntQueryParam(r, queryParamUserID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get user id query param: %w", err)
	}

	limit, err := core_http_request.GetIntQueryParam(r, queryParamLimit)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, queryParamOffset)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get offset query param: %w", err)
	}

	return userID, limit, offset, nil
}
