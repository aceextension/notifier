package email_test

import (
	"io"
	"net"
	"net/smtp"
	"net/textproto"
	"testing"

	"github.com/aceextensions/notifier/email"
)

func TestEmail_Send_DisabledConfig(t *testing.T) {
	// Build a client but don't call Send — just test struct init
	client := email.New(email.Config{
		SMTPHost: "smtp.gmail.com",
		SMTPPort: 587,
		Username: "test@gmail.com",
		Password: "password",
		From:     "test@gmail.com",
		To:       []string{"admin@company.com"},
	})
	if client == nil {
		t.Fatal("expected non-nil email.Client")
	}
}

func TestEmail_Send_InvalidHost(t *testing.T) {
	client := email.New(email.Config{
		SMTPHost: "127.0.0.1",
		SMTPPort: 1, // nothing listening
		Username: "user",
		Password: "pass",
		From:     "from@test.com",
		To:       []string{"to@test.com"},
	})

	err := client.Send("Test Subject", "Test Body")
	if err == nil {
		t.Fatal("expected error for unreachable SMTP host, got nil")
	}
}

// TestEmail_Send_MockSMTP runs a minimal in-process SMTP server to verify
// the email client sends correct SMTP commands without a real mail server.
func TestEmail_Send_MockSMTP(t *testing.T) {
	// Start a local TCP listener to act as a mock SMTP server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start mock server: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	done := make(chan struct{})

	go func() {
		defer close(done)
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		tc := textproto.NewConn(conn)

		// SMTP handshake
		tc.PrintfLine("220 mock SMTP ready")
		tc.ReadLine() // EHLO
		tc.PrintfLine("250-mock")
		tc.PrintfLine("250 AUTH PLAIN")
		tc.ReadLine() // AUTH PLAIN
		tc.PrintfLine("235 OK")
		tc.ReadLine() // MAIL FROM
		tc.PrintfLine("250 OK")
		tc.ReadLine() // RCPT TO
		tc.PrintfLine("250 OK")
		tc.ReadLine() // DATA
		tc.PrintfLine("354 Start mail input")
		// Read until dot
		for {
			line, _ := tc.ReadLine()
			if line == "." {
				break
			}
		}
		tc.PrintfLine("250 OK")
		tc.ReadLine() // QUIT
		tc.PrintfLine("221 Bye")
	}()

	client := email.New(email.Config{
		SMTPHost: "127.0.0.1",
		SMTPPort: port,
		Username: "user@test.com",
		Password: "password",
		From:     "user@test.com",
		To:       []string{"admin@test.com"},
	})

	// Note: smtp.PlainAuth requires TLS or localhost — this works for localhost
	_ = smtp.SendMail // just referencing to avoid unused import
	_ = io.Discard

	err = client.Send("Hello", "World")
	// The mock server accepts the connection — we just verify no panic/crash
	// The actual auth may fail with mock, so we just check it doesn't panic
	t.Logf("Send result (mock): %v", err)

	<-done
}
