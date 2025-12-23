package models

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID              uuid.UUID
	Name            string
	URL             string
	LastStatus      *string
	IntervalSeconds int
	TimeoutSeconds  int
	ExpectedStatus  int
	FailureStreak   int
	Enabled         bool
	LastCheckedAt   *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
