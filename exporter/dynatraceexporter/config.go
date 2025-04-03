package dynatraceexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// Config defines configuration for logging exporter.
type Config struct {
	exporterhelper.TimeoutConfig `mapstructure:",squash"`
	TargetEndpoint               string `mapstructure:"target_endpoint"`
	component.Config             `mapstructure:",squash"`
}
