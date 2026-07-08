package tasks_transport_http

import (
	"context"
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
)

type TasksHTTPHandler struct {
	tasksService TasksService
}

type TasksService interface {
	CreateTask(
		ctx context.Context,
		task domain.Task,
	) (domain.Task, error)
	GetTasks(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.Task, error)
	GetTask(
		ctx context.Context,
		taskID int,
	) (domain.Task, error)
	DeleteTask(
		ctx context.Context,
		taskID int,
	) error
	PatchTask(
		ctx context.Context,
		taskID int,
		patch domain.TaskPatch,
	) (domain.Task, error)
}

func NewTasksHTTPHandler(
	tasksService TasksService,
) *TasksHTTPHandler {
	return &TasksHTTPHandler{
		tasksService: tasksService,
	}
}

func (h *TasksHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodPost, "/tasks", h.CreateTask),
		core_http_server.NewRoute(http.MethodGet, "/tasks", h.GetTasks),
		core_http_server.NewRoute(http.MethodGet, "/tasks/{id}", h.GetTask),
		core_http_server.NewRoute(http.MethodDelete, "/tasks/{id}", h.DeleteTask),
		core_http_server.NewRoute(http.MethodPatch, "/tasks/{id}", h.PatchTask),
	}
}
