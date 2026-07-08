package tasks_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *TasksService) PatchTask(
	ctx context.Context,
	taskID int,
	patch domain.TaskPatch,
) (domain.Task, error) {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("get requester id: %w", err)
	}

	task, err := s.tasksRepository.GetTask(ctx, taskID)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task: %w", err)
	}

	if err := core_auth.EnsureTaskOwner(requesterID, task.AuthorUserID); err != nil {
		return domain.Task{}, fmt.Errorf("ensure task owner: %w", err)
	}

	if err := task.ApplyPatch(patch); err != nil {
		return domain.Task{}, fmt.Errorf("failed to apply patch to task: %w", err)
	}

	patchedTask, err := s.tasksRepository.PatchTask(ctx, taskID, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to patch task: %w", err)
	}

	return patchedTask, nil
}
