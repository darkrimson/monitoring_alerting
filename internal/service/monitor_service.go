package service

import (
	"context"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
)

type MonitorService struct {
	repo MonitorRepository
}

func NewMonitorService(repo MonitorRepository) *MonitorService {
	return &MonitorService{repo: repo}
}

func (s *MonitorService) Create(ctx context.Context, m *models.Monitor) error {
	return s.repo.Create(ctx, m)
}

func (s *MonitorService) Update(ctx context.Context, m *models.Monitor) error {
	return s.repo.Update(ctx, m)
}

func (s *MonitorService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *MonitorService) GetByID(ctx context.Context, id uuid.UUID) (*models.Monitor, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *MonitorService) List(ctx context.Context) ([]models.Monitor, error) {
	return s.repo.List(ctx)
}
