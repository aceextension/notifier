# 📬 notifier

A clean, zero-dependency Go package providing a common `Notifier` interface with ready-to-use **Slack** and **Email** implementations.

[![Go Reference](https://pkg.go.dev/badge/github.com/aceextensions/notifier.svg)](https://pkg.go.dev/github.com/aceextensions/notifier)

---

## Install

```bash
go get github.com/aceextensions/notifier@latest
```

---

## Usage

### Slack

```go
import "github.com/aceextensions/notifier/slack"

n := slack.New("https://hooks.slack.com/services/XXX/YYY/ZZZ")
err := n.Send("Server Alert", "CPU at 95% on prod-01")
```

### Email

```go
import "github.com/aceextensions/notifier/email"

n := email.New(email.Config{
    SMTPHost: "smtp.gmail.com",
    SMTPPort: 587,
    Username: "you@gmail.com",
    Password: "your-app-password",
    From:     "you@gmail.com",
    To:       []string{"admin@company.com"},
})
err := n.Send("Server Alert", "CPU at 95% on prod-01")
```

### Multiple channels at once

```go
import (
    "github.com/aceextensions/notifier"
    "github.com/aceextensions/notifier/slack"
    "github.com/aceextensions/notifier/email"
)

multi := notifier.NewMulti(
    slack.New(os.Getenv("SLACK_WEBHOOK_URL")),
    email.New(email.Config{ /* ... */ }),
)

multi.Send("Alert", "Disk is at 92% on prod-01")
```

### Use the interface in your own code

```go
import "github.com/aceextensions/notifier"

func SendAlert(n notifier.Notifier, msg string) error {
    return n.Send("DevOps Alert", msg)
}
```

---

## Testing

### Run all unit tests (no network required)

```bash
go test ./... -v
```

Expected output:
```
PASS  github.com/aceextensions/notifier         (4 tests)
PASS  github.com/aceextensions/notifier/email   (3 tests)
PASS  github.com/aceextensions/notifier/slack   (3 pass, 1 skipped)
```

The live Slack test is **skipped by default** — it only runs when you provide a real webhook URL.

---

### Test live Slack notification (sends a real message)

```bash
SLACK_WEBHOOK_URL="https://hooks.slack.com/services/TXXXXXXX/BXXXXXXX/XXXXXXXXX" \
  go test ./slack/... -run TestSlack_Send_Live -v
```

If successful you'll see:
```
--- PASS: TestSlack_Send_Live (0.31s)
    slack_test.go: ✅ Slack message sent successfully!
```

And a real message will appear in your Slack channel like this:

> **✅ DevOps Autopilot - Test Notification**
> This is a live test from `github.com/aceextensions/notifier/slack` package.
> If you see this in Slack, the notifier is working! 🚀

### Other useful test commands

```bash
go test ./...           # run all tests
go test ./... -v        # verbose output
go test ./... -cover    # show test coverage
go test ./slack/... -v  # test only Slack package
go test ./email/... -v  # test only Email package
```

---

## Package Layout

```
github.com/aceextensions/notifier       ← Notifier interface + Multi dispatcher
github.com/aceextensions/notifier/slack ← Slack Incoming Webhook
github.com/aceextensions/notifier/email ← SMTP Email
```

---

## Interface

```go
type Notifier interface {
    Send(subject, body string) error
}
```

Any type implementing `Send` is a valid `Notifier`. You can write your own (PagerDuty, SMS, Telegram, etc.) and plug it into `notifier.NewMulti`.

---

## License

MIT
