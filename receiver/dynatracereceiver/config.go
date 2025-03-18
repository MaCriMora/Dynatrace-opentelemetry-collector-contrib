package dynatracereceiver

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the Dynatrace receiver.
type Config struct {
	component.Config `mapstructure:",squash"`

	// APIEndpoint is the Dynatrace API URL.
	APIEndpoint string `mapstructure:"api_endpoint"`

	// APIToken is the token used for Dynatrace API authentication.
	APIToken string `mapstructure:"api_token"`
}
