// Package slack provides a Slack webhook notifier that implements
// the github.com/aceextensions/notifier.Notifier interface.
//
// Usage:
//
//	import "github.com/aceextensions/notifier/slack"
//
//	n := slack.New("https://hooks.slack.com/services/XXX/YYY/ZZZ")
//	err := n.Send("Server Alert", "CPU usage at 95% on prod-01")
package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client sends notifications to a Slack channel via an Incoming Webhook.
type Client struct {
	webhookURL string
	httpClient *http.Client
}

// Option is a functional option for configuring the Slack Client.
type Option func(*Client)

// WithHTTPClient sets a custom *http.Client (useful for testing or proxies).
func WithHTTPClient(c *http.Client) Option {
	return func(s *Client) {
		s.httpClient = c
	}
}

// New creates a new Slack notifier.
//
//	webhookURL: Your Slack Incoming Webhook URL
//	            (Settings → Integrations → Incoming Webhooks)
func New(webhookURL string, opts ...Option) *Client {
	c := &Client{
		webhookURL: webhookURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Send posts a message to Slack. The subject is used as the message title
// (bold) and body is appended below it.
func (c *Client) Send(subject, body string) error {
	text := fmt.Sprintf("*%s*\n%s", subject, body)
	payload := map[string]string{"text": text}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal payload: %w", err)
	}

	resp, err := c.httpClient.Post(c.webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("slack: http post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: unexpected status code %d", resp.StatusCode)
	}

	return nil
}
