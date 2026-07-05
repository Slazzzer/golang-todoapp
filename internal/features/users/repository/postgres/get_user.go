package users_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	"github.com/jackc/pgx/v5"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT user_id, user_version, user_full_name, user_phone_number
		FROM todoapp.users
		WHERE user_id = $1
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.UserID,
		&userModel.UserVersion,
		&userModel.UserFullName,
		&userModel.UserPhoneNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with ID='%d': %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomain := domain.NewUser(
		userModel.UserID,
		userModel.UserVersion,
		userModel.UserFullName,
		userModel.UserPhoneNumber,
	)
	return userDomain, nil

}
