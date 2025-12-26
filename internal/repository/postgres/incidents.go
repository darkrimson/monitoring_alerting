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

const (
	selectOpenIncidentByMonitorQuery = `
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

	insertIncidentQuery = `
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
	updateIncidentFailureQuery = `
		UPDATE incidents
		SET
			failure_count = failure_count + 1,
			last_check_id = $2,
			updated_at = now()
		WHERE id = $1 AND status = 'OPEN';
	`

	resolveIncidentQuery = `
		UPDATE incidents
		SET
			status = 'RESOLVED',
			resolved_at = $2,
			last_check_id = $3,
			updated_at = now()
		WHERE id = $1
		RETURNING
			id,
			monitor_id,
			status,
			started_at,
			resolved_at,
			failure_count
	`
)

func (r *IncidentRepository) GetOpenByMonitor(
	ctx context.Context,
	monitorID uuid.UUID,
) (*models.Incident, error) {

	var inc models.Incident

	err := r.pool.QueryRow(ctx, selectOpenIncidentByMonitorQuery, monitorID).Scan(
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

	_, err := r.pool.Exec(
		ctx,
		insertIncidentQuery,
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

	cmd, err := r.pool.Exec(ctx, updateIncidentFailureQuery, incidentID, lastCheckID)
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
) (*models.Incident, error) {

	var incident models.Incident

	err := r.pool.QueryRow(
		ctx,
		resolveIncidentQuery,
		incidentID,
		resolvedAt,
		lastCheckID,
	).Scan(
		&incident.ID,
		&incident.MonitorID,
		&incident.Status,
		&incident.StartedAt,
		&incident.ResolvedAt,
		&incident.FailureCount,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}

	return &incident, nil
}
