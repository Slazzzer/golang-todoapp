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
) ([]domain.User, error) {
	resolvedLimit, resolvedOffset, err := core_pagination.Resolve(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("resolve pagination: %w", err)
	}

	users, err := s.usersRepository.GetUsers(ctx, resolvedLimit, resolvedOffset)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return users, nil
}
