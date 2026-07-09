package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/types"
)

// PatchTaskRequest тело частичного обновления задачи.
type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`       // Новый заголовок (опционально)
	Description core_http_types.Nullable[string] `json:"description"` // Новое описание (опционально)
	Completed   core_http_types.Nullable[bool]   `json:"completed"`   // Статус выполнения (опционально)
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` is required")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 characters")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("`Description` must be between 1 and 1000 characters")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` is required")
		}
	}

	return nil
}

type PatchTaskUserResponse TaskDTOResponse

// PatchTask обновляет задачу.
//
// @Summary      Изменить задачу
// @Description  Частичное обновление. При completed=true автоматически проставляется completed_at.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID задачи"
// @Param        request body PatchTaskRequest true "Поля для обновления"
// @Success      200 {object} PatchTaskUserResponse "Обновлённая задача"
// @Failure      400 {object} map[string]string "Невалидные данные"
// @Failure      404 {object} map[string]string "Задача не найдена"
// @Router       /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task ID from path")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request body",
		)
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
		return
	}

	response := PatchTaskUserResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
