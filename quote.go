package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type LoveQuote struct {
	Quote  string `json:"quote"`
	Author string `json:"author"`
}

func (lq *LoveQuote) TranslateToPT() ([]string, error) {
	ctx := context.TODO()
	client, err := translate.NewClient(
		ctx,
		option.WithAPIKey(os.Getenv("GCP_API_KEY")),
	)
	if err != nil {
		return []string{}, err
	}

	translation, err := client.Translate(
		ctx,
		[]string{lq.Quote, lq.Author},
		language.BrazilianPortuguese,
		nil,
	)
	if err != nil {
		return []string{}, err
	}

	if len(translation) == 0 {
		return []string{}, fmt.Errorf(
			"Empty translation for text: %s",
			lq.Quote,
		)
	}

	var translations []string

	for _, translation := range translation {
		translations = append(translations, translation.Text)
	}

	return translations, nil
}

func GetLoveQuote() (LoveQuote, error) {
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

	var quote LoveQuote
	err = json.NewDecoder(res.Body).Decode(&quote)

	if err != nil {
		return LoveQuote{}, err
	}

	if quote.Quote == "" {
		return quote, fmt.Errorf("Quote is empty")
	}

	return quote, nil
}
