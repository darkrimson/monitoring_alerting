package dto

import (
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/models"
	"github.com/google/uuid"
)

type MonitorResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	URL             string     `json:"url"`
	IntervalSeconds int        `json:"interval_seconds"`
	TimeoutSeconds  int        `json:"timeout_seconds"`
	ExpectedStatus  int        `json:"expected_status"`
	Enabled         bool       `json:"enabled"`
	LastStatus      *string    `json:"last_status,omitempty"`
	LastCheckedAt   *time.Time `json:"last_checked_at,omitempty"`
}

func MonitorToResponse(m *models.Monitor) MonitorResponse {
	return MonitorResponse{
		ID:              m.ID,
		Name:            m.Name,
		URL:             m.URL,
		IntervalSeconds: m.IntervalSeconds,
		TimeoutSeconds:  m.TimeoutSeconds,
		ExpectedStatus:  m.ExpectedStatus,
		Enabled:         m.Enabled,
		LastStatus:      m.LastStatus,
		LastCheckedAt:   m.LastCheckedAt,
	}
}
