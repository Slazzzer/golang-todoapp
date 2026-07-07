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

	tasks, err := s.statisticsRepository.GetTasksStatistics(ctx, userID, fromDate, toDate)

	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"failed to get tasks statistics from repository: %w",
			err,
		)
	}

	statistics := calcStatistics(tasks)
	return statistics, nil
}

func calcStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.NewStatistics(0, 0, nil, nil)
	}

	tasksCreated := len(tasks)
	var totalCompletionDuration time.Duration
	tasksCompleted := 0
	for _, task := range tasks {
		if task.Completed {
			tasksCompleted++
		}

		completionDuration := task.GetCompletionDuration()
		if completionDuration != nil {
			totalCompletionDuration += *completionDuration
		}
	}

	taskCompletedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var tasksAverageCompletionTime *time.Duration
	if tasksCompleted > 0 && totalCompletionDuration != 0 {
		avg := totalCompletionDuration / time.Duration(tasksCompleted)
		tasksAverageCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&taskCompletedRate,
		tasksAverageCompletionTime,
	)

}
