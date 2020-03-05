package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var httpClient = http.Client{
	Timeout: 15 * time.Second,
}

type Slack struct {
	name       string
	webhookURL string
}

func NewSlack(name, webhook string) *Slack {
	return &Slack{
		name:       name,
		webhookURL: webhook,
	}
}

func (s *Slack) SendAlert(alert *Alert, data string) error {
	request := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": fmt.Sprintf("[ALERT] :exclamation: %s\n%s", alert.Name, alert.Message),
					"emoji": true,
				},
			},
			map[string]interface{}{
				"type": "context",
				"elements": []interface{}{
					map[string]interface{}{
						"type": "mrkdwn",
						"text": fmt.Sprintf("```%s```", data),
					},
				},
			},
		},
	}

	return s.sendMessage(request)

}

func (s *Slack) ResolveAlert(alert *Alert, message string) error {
	request := map[string]interface{}{
		"blocks": []interface{}{
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": fmt.Sprintf("[RESOLVED] :tick: %s\n%s", alert.Name, message),
					"emoji": true,
				},
			},
		},
	}

	return s.sendMessage(request)
}

func (s *Slack) Name() string {
	return s.name
}

func (s *Slack) sendMessage(r interface{}) error {
	req, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("couldn't prepare slack webhook request")
	}
	resp, err := httpClient.Post(s.webhookURL, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return fmt.Errorf("couldn't send slack message: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected http status from slack: %d", resp.StatusCode)
	}
	return nil
}
