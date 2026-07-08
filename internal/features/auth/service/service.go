package auth_service

import (
	"context"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

type AuthService struct {
	usersService UsersService
	tokenManager TokenManager
}

type UsersService interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id int) (domain.User, error)
}

type TokenManager interface {
	Issue(userID int) (string, error)
}

func NewAuthService(
	usersService UsersService,
	tokenManager TokenManager,
) *AuthService {
	return &AuthService{
		usersService: usersService,
		tokenManager: tokenManager,
	}
}
