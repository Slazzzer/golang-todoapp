package pgx

import (
	"context"
	"fmt"
	"time"

	core_postgres_pool "github.com/Slazzzer/golang-todoapp/internal/core/repository/postgres/pool"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	pool      *pgxpool.Pool
	opTimeout time.Duration
}

func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*ConnectionPool, error) {
	connectionString := buildConnectionString(config)

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

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapError(err)
	}

	return &pgxRows{rows: rows}, nil
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	return pgxRow{row: p.pool.QueryRow(ctx, sql, args...)}
}

func (p *ConnectionPool) Exec(ctx context.Context, sql string, args ...any) (core_postgres_pool.ExecResult, error) {
	tag, err := p.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapError(err)
	}

	return pgxExecResult{tag: tag}, nil
}

func (p *ConnectionPool) Ping(ctx context.Context) error {
	return mapError(p.pool.Ping(ctx))
}

func (p *ConnectionPool) Close() {
	p.pool.Close()
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.opTimeout
}
