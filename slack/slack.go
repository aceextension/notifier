package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type slackNotifier struct {
	webhookURL string
}

// New creates a new slack notifier
func New(webhookURL string) *slackNotifier {
	return &slackNotifier{webhookURL: webhookURL}
}

type payload struct {
	Text string `json:"text"`
}

// Send sends a slack message
func (s *slackNotifier) Send(subject, body string) error {
	p := payload{
		Text: fmt.Sprintf("*%s*\n%s", subject, body),
	}
	
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	
	resp, err := http.Post(s.webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack status conflict: %s", resp.Status)
	}
	
	return nil
}
