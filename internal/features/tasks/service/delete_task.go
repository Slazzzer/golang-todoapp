package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/actinguser"
)

func (s *TasksService) DeleteTask(
	ctx context.Context,
	taskID int,
) error {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return fmt.Errorf("get task before delete: %w", err)
	}

	if err := actinguser.EnsureSelf(ctx, task.AuthorUserID); err != nil {
		return err
	}

	if err := s.tasksRepository.DeleteTask(ctx, taskID); err != nil {
		return fmt.Errorf("delete task from repository: %w", err)
	}
	return nil
}
