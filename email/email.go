package email

import (
	"fmt"
	"net/smtp"
	"strings"
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
	
	// Construct the email message with explicit headers
	message := fmt.Sprintf("From: %s\r\n", e.cfg.From)
	message += fmt.Sprintf("To: %s\r\n", strings.Join(e.cfg.To, ","))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n" // Separate headers from body
	message += body

	addr := fmt.Sprintf("%s:%d", e.cfg.SMTPHost, e.cfg.SMTPPort)
	
	return smtp.SendMail(addr, auth, e.cfg.From, e.cfg.To, []byte(message))
}
