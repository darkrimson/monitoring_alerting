package postgres

import (
	"context"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/scheduler/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchedulerRepository struct {
	pool *pgxpool.Pool
}

func NewSchedulerRepository(pool *pgxpool.Pool) *SchedulerRepository {
	return &SchedulerRepository{pool: pool}
}

func (r *SchedulerRepository) SelectDueMonitors(ctx context.Context, now time.Time) ([]dto.DueMonitor, error) {
	const query = `
		SELECT
			id,
			url,
			timeout_seconds,
			expected_status
		FROM monitors
		WHERE
			enabled = true
			AND (
				last_checked_at IS NULL
				OR last_checked_at + (interval_seconds * interval '1 second') <= $1
			);
	`

	rows, err := r.pool.Query(ctx, query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []dto.DueMonitor

	for rows.Next() {
		var m dto.DueMonitor
		if err := rows.Scan(
			&m.ID,
			&m.URL,
			&m.TimeoutSeconds,
			&m.ExpectedStatusCode,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}

	return result, rows.Err()
}
