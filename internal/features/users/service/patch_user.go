package users_service

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch) (domain.User, error) {

	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return domain.User{}, fmt.Errorf("get requester id: %w", err)
	}

	if err := core_auth.EnsureSelfAccess(requesterID, id); err != nil {
		return domain.User{}, fmt.Errorf("ensure self access: %w", err)
	}

	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to get user: %w",
			err,
		)
	}

	if err := user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to apply patch to user: %w",
			err,
		)
	}

	patchedUser, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to patch user: %w",
			err,
		)
	}
	return patchedUser, nil
}
