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
	GetTasksStatistics(
		ctx context.Context,
		userID *int,
		fromDate *time.Time,
		toDate *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(
	statisticsRepository StatisticsRepository,
) *StatisticsService {
	return &StatisticsService{
		statisticsRepository: statisticsRepository,
	}
}
