package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE todoapp.users
		SET user_full_name = $1,
		    user_phone_number = $2,
		    user_version = user_version + 1
		WHERE user_id = $3
		  AND user_version = $4
		RETURNING user_id, user_version, user_full_name, user_phone_number;
	`

	row := r.pool.QueryRow(
		ctx,
		query,
		user.FullName,
		user.PhoneNumber,
		id,
		user.Version,
	)

	var userModel UserModel
	err := row.Scan(
		&userModel.UserID,
		&userModel.UserVersion,
		&userModel.UserFullName,
		&userModel.UserPhoneNumber,
	)
	if err != nil {
		if core_postgres_pool.IsErrNoRows(err) {
			return domain.User{}, fmt.Errorf(
				"user with id='%d' concurrency accessed: %w",
				id,
				core_errors.ErrConflict,
			)
		}
		return domain.User{}, fmt.Errorf(
			"failed to scan patched user: %w",
			err,
		)
	}

	userDomain := domain.NewUser(
		userModel.UserID,
		userModel.UserVersion,
		userModel.UserFullName,
		userModel.UserPhoneNumber,
	)
	return userDomain, nil
}
