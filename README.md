# Notify
A project to practice factory and abstract factory patterns by implementing
a notification tool that can send a notification via System notification, 
SMS, E-mail, Telegram or any other i want to.

## Design
CLI tool
input text
destination (SMS, E-mail, Telegram, etc)

```go
type NotifyService interface {
    notify(text string) bool
}

type Notifier interface {
    // Factory method
    createNotifier() NotifyService
}
```

## Todo
- SNS integration
- IAM best practices
- Go `context` package
- Go `net` package
- Go `http` package
- Go `json` package

## APIs Utilized
- [Love Quote API](https://rapidapi.com/colebidex-mO-Ew1CYzUS/api/love-quote)
