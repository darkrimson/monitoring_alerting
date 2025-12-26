package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/darkrimson/monitoring_alerting/internal/config"
	"github.com/darkrimson/monitoring_alerting/internal/models"
)

type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramNotifier(cfg config.TelegramConfig) *TelegramNotifier {
	return &TelegramNotifier{
		token:  cfg.Token,
		chatID: cfg.ChatID,
		client: &http.Client{},
	}
}

const (
	incidentOpened   = "INCIDENT_OPENED"
	incidentResolved = "INCIDENT_RESOLVED"
)

func (t *TelegramNotifier) Send(
	ctx context.Context,
	alert models.Alert,
) error {

	text := t.formatMessage(alert)

	body, _ := json.Marshal(map[string]string{
		"chat_id": t.chatID,
		"text":    text,
	})

	url := fmt.Sprintf(
		"https://api.telegram.org/bot%s/sendMessage",
		t.token,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("telegram send failed: %s", resp.Status)
	}

	return nil
}

func (t *TelegramNotifier) formatMessage(
	alert models.Alert,
) string {

	switch alert.Type {

	case incidentOpened:
		return "🚨 INCIDENT OPENED\n\n" + string(alert.Payload)

	case incidentResolved:
		return "✅ INCIDENT RESOLVED\n\n" + string(alert.Payload)

	default:
		return "ℹ️ ALERT\n\n" + string(alert.Payload)
	}
}
