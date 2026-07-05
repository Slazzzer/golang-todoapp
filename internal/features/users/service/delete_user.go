package users_service

import (
	"context"
	"fmt"
)

func (s *UsersService) DeleteUser(ctx context.Context, userID int) error {
	if err := s.usersRepository.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}
