package main

import (
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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/gen2brain/beeep"
	"github.com/joho/godotenv"
)

type NotifyService interface {
	notify(text string) error
}

type SMSService struct {
	SnsClient *sns.Client
}

func (S *SMSService) notify(text string) error {
	input := sns.PublishInput{
		TopicArn: aws.String(os.Getenv("AWS_SNS_TOPIC_ARN")),
		Message:  aws.String(text),
	}
	_, err := S.SnsClient.Publish(context.TODO(), &input)
	return err
}

type OSService struct{}

func (o *OSService) notify(text string) error {
	return beeep.Notify("SMS notification", text, "")
}

type NotifyServiceFactory struct {
	SnsClient *sns.Client
}

func (n NotifyServiceFactory) createService(service string) (NotifyService, error) {
	switch strings.ToLower(service) {
	case "sms":
		return &SMSService{n.SnsClient}, nil
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

	return quote, nil
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env file")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	snsClient := sns.NewFromConfig(cfg)
	factory := NotifyServiceFactory{snsClient}

	osService, err := factory.createService("os")
	if err != nil {
		log.Fatal(err)
	}

	quote, err := getLoveQuote()
	if err != nil {
		err = osService.notify(err.Error())
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	translatedText, err := quote.translateToPT()
	if err != nil {
		err = osService.notify(err.Error())
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	smsService, err := factory.createService("sms")
	if err != nil {
		err = osService.notify(err.Error())
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	err = smsService.notify(translatedText)
	if err != nil {
		err = osService.notify(err.Error())
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(err)
	}

	err = osService.notify("SMS delivered")
	if err != nil {
		log.Fatal(err)
	}
}
