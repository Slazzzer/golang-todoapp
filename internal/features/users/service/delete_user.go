package users_service

import (
	"context"
	"fmt"

	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *UsersService) DeleteUser(ctx context.Context, userID int) error {
	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get requester id: %w", err)
	}

	if err := core_auth.EnsureSelfAccess(requesterID, userID); err != nil {
		return fmt.Errorf("ensure self access: %w", err)
	}

	if err := s.usersRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
