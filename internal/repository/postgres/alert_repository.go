package postgres

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AlertRepository struct {
	pool *pgxpool.Pool
}

func NewAlertRepository(pool *pgxpool.Pool) *AlertRepository {
	return &AlertRepository{pool: pool}
}

func (r *AlertRepository) Create(
	ctx context.Context,
	alert *models.Alert,
) error {

	const query = `
		INSERT INTO alerts (
			id,
			incident_id,
			type,
			channel,
			payload
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		alert.ID,
		alert.IncidentID,
		alert.Type,
		alert.Channel,
		alert.Payload,
	)

	return err
}

func (r *AlertRepository) GetPending(
	ctx context.Context,
) ([]models.Alert, error) {

	const query = `
		SELECT
			id,
			incident_id,
			type,
			channel,
			payload,
			sent_at,
			created_at
		FROM alerts
		WHERE sent_at IS NULL
		ORDER BY created_at
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Alert

	for rows.Next() {
		var a models.Alert
		if err := rows.Scan(
			&a.ID,
			&a.IncidentID,
			&a.Type,
			&a.Channel,
			&a.Payload,
			&a.SentAt,
			&a.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, rows.Err()
}

func (r *AlertRepository) MarkSent(
	ctx context.Context,
	alertID uuid.UUID,
) error {

	const query = `
		UPDATE alerts
		SET sent_at = now()
		WHERE id = $1 AND sent_at IS NULL
	`

	cmd, err := r.pool.Exec(ctx, query, alertID)
	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
