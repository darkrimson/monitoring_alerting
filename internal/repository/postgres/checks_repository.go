package postgres

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/check"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChecksRepository struct {
	pool *pgxpool.Pool
}

func NewChecksRepository(pool *pgxpool.Pool) *ChecksRepository {
	return &ChecksRepository{
		pool: pool,
	}
}

func (r *ChecksRepository) Insert(ctx context.Context, result check.Result) error {
	const query = `
		INSERT INTO checks (
			monitor_id,
			ts,
			status,
			status_code,
			latency_ms,
			error
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(
		ctx,
		query,
		result.MonitorID,
		result.CheckedAt,
		result.Status,
		result.StatusCode,
		result.LatencyMs,
		nullIfEmpty(result.Error),
	)

	return err
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
