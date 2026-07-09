package users_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_pagination "github.com/Slazzzer/golang-todoapp/internal/core/pagination"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) (core_pagination.Page[domain.User], error) {
	resolvedLimit, resolvedOffset, err := core_pagination.Resolve(limit, offset)
	if err != nil {
		return core_pagination.Page[domain.User]{}, fmt.Errorf("resolve pagination: %w", err)
	}

	total, err := s.usersRepository.CountUsers(ctx)
	if err != nil {
		return core_pagination.Page[domain.User]{}, fmt.Errorf("count users: %w", err)
	}

	users, err := s.usersRepository.GetUsers(ctx, resolvedLimit, resolvedOffset)
	if err != nil {
		return core_pagination.Page[domain.User]{}, fmt.Errorf("get users from repository: %w", err)
	}

	return core_pagination.NewPage(users, total, resolvedLimit, resolvedOffset), nil
}
