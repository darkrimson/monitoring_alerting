package postgres

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/darkrimson/monitoring_alerting/internal/service"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type monitorRepo struct {
	pool *pgxpool.Pool
}

func NewMonitorRepository(pool *pgxpool.Pool) service.MonitorRepository {
	return &monitorRepo{pool: pool}
}

func (r *monitorRepo) Create(ctx context.Context, m *models.Monitor) error {
	const query = `
		INSERT INTO monitors (
		    name, url, interval_seconds, timeout_seconds, expected_status, enabled
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
		`
	return r.pool.QueryRow(
		ctx, query,
		m.Name,
		m.URL,
		m.IntervalSeconds,
		m.TimeoutSeconds,
		m.ExpectedStatus,
		m.Enabled,
	).Scan(&m.ID, &m.CreatedAt, &m.UpdatedAt)
}

func (r *monitorRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Monitor, error) {
	const query = `
		SELECT
			id, name, url,
			interval_seconds, timeout_seconds, 
			expected_status, enabled,
			last_status, last_checked_at,
			created_at, updated_at
		FROM monitors
		WHERE id = $1
	`

	var m models.Monitor

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&m.ID,
		&m.Name,
		&m.URL,
		&m.IntervalSeconds,
		&m.TimeoutSeconds,
		&m.ExpectedStatus,
		&m.Enabled,
		&m.LastStatus,
		&m.LastCheckedAt,
		&m.CreatedAt,
		&m.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *monitorRepo) List(ctx context.Context) ([]models.Monitor, error) {
	const query = `
		SELECT
			id, name, url,
			interval_seconds, timeout_seconds,
			expected_status, enabled,
			last_status, last_checked_at,
			created_at, updated_at
		FROM monitors
		ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Monitor

	for rows.Next() {
		var m models.Monitor
		if err := rows.Scan(
			&m.ID,
			&m.Name,
			&m.URL,
			&m.IntervalSeconds,
			&m.TimeoutSeconds,
			&m.ExpectedStatus,
			&m.Enabled,
			&m.LastStatus,
			&m.LastCheckedAt,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result, rows.Err()
}

func (r *monitorRepo) Update(ctx context.Context, m *models.Monitor) error {
	const query = `
		UPDATE monitors
		SET
			name = $2,
			url = $3,
			interval_seconds = $4,
			timeout_seconds = $5,
			expected_status = $6,
			enabled = $7,
			updated_at = now()
		WHERE id = $1 
		`

	cmd, err := r.pool.Exec(ctx, query,
		m.ID,
		m.Name,
		m.URL,
		m.IntervalSeconds,
		m.TimeoutSeconds,
		m.ExpectedStatus,
		m.Enabled,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *monitorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = `
		DELETE FROM monitors
		WHERE id = $1
	`

	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
