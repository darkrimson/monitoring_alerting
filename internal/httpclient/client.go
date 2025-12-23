package httpclient

import (
	"context"
	"net/http"
	"time"

	dto2 "github.com/darkrimson/monitoring_alerting/internal/httpclient/dto"
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

func (c *Client) Check(ctx context.Context, m dto.DueMonitor) dto2.Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return dto2.Result{
			MonitorID: m.ID,
			Status:    dto2.StatusDown,
			Error:     err.Error(),
			CheckedAt: time.Now(),
		}
	}

	client := *c.httpClient
	client.Timeout = time.Duration(m.TimeoutSeconds) * time.Second

	resp, err := client.Do(req)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return dto2.Result{
			MonitorID: m.ID,
			Status:    dto2.StatusDown,
			LatencyMs: int(latency),
			Error:     err.Error(),
			CheckedAt: time.Now(),
		}
	}
	defer resp.Body.Close()

	status := dto2.StatusDown
	if resp.StatusCode == m.ExpectedStatusCode {
		status = dto2.StatusUp
	}

	code := resp.StatusCode

	return dto2.Result{
		MonitorID:  m.ID,
		Status:     status,
		StatusCode: &code,
		LatencyMs:  int(latency),
		CheckedAt:  time.Now(),
	}
}
