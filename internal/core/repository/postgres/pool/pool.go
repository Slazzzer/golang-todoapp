package core_postgres_pool

import (
	"context"
	"errors"
	"time"
)

// ErrNoRows — строка не найдена (аналог sql.ErrNoRows).
var ErrNoRows = errors.New("no rows in result set")

func IsErrNoRows(err error) bool {
	return errors.Is(err, ErrNoRows)
}

// Row — одна строка результата запроса.
type Row interface {
	Scan(dest ...any) error
}

// Rows — итератор по строкам результата запроса.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
	Err() error
}

// ExecResult — результат INSERT/UPDATE/DELETE.
type ExecResult interface {
	RowsAffected() int64
}

// Pool — абстракция над пулом соединений к PostgreSQL.
// Не зависит от конкретной драйверной библиотеки.
type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, args ...any) (ExecResult, error)
	Close()

	OpTimeout() time.Duration
}
