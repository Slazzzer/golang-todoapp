package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks возвращает список задач текущего пользователя.
//
// @Summary      Список задач
// @Description  Возвращает только задачи авторизованного пользователя. Поддерживает пагинацию limit/offset (по умолчанию limit=50, max=100).
// @Tags         tasks
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Размер страницы"
// @Param        offset query int false "Смещение"
// @Success      200 {array} TaskDTOResponse "Список задач"
// @Failure      401 {object} map[string]string "Не авторизован"
// @Router       /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit and offset query params",
		)
		return
	}

	tasksDomains, err := h.tasksService.GetTasks(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasks",
		)
		return
	}

	response := GetTasksResponse(taskDTOsFromDomains(tasksDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		queryParamLimit  = "limit"
		queryParamOffset = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, queryParamLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, queryParamOffset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get offset query param: %w", err)
	}

	return limit, offset, nil
}
