package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	taskID int,
) (domain.Task, error) {
	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task from repository: %w", err)
	}

	return task, nil
}
