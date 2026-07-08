package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

const completedDurationFilter = `
	task_completed
	AND task_completed_at IS NOT NULL
	AND task_completed_at >= task_created_at
`

func (r *StatisticsRepository) GetStatistics(
	ctx context.Context,
	userID *int,
	fromDate *time.Time,
	toDate *time.Time,
) (domain.Statistics, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT
			COUNT(*)::int,
			COUNT(*) FILTER (WHERE task_completed)::int,
			COUNT(*) FILTER (WHERE ` + completedDurationFilter + `)::int,
			AVG(EXTRACT(EPOCH FROM (task_completed_at - task_created_at)))
				FILTER (WHERE ` + completedDurationFilter + `)
		FROM todoapp.tasks
	`)

	args := make([]any, 0, 3)
	conditions := make([]string, 0, 3)

	if userID != nil {
		conditions = append(conditions, "author_user_id = $"+strconv.Itoa(len(args)+1))
		args = append(args, *userID)
	}

	if fromDate != nil {
		conditions = append(conditions, "task_created_at >= $"+strconv.Itoa(len(args)+1))
		args = append(args, *fromDate)
	}

	if toDate != nil {
		conditions = append(conditions, "task_created_at < $"+strconv.Itoa(len(args)+1))
		args = append(args, *toDate)
	}

	if len(conditions) > 0 {
		queryBuilder.WriteString(" WHERE ")
		queryBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	queryBuilder.WriteString(";")

	var (
		tasksCreated        int
		tasksCompleted      int
		tasksWithDuration   int
		avgSeconds          *float64
	)

	row := r.pool.QueryRow(ctx, queryBuilder.String(), args...)
	if err := row.Scan(&tasksCreated, &tasksCompleted, &tasksWithDuration, &avgSeconds); err != nil {
		return domain.Statistics{}, fmt.Errorf("scan statistics: %w", err)
	}

	if tasksCreated == 0 {
		return domain.NewStatistics(0, 0, nil, nil), nil
	}

	completedRate := float64(tasksCompleted) / float64(tasksCreated) * 100

	var avgCompletionTime *time.Duration
	if tasksWithDuration > 0 && avgSeconds != nil {
		avg := time.Duration(int64(*avgSeconds * float64(time.Second)))
		avgCompletionTime = &avg
	}

	return domain.NewStatistics(
		tasksCreated,
		tasksCompleted,
		&completedRate,
		avgCompletionTime,
	), nil
}
