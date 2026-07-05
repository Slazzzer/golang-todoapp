package core_postgres_pool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type pgxRow struct {
	row pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	if err := r.row.Scan(dest...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNoRows
		}

		return err
	}

	return nil
}

type pgxRows struct {
	rows pgx.Rows
}

func (r *pgxRows) Next() bool {
	return r.rows.Next()
}

func (r *pgxRows) Scan(dest ...any) error {
	return r.rows.Scan(dest...)
}

func (r *pgxRows) Close() {
	r.rows.Close()
}

func (r *pgxRows) Err() error {
	return r.rows.Err()
}

type pgxExecResult struct {
	tag pgconn.CommandTag
}

func (r pgxExecResult) RowsAffected() int64 {
	return r.tag.RowsAffected()
}

type ConnectionPool struct {
	pool      *pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxconfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxconfig)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping pgx pool: %w", err)
	}

	return &ConnectionPool{
		pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return &pgxRows{rows: rows}, nil
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return pgxRow{row: p.pool.QueryRow(ctx, sql, args...)}
}

func (p *ConnectionPool) Exec(ctx context.Context, sql string, args ...any) (ExecResult, error) {
	tag, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	return pgxExecResult{tag: tag}, nil
}

func (p *ConnectionPool) Close() {
	p.pool.Close()
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
