package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
	core_auth "github.com/Slazzzer/golang-todoapp/internal/core/auth"
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

	requesterID, err := core_auth.UserIDFromContext(ctx)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf("get requester id: %w", err)
	}

	if userID != nil && *userID != requesterID {
		return domain.Statistics{}, fmt.Errorf(
			"statistics for another user are forbidden: %w",
			core_errors.ErrForbidden,
		)
	}

	statistics, err := s.statisticsRepository.GetStatistics(ctx, &requesterID, fromDate, toDate)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"failed to get statistics from repository: %w",
			err,
		)
	}

	return statistics, nil
}
