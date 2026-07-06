package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) GetTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT
			task_id,
			task_version,
			task_title,
			task_description,
			task_completed,
			task_created_at,
			task_completed_at,
			author_user_id
		FROM todoapp.tasks
		WHERE task_id = $1;
		`

	row := r.pool.QueryRow(ctx, query, taskID)

	var taskModel TaskModel

	err := row.Scan(
		&taskModel.ID,
		&taskModel.Version,
		&taskModel.Title,
		&taskModel.Description,
		&taskModel.Completed,
		&taskModel.CreatedAt,
		&taskModel.CompletedAt,
		&taskModel.AuthorUserID,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Task{}, fmt.Errorf(
				"task with id='%d': %w",
				taskID,
				core_errors.ErrNotFound,
			)
		}
		return domain.Task{}, fmt.Errorf("failed to scan task model: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)

	return taskDomain, nil
}
