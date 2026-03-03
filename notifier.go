// Package notifier provides a common interface for sending notifications
// across multiple channels (Slack, Email, etc.).
//
// Usage:
//
//	import "github.com/aceextensions/notifier"
//
//	var n notifier.Notifier = slack.New("https://hooks.slack.com/...")
//	n.Send("Alert", "CPU is at 95%")
package notifier

// Notifier is the common interface all notification channels must implement.
// Any type that implements Send can be used wherever a Notifier is expected.
type Notifier interface {
	// Send sends a notification with the given subject and body.
	// Returns an error if the delivery failed.
	Send(subject, body string) error
}

// Multi dispatches a notification to multiple Notifiers.
// It continues even if one fails and returns all errors combined.
type Multi struct {
	notifiers []Notifier
}

// NewMulti creates a Multi notifier from a list of Notifiers.
func NewMulti(notifiers ...Notifier) *Multi {
	return &Multi{notifiers: notifiers}
}

// Send sends the notification to all registered notifiers.
// Returns the first error encountered (but still calls all notifiers).
func (m *Multi) Send(subject, body string) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Send(subject, body); err != nil {
			lastErr = err
		}
	}
	return lastErr
}

// Ensure Multi implements Notifier at compile time.
var _ Notifier = (*Multi)(nil)
