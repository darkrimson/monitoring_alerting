package worker

import (
	"context"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient/dto"
	"github.com/google/uuid"
)

type ChecksRepository interface {
	Insert(ctx context.Context, result dto.Result) (uuid.UUID, error)
}

type MonitorStateRepository interface {
	UpdateStatus(
		ctx context.Context,
		monitorID uuid.UUID,
		status string,
		checkedAt time.Time,
	) error
}
