package alerts

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, alert *models.Alert) error
	GetPending(ctx context.Context) ([]models.Alert, error)
	MarkSent(ctx context.Context, alertID uuid.UUID) error
}
