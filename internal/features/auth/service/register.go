package auth_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

func (s *AuthService) Register(
	ctx context.Context,
	user domain.User,
) (domain.User, string, error) {
	createdUser, err := s.usersService.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, "", fmt.Errorf("create user: %w", err)
	}

	token, err := s.tokenManager.Issue(createdUser.ID)
	if err != nil {
		return domain.User{}, "", fmt.Errorf("issue token: %w", err)
	}

	return createdUser, token, nil
}

func (s *AuthService) Me(ctx context.Context, userID int) (domain.User, error) {
	user, err := s.usersService.GetUser(ctx, userID)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user: %w", err)
	}

	return user, nil
}
