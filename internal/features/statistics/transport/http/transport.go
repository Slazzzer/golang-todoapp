package statistics_transport_http

import (
	"context"
	"net/http"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_http_server "github.com/Slazzzer/golang-todoapp/internal/core/transport/http/server"
)

type StatisticsHTTPHandler struct {
	statisticsService StatisticsService
}

type StatisticsService interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		fromDate *time.Time,
		toDate *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsHTTPHandler(
	statisticsService StatisticsService,
) *StatisticsHTTPHandler {
	return &StatisticsHTTPHandler{
		statisticsService: statisticsService,
	}
}

func (h *StatisticsHTTPHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		core_http_server.NewRoute(http.MethodGet, "/statistics", h.GetStatistics),
	}
}
