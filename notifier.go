package notifier

// Notifier defines the interface for sending notifications
type Notifier interface {
	Send(subject, body string) error
}

// MultiNotifier dispatches to multiple notifiers
type MultiNotifier struct {
	notifiers []Notifier
}

// NewMulti creates a new multi notifier
func NewMulti(notifiers ...Notifier) *MultiNotifier {
	return &MultiNotifier{notifiers: notifiers}
}

// Send sends to all registered notifiers
func (m *MultiNotifier) Send(subject, body string) error {
	var lastErr error
	for _, n := range m.notifiers {
		if err := n.Send(subject, body); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
