package models

import "github.com/google/uuid"

type MonitorAlert struct {
	MonitorID     uuid.UUID
	AlertTargetID uuid.UUID
}
