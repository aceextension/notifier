// Package email provides an SMTP email notifier that implements
// the github.com/aceextensions/notifier.Notifier interface.
//
// Usage:
//
//	import "github.com/aceextensions/notifier/email"
//
//	n := email.New(email.Config{
//	    SMTPHost: "smtp.gmail.com",
//	    SMTPPort: 587,
//	    Username: "you@gmail.com",
//	    Password: "your-app-password",
//	    From:     "you@gmail.com",
//	    To:       []string{"admin@company.com"},
//	})
//	err := n.Send("Server Alert", "CPU usage at 95% on prod-01")
package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// Config holds the SMTP connection and routing configuration.
type Config struct {
	SMTPHost string   // e.g. "smtp.gmail.com"
	SMTPPort int      // e.g. 587
	Username string   // SMTP auth username
	Password string   // SMTP auth password (use App Password for Gmail)
	From     string   // Sender address
	To       []string // Recipient address(es)
}

// Client sends notifications via SMTP email.
type Client struct {
	cfg Config
}

// New creates a new SMTP email notifier with the provided Config.
func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

// Send sends an email with the given subject and plain-text body.
// Compatible with Gmail, SendGrid, Mailgun, AWS SES SMTP, and any standard SMTP.
func (c *Client) Send(subject, body string) error {
	addr := fmt.Sprintf("%s:%d", c.cfg.SMTPHost, c.cfg.SMTPPort)
	auth := smtp.PlainAuth("", c.cfg.Username, c.cfg.Password, c.cfg.SMTPHost)

	to := strings.Join(c.cfg.To, ", ")

	msg := strings.Join([]string{
		"From: " + c.cfg.From,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	if err := smtp.SendMail(addr, auth, c.cfg.From, c.cfg.To, []byte(msg)); err != nil {
		return fmt.Errorf("email: send mail: %w", err)
	}

	return nil
}
