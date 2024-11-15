package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gen2brain/beeep"
)

type telegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

type telegramReturn struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type TelegramService struct {
	ChatID   string
	APIToken string
}

func (t *TelegramService) Notify(text string) error {
	message := telegramMessage{
		t.ChatID,
		text,
		"HTML",
	}

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode(message)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage",
			t.APIToken,
		),
		&buf,
	)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var resp telegramReturn
	err = json.NewDecoder(res.Body).Decode(&resp)

	if !resp.Ok {
		return fmt.Errorf("Failed to send to message to Telegram, %s", resp.Description)
	}

	return nil
}

type OSService struct{}

func (o *OSService) Notify(text string) error {
	return beeep.Notify("LoveQuote:", text, "")
}
