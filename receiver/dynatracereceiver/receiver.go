package dynatracereceiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/collector/component"
)

// Receiver represents your Dynatrace receiver instance.
type Receiver struct {
	Config   *Config
	ticker   *time.Ticker
	stopChan chan struct{}
}

type DynatraceMetric struct {
	MetricID    string `json:"metricId"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
}

type DynatraceResponse struct {
	TotalCount int               `json:"totalCount"`
	Metrics    []DynatraceMetric `json:"metrics"`
}

type OTMetric struct {
	Name        string
	Description string
	Unit        string
}

// Start initiates periodic fetching from Dynatrace.
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	fmt.Println("Dynatrace Receiver started with config:", r.Config)

	r.ticker = time.NewTicker(30 * time.Second) // Fetch data every 30 seconds
	r.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-r.ticker.C:
				err := SimpleTestPull(r.Config.APIEndpoint, r.Config.APIToken)
				if err != nil {
					fmt.Println("Error pulling metrics:", err)
				}
			case <-r.stopChan:
				fmt.Println("Stopping Dynatrace Receiver polling loop.")
				return
			}
		}
	}()

	return nil
}

// Shutdown stops the periodic fetching loop.
func (r *Receiver) Shutdown(ctx context.Context) error {
	fmt.Println("Dynatrace Receiver shutting down.")
	r.ticker.Stop()
	close(r.stopChan)
	return nil
}

// SimpleTestPull fetches Dynatrace metrics and converts them to OpenTelemetry format.
func SimpleTestPull(apiEndpoint string, apiToken string) error {
	client := &http.Client{}

	// Make a request to Dynatrace API
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Api-Token "+apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dtResponse DynatraceResponse
	if err := json.Unmarshal(body, &dtResponse); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// Convert Dynatrace Metrics to OpenTelemetry format
	var otelMetrics []OTMetric

	for _, metric := range dtResponse.Metrics {
		otMetric := OTMetric{
			Name:        metric.MetricID,
			Description: metric.Description,
			Unit:        metric.Unit,
		}
		otelMetrics = append(otelMetrics, otMetric)
	}

	// Print converted OpenTelemetry metrics
	fmt.Printf("\nConverted %d OpenTelemetry Metrics:\n", len(otelMetrics))
	for _, otMetric := range otelMetrics {
		fmt.Printf("Name: %s | Description: %s | Unit: %s\n", otMetric.Name, otMetric.Description, otMetric.Unit)
	}

	return nil
}
