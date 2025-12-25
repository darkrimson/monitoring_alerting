package postgres

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient/dto"
	"github.com/google/uuid"
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

const (
	insertChecksQuery = `
		INSERT INTO checks (
			monitor_id,
			ts,
			status,
			status_code,
			latency_ms,
			error
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
)

func (r *ChecksRepository) Insert(
	ctx context.Context,
	result dto.Result,
) (uuid.UUID, error) {

	var checkID uuid.UUID

	err := r.pool.QueryRow(
		ctx,
		insertChecksQuery,
		result.MonitorID,
		result.CheckedAt,
		result.Status,
		result.StatusCode,
		result.LatencyMs,
		nullIfEmpty(result.Error),
	).Scan(&checkID)

	return checkID, err
}

func nullIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
