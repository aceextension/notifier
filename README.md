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
