package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/MaCriMora/Dynatrace-opentelemetry-collector-contrib/receiver/dynatracereceiver"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	apiEndpoint := os.Getenv("API_ENDPOINT")
	apiToken := os.Getenv("API_TOKEN")

	if apiEndpoint == "" || apiToken == "" {
		log.Fatal("API credentials missing. Check your .env file.")
	}

	config := &dynatracereceiver.Config{
		APIEndpoint: apiEndpoint,
		APIToken:    apiToken,
	}

	receiver := &dynatracereceiver.Receiver{Config: config}

	err = receiver.Start(context.Background(), nil)
	if err != nil {
		log.Fatal("Error starting receiver:", err)
	}

	time.Sleep(2 * time.Minute)

	err = receiver.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Error shutting down receiver:", err)
	}

	log.Println("Receiver stopped successfully.")
}
