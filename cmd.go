package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var debug bool

func notifyError(service NotifyService, err error) {
	if err != nil {
		if serviceErr := service.Notify(err.Error()); serviceErr != nil {
			log.Fatal(serviceErr)
		}
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use:   "lqt",
	Short: "lqt is a quote notifier service",
}

var notifyCmd = &cobra.Command{
	Use:   "notify",
	Short: "Sends a love quote to Telegram",
	Run: func(cmd *cobra.Command, args []string) {
		factory := NotifyServiceFactory{}

		osService, err := factory.CreateService("os")
		if err != nil {
			log.Fatal(err)
		}

		start := time.Now()
		quote, err := GetLoveQuote()
		if debug {
			fmt.Println("GetLoveQuote: ", time.Since(start))
		}
		notifyError(osService, err)

		start = time.Now()
		translatedText, err := quote.TranslateToPT()
		if debug {
			fmt.Println("TranslateToPT: ", time.Since(start))
		}
		notifyError(osService, err)

		telegramService, err := factory.CreateService("telegram")
		notifyError(osService, err)

		formatedMessage := fmt.Sprintf(
			"<blockquote>%s</blockquote><b>Por: </b>%s",
			translatedText[0],
			translatedText[1],
		)

		start = time.Now()
		err = telegramService.Notify(formatedMessage)
		if debug {
			fmt.Println("Telegram service notify: ", time.Since(start))
		}
		notifyError(osService, err)

		err = osService.Notify("Telegram message delivered")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func Init() {
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "execute in debug mode")
	rootCmd.AddCommand(notifyCmd)
}
