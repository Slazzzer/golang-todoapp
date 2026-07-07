package statistics_service

import (
	"context"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	statisticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetStatistics(
		ctx context.Context,
		userID *int,
		fromDate *time.Time,
		toDate *time.Time,
	) (domain.Statistics, error)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
