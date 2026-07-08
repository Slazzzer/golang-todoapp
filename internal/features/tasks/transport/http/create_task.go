package tasks_transport_http

import (
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

// CreateTaskRequest тело запроса создания задачи.
type CreateTaskRequest struct {
	Title       string  `json:"title" validate:"required,min=1,max=100" example:"Купить молоко"` // Заголовок (1–100 символов)
	Description *string `json:"description" validate:"omitempty,min=1,max=1000" example:"2 литра"` // Описание (опционально)
}

type CreateTaskResponse TaskDTOResponse

// CreateTask создаёт задачу для текущего пользователя.
//
// @Summary      Создать задачу
// @Description  author_user_id берётся из JWT автоматически, передавать в JSON не нужно.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateTaskRequest true "Данные новой задачи"
// @Success      201 {object} CreateTaskResponse "Задача создана"
// @Failure      400 {object} map[string]string "Невалидные данные"
// @Failure      401 {object} map[string]string "Не авторизован"
// @Router       /tasks [post]
func (h *TasksHTTPHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request body")
		return
	}

	taskDomain := domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		0,
	)

	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create task")
		return
	}

	response := CreateTaskResponse(taskDTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}
