package models

import (
	"time"

	"github.com/google/uuid"
)

type Incident struct {
	ID           uuid.UUID
	MonitorID    uuid.UUID
	Status       string // OPEN / RESOLVED
	StartedAt    time.Time
	ResolvedAt   *time.Time
	FailureCount int
	LastCheckID  *uuid.UUID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
