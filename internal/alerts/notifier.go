package alerts

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
)

type Notifier interface {
	Send(ctx context.Context, alert models.Alert) error
}
