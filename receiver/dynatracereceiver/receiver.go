package dynatracereceiver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
)

type Receiver struct {
	Config   *Config
	ticker   *time.Ticker
	stopChan chan struct{}
}

type DynatraceResponse struct {
	TotalCount  int                   `json:"totalCount"`
	NextPageKey string                `json:"nextPageKey"`
	Resolution  string                `json:"resolution"`
	Result      []DynatraceMetricData `json:"result"`
}

type DynatraceMetricData struct {
	MetricID string         `json:"metricId"`
	Data     []MetricValues `json:"data"`
}

type MetricValues struct {
	Timestamps []int64   `json:"timestamps"`
	Values     []float64 `json:"values"`
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

	metrics, err := fetchAllDynatraceMetrics(apiEndpoint, apiToken)
	if err != nil {
		return err
	}

	printOTelMetrics(metrics)

	return nil
}

func fetchAllDynatraceMetrics(apiEndpoint string, apiToken string) ([]DynatraceMetricData, error) {

	url := createMetricsQuery(apiEndpoint)

	resp, err := makeHttpRequest(url, apiToken)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	body, err := readResponseBody(resp)
	if err != nil {
		return nil, err
	}

	var dtResponse DynatraceResponse
	if err := json.Unmarshal(body, &dtResponse); err != nil {
		return nil, fmt.Errorf("json unmarshal failed: %w", err)
	}

	fmt.Println("Fetching data from:", url)
	fmt.Println("Raw Response:", string(body))

	return dtResponse.Result, nil
}

func createMetricsQuery(apiEndpoint string) (url string) {
	metrics := []string{
		"dsfm:active_gate.jvm.cpu_usage",
		"builtin:billing.log.ingest.usage",
	}

	metricSelector := strings.Join(metrics, ",")
	url = fmt.Sprintf("%s?metricSelector=%s&resolution=1h&from=now-24h&to=now", apiEndpoint, metricSelector)

	fmt.Println("Fetching data from:", url)

	return

}

func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body failed: %w", err)
	}
	return body, nil
}

func makeHttpRequest(url, apiToken string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Api-Token "+apiToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func printOTelMetrics(metrics []DynatraceMetricData) {
	fmt.Printf("\nConverted OpenTelemetry Metrics (%d metrics found):\n", len(metrics))

	for _, metric := range metrics {
		fmt.Printf("\nMetric Name: %s\n", metric.MetricID)
		fmt.Println("---------------------------------")

		for _, data := range metric.Data {
			for i, ts := range data.Timestamps {
				if i < len(data.Values) {
					fmt.Printf("Timestamp: %d | Value: %.2f\n", ts, data.Values[i])
				}
			}
		}
	}
}
