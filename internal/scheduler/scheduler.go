package scheduler

import (
	"context"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/scheduler/dto"
)

type Repository interface {
	SelectDueMonitors(ctx context.Context, now time.Time) ([]dto.DueMonitor, error)
}

type Scheduler struct {
	repo Repository
}

func NewScheduler(repo Repository) *Scheduler {
	return &Scheduler{repo: repo}
}

func (s *Scheduler) DueMonitors(ctx context.Context, now time.Time) ([]dto.DueMonitor, error) {
	return s.repo.SelectDueMonitors(ctx, now)
}
