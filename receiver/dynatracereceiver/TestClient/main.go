package main

import (
	"context"
	"log"
	"time"

	"github.com/macrimo/opentelemetry-collector-contrib/receiver/dynatracereceiver"
)

func main() {
	config := &dynatracereceiver.Config{
		APIEndpoint: "", // tofo -> find correct api enpoint and edit query to fetch correct data
		APIToken:    "",
	}
	// add dynatrace endpoint & APItoken

	receiver := &dynatracereceiver.Receiver{Config: config}

	err := receiver.Start(context.Background(), nil)
	if err != nil {
		log.Fatal("Error starting receiver:", err)
	}

	// currently just set to 2min for less logs but can be increased or just left out. Was just for testing
	time.Sleep(2 * time.Minute)

	err = receiver.Shutdown(context.Background())
	if err != nil {
		log.Fatal("Error shutting down receiver:", err)
	}

	log.Println("Receiver stopped successfully.")
}
