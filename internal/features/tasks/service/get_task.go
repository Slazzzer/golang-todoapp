package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get requester id: %w", err)
	}

	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task from repository: %w", err)
	}

	if err := core_auth.EnsureTaskOwner(requesterID, task.AuthorUserID); err != nil {
		return domain.Task{}, fmt.Errorf("ensure task owner: %w", err)
	}

	return task, nil
}
