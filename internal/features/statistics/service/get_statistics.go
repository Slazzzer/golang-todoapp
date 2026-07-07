package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_errors "github.com/Slazzzer/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	fromDate *time.Time,
	toDate *time.Time,
) (domain.Statistics, error) {

	if fromDate != nil && toDate != nil {
		if toDate.Before(*fromDate) || toDate.Equal(*fromDate) {
			return domain.Statistics{}, fmt.Errorf(
				"`toDate` must be after `fromDate`: %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	statistics, err := s.statisticsRepository.GetStatistics(ctx, userID, fromDate, toDate)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"failed to get statistics from repository: %w",
			err,
		)
	}

	return statistics, nil
}
