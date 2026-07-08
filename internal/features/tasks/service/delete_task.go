package tasks_service

import (
	"context"
	"fmt"

	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	taskID int,
) error {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get requester id: %w", err)
	}

	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task: %w", err)
	}

	if err := core_auth.EnsureTaskOwner(requesterID, task.AuthorUserID); err != nil {
		return fmt.Errorf("ensure task owner: %w", err)
	}

	if err := s.tasksRepository.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}
	return nil
}
