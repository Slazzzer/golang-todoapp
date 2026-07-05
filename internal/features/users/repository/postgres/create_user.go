package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	user domain.User) (domain.User, error) {

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `INSERT INTO todoapp.users (user_full_name, user_phone_number)
		VALUES ($1, $2)
		RETURNING user_id, user_version, user_full_name, user_phone_number`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)

	var userModel UserModel
	err := row.Scan(
		&userModel.UserID,
		&userModel.UserVersion,
		&userModel.UserFullName,
		&userModel.UserPhoneNumber,
	)
	if err != nil {
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
