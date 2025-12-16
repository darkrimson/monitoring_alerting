package worker

import (
	"context"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/check"
	"github.com/google/uuid"
)

type ChecksRepository interface {
	Insert(ctx context.Context, result check.Result) error
}

type MonitorStateRepository interface {
	UpdateStatus(
		ctx context.Context,
		monitorID uuid.UUID,
		status string,
		checkedAt time.Time,
	) error
}
