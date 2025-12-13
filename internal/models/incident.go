package models

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID uuid.UUID

	MonitorID uuid.UUID

	StartedAt  time.Time
	ResolvedAt *time.Time

	CurrentStatus string // OPEN / RESOLVED
}
