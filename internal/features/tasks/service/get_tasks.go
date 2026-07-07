package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_pagination "github.com/Slazzzer/golang-todoapp/internal/core/pagination"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	resolvedLimit, resolvedOffset, err := core_pagination.Resolve(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("resolve pagination: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, resolvedLimit, resolvedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}
