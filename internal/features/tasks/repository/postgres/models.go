package tasks_postgres_repository

import (
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

type TaskModel struct {
	ID           int        `db:"task_id"`
	Version      int        `db:"task_version"`
	Title        string     `db:"task_title"`
	Description  *string    `db:"task_description"`
	Completed    bool       `db:"task_completed"`
	CreatedAt    time.Time  `db:"task_created_at"`
	CompletedAt  *time.Time `db:"task_completed_at"`
	AuthorUserID int        `db:"author_user_id"`
}

func taskDomainFromModel(taskModel TaskModel) domain.Task {
	return domain.NewTask(
		taskModel.ID,
		taskModel.Version,
		taskModel.Title,
		taskModel.Description,
		taskModel.Completed,
		taskModel.CreatedAt,
		taskModel.CompletedAt,
		taskModel.AuthorUserID,
	)
}

func taskDomainsFromModels(taskModels []TaskModel) []domain.Task {
	domains := make([]domain.Task, len(taskModels))

	for i, model := range taskModels {
		domains[i] = taskDomainFromModel(model)
	}

	return domains
}
