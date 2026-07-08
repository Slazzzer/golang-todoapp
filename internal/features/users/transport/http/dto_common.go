package users_transport_http

import "github.com/Slazzzer/golang-todoapp/internal/core/domain"

// UserDTOResponse данные пользователя в ответе API.
type UserDTOResponse struct {
	ID          int     `json:"id" example:"1"`                         // Идентификатор пользователя
	Version     int     `json:"version" example:"1"`                  // Версия записи (для оптимистичной блокировки)
	FullName    string  `json:"full_name" example:"Иван Иванов"`      // Полное имя
	PhoneNumber *string `json:"phone_number" example:"+79991234567"` // Телефон или null
}

func UserDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOFromDomain(user)
}
