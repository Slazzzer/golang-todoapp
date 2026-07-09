package users_transport_http

import (
	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_pagination "github.com/Slazzzer/golang-todoapp/internal/core/pagination"
)

type UsersPageResponse = core_pagination.Page[UserDTOResponse]

func usersPageFromDomains(
	page core_pagination.Page[domain.User],
) UsersPageResponse {
	items := make([]UserDTOResponse, len(page.Items))
	for i, user := range page.Items {
		items[i] = userDTOFromDomain(user)
	}

	return core_pagination.NewPage(items, page.Total, page.Limit, page.Offset)
}
