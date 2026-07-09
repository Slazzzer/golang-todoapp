package tasks_postgres_repository

import (
	"context"
	"fmt"
)

func (r *TasksRepository) CountTasks(
	ctx context.Context,
	userID *int,
) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `SELECT COUNT(*) FROM todoapp.tasks %s;`
	args := []any{}

	if userID != nil {
		query = fmt.Sprintf(query, "WHERE author_user_id = $1")
		args = append(args, *userID)
	} else {
		query = fmt.Sprintf(query, "")
	}

	var total int
	if err := r.pool.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, fmt.Errorf("count tasks: %w", err)
	}

	return total, nil
}
