package dynatracereceiver

import (
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	component.Config `mapstructure:",squash"`

	APIEndpoint string `mapstructure:"api_endpoint"`

	APIToken string `mapstructure:"api_token"`
}
