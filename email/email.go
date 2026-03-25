package email

import (
	"fmt"
	"net/smtp"
)

// Config represents email configuration
type Config struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
	To       []string
}

type emailNotifier struct {
	cfg Config
}

// New creates a new email notifier
func New(cfg Config) *emailNotifier {
	return &emailNotifier{cfg: cfg}
}

// Send sends an email
func (e *emailNotifier) Send(subject, body string) error {
	auth := smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.SMTPHost)
	
	msg := []byte(fmt.Sprintf("To: %v\r\nSubject: %s\r\n\r\n%s\r\n", e.cfg.To, subject, body))
	
	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)
	
	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, msg)
}
