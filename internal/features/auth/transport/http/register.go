package auth_transport_http

import (
	"net/http"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
	users_transport_http "github.com/Slazzzer/golang-todoapp/internal/features/users/transport/http"
)

// RegisterRequest тело запроса регистрации.
type RegisterRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"Иван Иванов"` // Полное имя (3–100 символов)
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79991234567"` // Телефон в формате +XXXXXXXXXXX (опционально)
}

// RegisterResponse ответ при успешной регистрации.
type RegisterResponse struct {
	Token string                               `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."` // JWT-токен для последующих запросов
	User  users_transport_http.UserDTOResponse `json:"user"`                                    // Созданный пользователь
}

// Register регистрация пользователя и выдача JWT.
//
// @Summary      Регистрация пользователя
// @Description  Публичный эндпоинт. Создаёт пользователя и возвращает JWT-токен (срок жизни задаётся JWT_TTL, по умолчанию 30 дней). Пароль не требуется.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Данные для регистрации"
// @Success      201 {object} RegisterResponse "Пользователь создан, токен выдан"
// @Failure      400 {object} map[string]string "Невалидные данные"
// @Failure      429 {object} map[string]string "Превышен лимит регистраций с одного IP"
// @Router       /auth/register [post]
func (h *AuthHTTPHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userDomain := domain.NewUserUninitialized(request.FullName, request.PhoneNumber)

	userDomain, token, err := h.authService.Register(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to register user")
		return
	}

	responseHandler.JSONResponse(RegisterResponse{
		Token: token,
		User:  users_transport_http.UserDTOFromDomain(userDomain),
	}, http.StatusCreated)
}

// Me возвращает профиль текущего пользователя.
//
// @Summary      Текущий пользователь
// @Description  Возвращает профиль пользователя из JWT-токена. Удобная альтернатива GET /users/{id}, когда id неизвестен.
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} users_transport_http.UserDTOResponse "Профиль текущего пользователя"
// @Failure      401 {object} map[string]string "Токен отсутствует или невалиден"
// @Router       /auth/me [get]
func (h *AuthHTTPHandler) Me(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get authenticated user")
		return
	}

	userDomain, err := h.authService.Me(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get current user")
		return
	}

	responseHandler.JSONResponse(
		users_transport_http.UserDTOFromDomain(userDomain),
		http.StatusOK,
	)
}
