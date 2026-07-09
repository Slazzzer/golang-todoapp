package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/Slazzzer/golang-todoapp/internal/core/logger"
	core_http_request "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/response"
)

type GetUsersResponse UsersPageResponse

// GetUsers возвращает список пользователей.
//
// @Summary      Список пользователей
// @Description  Возвращает страницу пользователей с total/limit/offset (по умолчанию limit=50, max=100).
// @Tags         users
// @Produce      json
// @Param        limit query int false "Размер страницы"
// @Param        offset query int false "Смещение"
// @Success      200 {object} GetUsersResponse "Страница пользователей"
// @Failure      400 {object} map[string]string "Невалидные параметры"
// @Router       /users [get]
func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get limit and offset query params",
		)
		return
	}

	page, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")

		return
	}

	response := GetUsersResponse(usersPageFromDomains(page))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		queryParamLimit  = "limit"
		queryParamOffset = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, queryParamLimit)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get limit query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, queryParamOffset)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get offset query param: %w", err)
	}

	return limit, offset, nil
}
