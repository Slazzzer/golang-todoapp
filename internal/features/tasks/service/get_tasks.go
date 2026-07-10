package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/actinguser"
	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_pagination "github.com/Slazzzer/golang-todoapp/internal/core/pagination"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) (core_pagination.Page[domain.Task], error) {
	filterUserID := userID

	if !actinguser.IsAdmin(ctx) {
		actingID, err := actinguser.Require(ctx)
		if err != nil {
			return core_pagination.Page[domain.Task]{}, err
		}
		filterUserID = &actingID
	}

	resolvedLimit, resolvedOffset, err := core_pagination.Resolve(limit, offset)
	if err != nil {
		return core_pagination.Page[domain.Task]{}, fmt.Errorf("resolve pagination: %w", err)
	}

	total, err := s.tasksRepository.CountTasks(ctx, filterUserID)
	if err != nil {
		return core_pagination.Page[domain.Task]{}, fmt.Errorf("count tasks: %w", err)
	}

	tasks, err := s.tasksRepository.GetTasks(ctx, filterUserID, resolvedLimit, resolvedOffset)
	if err != nil {
		return core_pagination.Page[domain.Task]{}, fmt.Errorf("failed to get tasks: %w", err)
	}

	return core_pagination.NewPage(tasks, total, resolvedLimit, resolvedOffset), nil
}
