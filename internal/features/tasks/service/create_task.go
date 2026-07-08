package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *TasksService) CreateTask(
	ctx context.Context,
	task domain.Task,
) (domain.Task, error) {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get requester id: %w", err)
	}

	task.AuthorUserID = requesterID

	if err := task.Validate(); err != nil {
		return domain.Task{}, fmt.Errorf("invalid task domain: %w", err)
	}

	task, err = s.tasksRepository.CreateTask(ctx, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
