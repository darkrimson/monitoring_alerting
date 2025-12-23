package alerts

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/darkrimson/monitoring_alerting/internal/models"
)

type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramNotifier() *TelegramNotifier {
	return &TelegramNotifier{
		token:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		chatID: os.Getenv("TELEGRAM_CHAT_ID"),
		client: &http.Client{},
	}
}

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

	if resp.StatusCode >= 300 {
		return fmt.Errorf("telegram send failed: %s", resp.Status)
	}

	return nil
}

func (t *TelegramNotifier) formatMessage(
	alert models.Alert,
) string {

	switch alert.Type {

	case "INCIDENT_OPENED":
		return "🚨 INCIDENT OPENED\n\n" + string(alert.Payload)

	case "INCIDENT_RESOLVED":
		return "✅ INCIDENT RESOLVED\n\n" + string(alert.Payload)

	default:
		return "ℹ️ ALERT\n\n" + string(alert.Payload)
	}
}
