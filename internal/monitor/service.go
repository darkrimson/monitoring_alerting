package monitor

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewMonitorService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, m *models.Monitor) error {
	return s.repo.Create(ctx, m)
}

func (s *Service) Update(ctx context.Context, m *models.Monitor) error {
	return s.repo.Update(ctx, m)
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*models.Monitor, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]models.Monitor, error) {
	return s.repo.List(ctx)
}
