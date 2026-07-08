package tasks_transport_http

import (
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

// TaskDTOResponse данные задачи в ответе API.
type TaskDTOResponse struct {
	ID           int        `json:"id" example:"1"`                                    // Идентификатор задачи
	Version      int        `json:"version" example:"1"`                               // Версия записи
	Title        string     `json:"title" example:"Купить молоко"`                     // Заголовок
	Description  *string    `json:"description" example:"2 литра"`                     // Описание или null
	Completed    bool       `json:"completed" example:"false"`                         // Выполнена ли задача
	CreatedAt    time.Time  `json:"created_at" example:"2026-07-08T12:00:00Z"`         // Дата создания
	CompletedAt  *time.Time `json:"completed_at"`                                      // Дата завершения или null
	AuthorUserID int        `json:"author_user_id" example:"1"`                        // ID автора (из JWT)
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		CompletedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func taskDTOsFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
