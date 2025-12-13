package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AlertTarget struct {
	ID uuid.UUID

	Type string // telegram

	Payload json.RawMessage

	CreatedAt time.Time
}
