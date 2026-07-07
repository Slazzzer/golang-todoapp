package statistics_postgres_repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Slazzzer/golang-todoapp/internal/core/domain"
)

func (r *StatisticsRepository) GetTasksStatistics(
	ctx context.Context,
	userID *int,
	fromDate *time.Time,
	toDate *time.Time,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`
		SELECT
			task_id,
			task_version,
			task_title,
			task_description,
			task_completed,
			task_created_at,
			task_completed_at,
			author_user_id
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

	queryBuilder.WriteString(" ORDER BY task_id ASC;")

	rows, err := r.pool.Query(ctx, queryBuilder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel
		err := rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task model: %w", err)
		}
		taskModels = append(taskModels, taskModel)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over tasks: %w", err)
	}

	return taskDomainsFromModels(taskModels), nil
}
