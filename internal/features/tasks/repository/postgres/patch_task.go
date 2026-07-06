package tasks_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *TasksRepository) PatchTask(
	ctx context.Context,
	taskID int,
	task domain.Task,
) (domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.tasks
		SET
			task_title = $1,
			task_description = $2,
			task_completed = $3,
			task_completed_at = $4,
			task_version = task_version + 1
		WHERE task_id = $5 AND task_version = $6
		RETURNING 
		    task_id,
			task_version,
			task_title,
			task_description,
			task_completed,
			task_created_at,
			task_completed_at,
			author_user_id;
	`
	row := r.pool.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		task.CompletedAt,
		taskID,
		task.Version,
	)

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
				"task with id ='%d' concurrently accessed: %w",
				taskID,
				core_errors.ErrConflict,
			)
		}
		return domain.Task{}, fmt.Errorf("failed to patch task: %w", err)
	}

	taskDomain := taskDomainFromModel(taskModel)
	return taskDomain, nil
}
