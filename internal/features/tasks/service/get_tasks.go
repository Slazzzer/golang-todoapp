package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
	core_pagination "github.com/Slazzzer/golang-todoapp/internal/core/pagination"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get requester id: %w", err)
	}

	resolvedLimit, resolvedOffset, err := core_pagination.Resolve(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("resolve pagination: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, &requesterID, resolvedLimit, resolvedOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}
