package users_postgres_repository

import "github.com/Slazzzer/golang-todoapp/internal/core/domain"

type UserModel struct {
	UserID          int
	UserVersion     int
	UserFullName    string
	UserPhoneNumber *string
}

func userDomainsFromModels(users []UserModel) []domain.User {
	userDomains := make([]domain.User, len(users))
	for i, user := range users {
		userDomains[i] = domain.NewUser(
			user.UserID,
			user.UserVersion,
			user.UserFullName,
			user.UserPhoneNumber,
		)
	}
	return userDomains
}
