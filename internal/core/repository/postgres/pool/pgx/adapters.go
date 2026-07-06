package pgx

import (
	"errors"
	"fmt"

	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgForeignKeyViolation = "23503"

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return core_postgres_pool.ErrNoRows
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) && pgErr.Code == pgForeignKeyViolation {
		return fmt.Errorf(
			"%v: %w",
			err,
			core_postgres_pool.ErrViolatesForeignKey,
		)
	}

	return err
}

type pgxRow struct {
	row pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.row.Scan(dest...)
	if err != nil {
		return mapError(err)
	}
	return nil
}

type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Next() bool {
	next := r.rows.Next()
	if !next {
		return false
	}
	return next
}

func (r *pgxRows) Scan(dest ...any) error {
	err := r.rows.Scan(dest...)
	if err != nil {
		return mapError(err)
	}
	return nil
}

func (r *pgxRows) Close() {
	r.rows.Close()
}

func (r *pgxRows) Err() error {
	err := r.rows.Err()
	if err != nil {
		return mapError(err)
	}
	return nil
}

type pgxExecResult struct {
	tag pgconn.CommandTag
}

func (r pgxExecResult) RowsAffected() int64 {
	return r.tag.RowsAffected()
}
