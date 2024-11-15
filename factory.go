package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type NotifyServiceFactory struct{}

func (n NotifyServiceFactory) CreateService(service string) (NotifyService, error) {
	switch strings.ToLower(service) {
	case "telegram":
		var chatID string
		if debug {
			chatID = os.Getenv("TELEGRAM_TEST_CHAT_ID")
		} else {
			chatID = os.Getenv("TELEGRAM_GROUP_ID")
		}

		return &TelegramService{
			chatID,
			os.Getenv("TELEGRAM_BOT_API_TOKEN"),
		}, nil
	case "os":
		return &OSService{}, nil
	default:
		return nil, errors.New(
			fmt.Sprintf("Service %s isn't valid\n", service),
		)
	}
}
