package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/actinguser"
	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	actingID, err := actinguser.Require(ctx)
	if err != nil {
		return domain.Task{}, err
	}

	task.AuthorUserID = actingID

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("invalid task domain: %w", err)
	}

	task, err = s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
