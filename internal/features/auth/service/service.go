package auth_service

import (
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/adminauth"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

type LoginResult struct {
	Role    string
	Session string
}

type AuthService struct {
	adminCfg adminauth.Config
}

func NewAuthService(adminCfg adminauth.Config) *AuthService {
	return &AuthService{adminCfg: adminCfg}
}

func (s *AuthService) Login(login, password string) (LoginResult, error) {
	if !s.adminCfg.ValidateCredentials(login, password) {
		return LoginResult{}, fmt.Errorf("invalid credentials: %w", core_errors.ErrUnauthorized)
	}

	return LoginResult{
		Role:    "admin",
		Session: s.adminCfg.SessionToken,
	}, nil
}
