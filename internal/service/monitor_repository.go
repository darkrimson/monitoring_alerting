package service

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
)

type MonitorRepository interface {
	Create(ctx context.Context, m *models.Monitor) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Monitor, error)
	Update(ctx context.Context, m *models.Monitor) error
	List(ctx context.Context) ([]models.Monitor, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
