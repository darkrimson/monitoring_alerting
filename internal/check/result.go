package check

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

type Result struct {
	MonitorID  uuid.UUID
	Status     Status
	StatusCode *int
	LatencyMs  int
	Error      string
	CheckedAt  time.Time
}
