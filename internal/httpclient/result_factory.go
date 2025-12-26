package httpclient

import (
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient/dto"
	dtoDueMonitor "github.com/darkrimson/monitoring_alerting/internal/scheduler/dto"
)

func newErrorResult(m dtoDueMonitor.DueMonitor, err error) dto.Result {
	return dto.Result{
		MonitorID: m.ID,
		Status:    dto.StatusDown,
		Error:     err.Error(),
		CheckedAt: time.Now(),
	}
}

func newSuccessResult(m dtoDueMonitor.DueMonitor, statusCode int, latency time.Duration) dto.Result {
	return dto.Result{
		MonitorID:  m.ID,
		Status:     dto.StatusUp,
		StatusCode: &statusCode,
		LatencyMs:  int(latency.Milliseconds()),
		CheckedAt:  time.Now(),
	}
}

func newFailureResult(m dtoDueMonitor.DueMonitor, statusCode int, latency time.Duration) dto.Result {
	return dto.Result{
		MonitorID:  m.ID,
		Status:     dto.StatusDown,
		StatusCode: &statusCode,
		LatencyMs:  int(latency.Milliseconds()),
		CheckedAt:  time.Now(),
	}
}
