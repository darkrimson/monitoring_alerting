package models

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID uuid.UUID

	Name string
	URL  string

	IntervalSeconds int
	TimeoutSeconds  int
	ExpectedStatus  int

	Enabled bool

	LastStatus    *string
	LastCheckedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
