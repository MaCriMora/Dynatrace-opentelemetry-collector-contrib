package dynatraceexporter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type Exporter struct {
	target string
}

func NewSimpleExporter() *Exporter {
	return &Exporter{}
}

func (e *Exporter) Start(_ context.Context, _ component.Host) error {
	fmt.Println("Simple Exporter started. Target:", e.target)
	return nil
}

func (e *Exporter) Shutdown(_ context.Context) error {
	fmt.Println("Simple Exporter shutting down.")
	return nil
}

func (e *Exporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	fmt.Println("Simple Exporter received metrics:", md.MetricCount())
	fmt.Println(md) // Just for testing, print the incoming data from dynatracereceiver
	return nil
}
