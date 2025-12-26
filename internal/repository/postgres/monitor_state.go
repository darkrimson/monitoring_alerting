package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MonitorStateRepository struct {
	pool *pgxpool.Pool
}

func NewMonitorStateRepository(pool *pgxpool.Pool) *MonitorStateRepository {
	return &MonitorStateRepository{
		pool: pool,
	}
}

const (
	updateMonitorStateQuery = `
		UPDATE monitors
		SET
			last_status = $2,
			last_checked_at = $3,
			failure_streak =
				CASE
					WHEN $2 = 'DOWN' AND $4 = false
						THEN failure_streak + 1
					WHEN $2 = 'UP'
						THEN 0
					ELSE failure_streak
				END
		WHERE id = $1
	`

	incrementFailureStreakQuery = `
		UPDATE monitors
		SET
			failure_streak = failure_streak + 1,
			updated_at = now()
		WHERE id = $1
	`

	resetFailureStreakQuery = `
		UPDATE monitors
		SET
			failure_streak = 0,
			updated_at = now()
		WHERE id = $1
	`
)

func (r *MonitorStateRepository) UpdateStatus(
	ctx context.Context,
	monitorID uuid.UUID,
	status string,
	checkedAt time.Time,
	hasOpenIncident bool,
) error {

	cmd, err := r.pool.Exec(
		ctx,
		updateMonitorStateQuery,
		monitorID,
		status,
		checkedAt,
		hasOpenIncident,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *MonitorStateRepository) IncrementFailureStreak(
	ctx context.Context,
	monitorID uuid.UUID,
) error {

	_, err := r.pool.Exec(ctx, incrementFailureStreakQuery, monitorID)
	return err
}

func (r *MonitorStateRepository) ResetFailureStreak(
	ctx context.Context,
	monitorID uuid.UUID,
) error {

	_, err := r.pool.Exec(ctx, resetFailureStreakQuery, monitorID)
	return err
}
