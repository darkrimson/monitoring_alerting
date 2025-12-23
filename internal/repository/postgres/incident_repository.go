package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IncidentRepository struct {
	pool *pgxpool.Pool
}

func NewIncidentRepository(pool *pgxpool.Pool) *IncidentRepository {
	return &IncidentRepository{pool: pool}
}

func (r *IncidentRepository) GetOpenByMonitor(
	ctx context.Context,
	monitorID uuid.UUID,
) (*models.Incident, error) {

	const query = `
		SELECT
			id,
			monitor_id,
			status,
			started_at,
			resolved_at,
			failure_count,
			last_check_id,
			created_at,
			updated_at
		FROM incidents
		WHERE monitor_id = $1 AND status = 'OPEN'
		LIMIT 1;
	`

	var inc models.Incident

	err := r.pool.QueryRow(ctx, query, monitorID).Scan(
		&inc.ID,
		&inc.MonitorID,
		&inc.Status,
		&inc.StartedAt,
		&inc.ResolvedAt,
		&inc.FailureCount,
		&inc.LastCheckID,
		&inc.CreatedAt,
		&inc.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &inc, nil
}

func (r *IncidentRepository) CreateIncident(
	ctx context.Context,
	incident *models.Incident,
) error {

	const query = `
		INSERT INTO incidents (
			id,
			monitor_id,
			status,
			started_at,
			failure_count,
			last_check_id
		)
		VALUES ($1, $2, 'OPEN', $3, $4, $5);
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		incident.ID,
		incident.MonitorID,
		incident.StartedAt,
		incident.FailureCount,
		incident.LastCheckID,
	)

	return err
}

func (r *IncidentRepository) UpdateFailure(
	ctx context.Context,
	incidentID uuid.UUID,
	lastCheckID uuid.UUID,
) error {

	const query = `
		UPDATE incidents
		SET
			failure_count = failure_count + 1,
			last_check_id = $2,
			updated_at = now()
		WHERE id = $1 AND status = 'OPEN';
	`

	cmd, err := r.pool.Exec(ctx, query, incidentID, lastCheckID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *IncidentRepository) ResolveIncident(
	ctx context.Context,
	incidentID uuid.UUID,
	lastCheckID uuid.UUID,
	resolvedAt time.Time,
) error {

	const query = `
		UPDATE incidents
		SET
			status = 'RESOLVED',
			resolved_at = $2,
			last_check_id = $3,
			updated_at = now()
		WHERE id = $1 AND status = 'OPEN';
	`

	cmd, err := r.pool.Exec(
		ctx,
		query,
		incidentID,
		resolvedAt,
		lastCheckID,
	)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
