package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		DELETE FROM todoapp.users
		WHERE user_id = $1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		if core_postgres_pool.IsErrViolatesForeignKey(err) {
			return fmt.Errorf(
				"user with id='%d' has related tasks: %w",
				id,
				core_errors.ErrConflict,
			)
		}
		return fmt.Errorf("exec query: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil
}
