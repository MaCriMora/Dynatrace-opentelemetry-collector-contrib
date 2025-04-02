package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/MaCriMora/Dynatrace-opentelemetry-collector-contrib/receiver/dynatracereceiver"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type DummyConsumer struct{}

func (d *DummyConsumer) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	log.Printf("DummyConsumer received metrics: %d\n", md.MetricCount())
	return nil
}

func (d *DummyConsumer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

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

	viper.SetConfigFile("../config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}

	metricSelectors := viper.GetStringSlice("receivers.dynatrace.metric_selectors")
	if len(metricSelectors) == 0 {
		log.Fatal("No metric selectors found in config.yaml")
	}

	config := &dynatracereceiver.Config{
		APIEndpoint:     apiEndpoint,
		APIToken:        apiToken,
		MetricSelectors: metricSelectors,
	}

	receiver := &dynatracereceiver.Receiver{
		Config:     config,
		NextMetric: &DummyConsumer{},
	}

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
