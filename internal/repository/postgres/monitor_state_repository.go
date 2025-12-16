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

func (r *MonitorStateRepository) UpdateStatus(
	ctx context.Context,
	monitorID uuid.UUID,
	status string,
	checkedAt time.Time,
) error {

	const query = `
		UPDATE monitors
		SET
			last_status = $2,
			last_checked_at = $3
		WHERE id = $1
	`

	cmd, err := r.pool.Exec(
		ctx,
		query,
		monitorID,
		status,
		checkedAt,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
