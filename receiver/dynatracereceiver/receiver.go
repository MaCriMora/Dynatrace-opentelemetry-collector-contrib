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

// todo: add pagination variable for more data -> current cap at 100
type DynatraceResponse struct {
	TotalCount int               `json:"totalCount"`
	Metrics    []DynatraceMetric `json:"metrics"`
}

type OTMetric struct {
	Name        string
	Description string
	Unit        string
}

// start polling from Dynatrace.
func (r *Receiver) Start(ctx context.Context, host component.Host) error {
	fmt.Println("Dynatrace Receiver started with config:", r.Config)

	r.ticker = time.NewTicker(30 * time.Second) // Fetch data every 30 seconds
	r.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-r.ticker.C:
				err := pullDynatraceMetrics(r.Config.APIEndpoint, r.Config.APIToken)
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

func pullDynatraceMetrics(apiEndpoint string, apiToken string) error {
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

	// Read and parse the responseS
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var dtResponse DynatraceResponse
	if err := json.Unmarshal(body, &dtResponse); err != nil {
		return fmt.Errorf("json unmarshal failed: %w", err)
	}

	// Set to false to print raw Dynatrace metrics or true for formatted OTel metrics depending what you want/need
	// not really neccessary later but good to understand the metrics from dynatrace
	formatOTelMetrics := true

	if formatOTelMetrics {
		PrintOTelMetrics(dtResponse)
	} else {
		PrintRawDynatraceMetrics(dtResponse)
	}

	return nil
}

func PrintRawDynatraceMetrics(response DynatraceResponse) {
	fmt.Printf("\nRaw %d Dynatrace Metrics:\n", response.TotalCount)
	for _, metric := range response.Metrics {
		fmt.Printf("MetricID: %s | DisplayName: %s | Description: %s | Unit: %s\n",
			metric.MetricID, metric.DisplayName, metric.Description, metric.Unit)
	}
}

func PrintOTelMetrics(response DynatraceResponse) {
	var otelMetrics []OTMetric

	for _, metric := range response.Metrics {
		otMetric := OTMetric{
			Name:        metric.MetricID,
			Description: metric.Description,
			Unit:        metric.Unit,
		}
		otelMetrics = append(otelMetrics, otMetric)
	}

	fmt.Printf("\nConverted %d OpenTelemetry Metrics:\n", len(otelMetrics))
	for _, otMetric := range otelMetrics {
		fmt.Printf("Name: %s | Description: %s | Unit: %s\n", otMetric.Name, otMetric.Description, otMetric.Unit)
	}
}
