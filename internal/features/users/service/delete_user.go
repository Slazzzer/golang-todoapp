package users_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/actinguser"
)

func (s *UsersService) DeleteUser(ctx context.Context, userID int) error {
	if err := actinguser.EnsureSelf(ctx, userID); err != nil {
		return err
	}

	if err := s.usersRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
