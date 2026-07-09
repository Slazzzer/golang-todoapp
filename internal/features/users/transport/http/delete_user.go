package users_transport_http

import (
	"net/http"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
)

// DeleteUser удаляет пользователя.
//
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя по id. Если есть связанные задачи, может вернуться 409.
// @Tags         users
// @Param        id path int true "ID пользователя"
// @Success      204 "Пользователь удалён"
// @Failure      404 {object} map[string]string "Пользователь не найден"
// @Failure      409 {object} map[string]string "Конфликт (например, есть связанные данные)"
// @Router       /users/{id} [delete]
func (h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
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

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
		)
		return
	}

	responseHandler.NoContentResponse()
}
