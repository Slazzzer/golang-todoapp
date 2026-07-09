package tasks_transport_http

import (
	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_pagination "github.com/Slazzzer/golang-todoapp/internal/core/pagination"
)

type TasksPageResponse = core_pagination.Page[TaskDTOResponse]

func tasksPageFromDomains(
	page core_pagination.Page[domain.Task],
) TasksPageResponse {
	items := make([]TaskDTOResponse, len(page.Items))
	for i, task := range page.Items {
		items[i] = taskDTOFromDomain(task)
	}

	return core_pagination.NewPage(items, page.Total, page.Limit, page.Offset)
}
