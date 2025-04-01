package dynatracereceiver

import (
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	component.Config `mapstructure:",squash"`

	APIEndpoint string `mapstructure:"API_ENDPOINT"`

	APIToken string `mapstructure:"api_token"`

	MetricSelectors []string `mapstructure:"metric_selectors"`
}
