package dto

import "github.com/google/uuid"

type DueMonitor struct {
	ID                 uuid.UUID
	URL                string
	TimeoutSeconds     int
	ExpectedStatusCode int
}
