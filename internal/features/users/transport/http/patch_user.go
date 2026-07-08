package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/types"
)

// PatchUserRequest тело частичного обновления профиля.
type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`    // Новое полное имя (опционально)
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"` // Новый телефон (опционально, null — удалить)
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))

		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set && r.PhoneNumber.Value != nil {
		phoneNumberLen := len([]rune(*r.PhoneNumber.Value))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
		}

		if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
			return fmt.Errorf("`PhoneNumber` must start with '+' symbol")
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser обновляет профиль текущего пользователя.
//
// @Summary      Изменить пользователя
// @Description  Частичное обновление профиля. Доступ только к своему id.
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "ID пользователя"
// @Param        request body PatchUserRequest true "Поля для обновления"
// @Success      200 {object} PatchUserResponse "Обновлённый профиль"
// @Failure      400 {object} map[string]string "Невалидные данные"
// @Failure      401 {object} map[string]string "Не авторизован"
// @Failure      403 {object} map[string]string "Доступ к чужому профилю запрещён"
// @Router       /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user ID from path value",
		)
		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate request",
		)
		return
	}

	userPatch := domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
