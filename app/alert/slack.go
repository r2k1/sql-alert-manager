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

type SlackWebhookRequest struct {
	Text        string            `json:"text"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Color      string `json:"color"`
	AuthorName string `json:"author_name"`
	AuthorLink string `json:"author_link"`
	AuthorIcon string `json:"author_icon"`
	Title      string `json:"title"`
	TitleLink  string `json:"title_link"`
	Text       string `json:"text"`
	Pretext    string `json:"pretext"`
	Footer     string `json:"footer"`
	ImageURL   string `json:"image_url"`
	ThumbURL   string `json:"thumb_url"`
	FooterIcon string `json:"footer_icon"`
	Timestamp  int64  `json:"ts"`
	Fallback   string `json:"fallback"`
}

func NewSlack(name, webhook string) *Slack {
	return &Slack{
		name:       name,
		webhookURL: webhook,
	}
}

func (s *Slack) SendAlert(alert *Alert, data string) error {
	return s.sendMessage(SlackWebhookRequest{
		Attachments: []SlackAttachment{
			{
				Color: "#D5212E",
				Text:  fmt.Sprintf("[ALERT] :exclamation: %s\n%s\n%s", alert.Name, alert.Message, data),
			},
		},
	})

}

func (s *Slack) ResolveAlert(alert *Alert, message string) error {
	return s.sendMessage(SlackWebhookRequest{
		Attachments: []SlackAttachment{
			{
				Color: "#3bb143",
				Text:  fmt.Sprintf("[RESOLVED] :tick: %s\n%s", alert.Name, message),
			},
		},
	})
}

func (s *Slack) Name() string {
	return s.name
}

func (s *Slack) sendMessage(r SlackWebhookRequest) error {
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
