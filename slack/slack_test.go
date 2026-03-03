package slack_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aceextensions/notifier/slack"
)

// TestSlack_Send_Success spins up a local HTTP server to simulate Slack's
// webhook endpoint — no real network call needed.
func TestSlack_Send_Success(t *testing.T) {
	// Create a mock HTTP server that returns 200 OK
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := slack.New(server.URL)
	err := client.Send("Test Alert", "CPU is at 95%")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSlack_Send_ServerError(t *testing.T) {
	// Mock server returns 500
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := slack.New(server.URL)
	err := client.Send("Test", "Body")
	if err == nil {
		t.Fatal("expected error for non-200 response, got nil")
	}
}

func TestSlack_Send_InvalidURL(t *testing.T) {
	client := slack.New("http://127.0.0.1:0") // nothing listening here
	err := client.Send("Test", "Body")
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}

// TestSlack_Send_Live sends a REAL message to Slack.
// Only runs when SLACK_WEBHOOK_URL environment variable is set.
//
// Run with:
//
//	SLACK_WEBHOOK_URL="https://hooks.slack.com/..." go test ./slack/... -run TestSlack_Send_Live -v
func TestSlack_Send_Live(t *testing.T) {
	webhookURL := os.Getenv("SLACK_WEBHOOK_URL")
	if webhookURL == "" {
		t.Skip("Skipping live test: SLACK_WEBHOOK_URL not set")
	}

	client := slack.New(webhookURL)
	err := client.Send(
		"✅ DevOps Autopilot - Test Notification",
		"This is a live test from `github.com/aceextensions/notifier/slack` package.\nIf you see this in Slack, the notifier is working! 🚀",
	)
	if err != nil {
		t.Fatalf("live slack send failed: %v", err)
	}

	t.Log("✅ Slack message sent successfully!")
}
