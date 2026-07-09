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
) (core_pagination.Page[domain.Task], error) {
	resolvedLimit, resolvedOffset, err := core_pagination.Resolve(limit, offset)
	if err != nil {
		return core_pagination.Page[domain.Task]{}, fmt.Errorf("resolve pagination: %w", err)
	}

	total, err := s.tasksRepository.CountTasks(ctx, userID)
	if err != nil {
		return core_pagination.Page[domain.Task]{}, fmt.Errorf("count tasks: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, userID, resolvedLimit, resolvedOffset)
	if err != nil {
		return core_pagination.Page[domain.Task]{}, fmt.Errorf("failed to get tasks: %w", err)
	}

	return core_pagination.NewPage(tasks, total, resolvedLimit, resolvedOffset), nil
}
