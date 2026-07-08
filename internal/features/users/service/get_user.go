package users_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("get requester id: %w", err)
	}

	if err := core_auth.EnsureSelfAccess(requesterID, id); err != nil {
		return domain.User{}, fmt.Errorf("ensure self access: %w", err)
	}

	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user by ID from repository: %w", err)
	}

	return user, nil
}
