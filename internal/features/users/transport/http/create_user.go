package users_transport_http

import (
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

// CreateUserRequest тело запроса создания пользователя.
type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"Иван Иванов"`                  // Полное имя (3–100 символов)
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79991234567"` // Телефон в формате +7... (опционально)
}

type CreateUserResponse UserDTOResponse

// CreateUser создаёт нового пользователя.
//
// @Summary      Создать пользователя
// @Description  Регистрирует нового пользователя. Телефон опционален.
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "Данные нового пользователя"
// @Success      201 {object} CreateUserResponse "Пользователь создан"
// @Failure      400 {object} map[string]string "Невалидные данные"
// @Router       /users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
