package main

import (
	"log"

	"github.com/joho/godotenv"
)

type NotifyService interface {
	Notify(text string) error
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error while loading env file")
	}
	Init()
	Execute()
}
