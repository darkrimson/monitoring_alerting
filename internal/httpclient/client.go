package httpclient

import (
	"context"
	"net/http"
	"time"

	"github.com/darkrimson/monitoring_alerting/internal/httpclient/dto"
	dtoDueMonitor "github.com/darkrimson/monitoring_alerting/internal/scheduler/dto"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) Check(ctx context.Context, m dtoDueMonitor.DueMonitor) dto.Result {
	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.URL, nil)
	if err != nil {
		return newErrorResult(m, err)
	}

	client := *c.httpClient
	client.Timeout = time.Duration(m.TimeoutSeconds) * time.Second

	resp, err := client.Do(req)
	if err != nil {
		return newErrorResult(m, err)
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	if resp.StatusCode != m.ExpectedStatusCode {
		return newFailureResult(m, resp.StatusCode, latency)
	}

	return newSuccessResult(m, resp.StatusCode, latency)
}
