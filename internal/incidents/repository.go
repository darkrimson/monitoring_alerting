package incidents

import (
	"context"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetOpenByMonitor(ctx context.Context, monitorID uuid.UUID) (*models.Incident, error)

	CreateIncident(
		ctx context.Context,
		incident *models.Incident,
	) error

	UpdateFailure(
		ctx context.Context,
		incidentID uuid.UUID,
		lastCheckID uuid.UUID,
	) error

	ResolveIncident(
		ctx context.Context,
		incidentID uuid.UUID,
		lastCheckID uuid.UUID,
		resolvedAt time.Time,
	) error
}
