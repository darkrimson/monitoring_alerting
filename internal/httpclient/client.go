package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/check"
	"github.com/darkrimson/monitoring_alerting/internal/scheduler/dto"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) Check(ctx context.Context, m dto.DueMonitor) check.Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return check.Result{
			MonitorID: m.ID,
			Status:    check.StatusDown,
			Error:     err.Error(),
			CheckedAt: time.Now(),
		}
	}

	client := *c.httpClient
	client.Timeout = time.Duration(m.TimeoutSeconds) * time.Second

	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return check.Result{
			MonitorID: m.ID,
			Status:    check.StatusDown,
			LatencyMs: int(latency),
			Error:     err.Error(),
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	status := check.StatusDown
	if resp.StatusCode == m.ExpectedStatusCode {
		status = check.StatusUp
	}

	code := resp.StatusCode

	return check.Result{
		MonitorID:  m.ID,
		Status:     status,
		StatusCode: &code,
		LatencyMs:  int(latency),
		CheckedAt:  time.Now(),
	}
}
