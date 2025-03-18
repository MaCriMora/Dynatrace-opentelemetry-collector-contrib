package main

import (
	"context"
	"log"
	"time"

	"github.com/macrimo/opentelemetry-collector-contrib/receiver/dynatracereceiver"
)

func main() {
	config := &dynatracereceiver.Config{
		APIEndpoint: "", // put in your endpoint
		APIToken:    "", // put in your token
	}

	receiver := &dynatracereceiver.Receiver{Config: config}

	// Start the receiver
	err := receiver.Start(context.Background(), nil)
	if err != nil {
		log.Fatal("Error starting receiver:", err)
	}

	// Let it run for a short period (e.g., 2 minutes) to see periodic pulls
	time.Sleep(2 * time.Minute)

	// Shutdown receiver
	err = receiver.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Error shutting down receiver:", err)
	}

	log.Println("Receiver stopped successfully.")
}
