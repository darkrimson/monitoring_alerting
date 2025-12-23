package alerts

import (
	"encoding/json"

	"github.com/darkrimson/monitoring_alerting/internal/models"
)

func BuildIncidentPayload(incident *models.Incident) []byte {
	payload, _ := json.Marshal(map[string]any{
		"incident_id": incident.ID,
		"monitor_id":  incident.MonitorID,
		"status":      incident.Status,
		"started_at":  incident.StartedAt,
		"resolved_at": incident.ResolvedAt,
		"failures":    incident.FailureCount,
	})
	return payload
}
