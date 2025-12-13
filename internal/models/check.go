package models

import (
	"time"

	"github.com/google/uuid"
)

type Check struct {
	ID int64

	MonitorID uuid.UUID

	Timestamp time.Time

	Status string // UP / DOWN

	StatusCode *int
	LatencyMs  *int
	Error      *string
}
