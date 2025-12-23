package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Alert struct {
	ID         uuid.UUID
	IncidentID uuid.UUID
	Type       string // INCIDENT_OPENED / INCIDENT_RESOLVED
	Channel    string // TELEGRAM
	Payload    json.RawMessage
	SentAt     *time.Time
	CreatedAt  time.Time
}
