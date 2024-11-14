package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

type NotifyService interface {
	notify(text string) error
}

type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type TelegramService struct{}

func (t *TelegramService) notify(text string) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(TelegramMessage{
		os.Getenv("TELEGRAM_GROUP_ID"),
		text,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage",
			os.Getenv("TELEGRAM_BOT_API_TOKEN"),
		),
		&buf,
	)

	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("here", string(payload))
	return nil
}

type OSService struct{}

func (o *OSService) notify(text string) error {
	return beeep.Notify("SMS notification", text, "")
}

type NotifyServiceFactory struct{}

func (n NotifyServiceFactory) createService(service string) (NotifyService, error) {
	switch strings.ToLower(service) {
	case "telegram":
		return &TelegramService{}, nil
	case "os":
		return &OSService{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Service %s isn't valid\n", service))
	}
}

type LoveQuote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

func (lq *LoveQuote) translateToPT() (string, error) {
	ctx := context.TODO()
	client, err := translate.NewClient(ctx, option.WithAPIKey(os.Getenv("GCP_API_KEY")))
	if err != nil {
		return "", err
	}

	translation, err := client.Translate(
		ctx,
		[]string{lq.Quote},
		language.BrazilianPortuguese,
		nil,
	)
	if err != nil {
		return "", err
	}

	if len(translation) == 0 {
		return "", fmt.Errorf("Empty translation for text: %s", lq.Quote)
	}

	return translation[0].Text, nil
}

func getLoveQuote() (LoveQuote, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		"https://love-quote.p.rapidapi.com/lovequote",
		nil,
	)
	if err != nil {
		return LoveQuote{}, err
	}

	req.Header.Add("X-RapidAPI-Key", os.Getenv("X_RAPIDAPI_KEY"))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return LoveQuote{}, err
	}
	defer res.Body.Close()

	payload, err := io.ReadAll(res.Body)
	if err != nil {
		return LoveQuote{}, err
	}

	quote := LoveQuote{}
	err = json.Unmarshal(payload, &quote)

	if err != nil {
		return LoveQuote{}, err
	}

	if quote.Quote == "" {
		return quote, fmt.Errorf("Quote is empty")
	}

	return quote, nil
}

func notifyError(service NotifyService, err error) {
	if err != nil {
		if serviceErr := service.notify(err.Error()); serviceErr != nil {
			log.Fatal(serviceErr)
		}
		log.Fatal(err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env file")
	}

	factory := NotifyServiceFactory{}

	osService, err := factory.createService("os")
	if err != nil {
		log.Fatal(err)
	}

	quote, err := getLoveQuote()
	notifyError(osService, err)

	translatedText, err := quote.translateToPT()
	notifyError(osService, err)

	telegramService, err := factory.createService("telegram")
	notifyError(osService, err)

	err = telegramService.notify(translatedText)
	notifyError(osService, err)

	err = osService.notify("SMS delivered")
	if err != nil {
		log.Fatal(err)
	}
}
