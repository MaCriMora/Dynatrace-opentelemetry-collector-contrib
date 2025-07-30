#!/bin/bash

# Ensure .env uses Unix line endings (prevents \r bugs on Linux/macOS)
tr -d '\r' < .env > .env.unix && mv .env.unix .env

# Load environment variables from .env file
export $(grep -v '^#' .env | xargs)

# Start the OpenTelemetry Collector
./bin/otelcontribcol_linux_amd64 --config=receiver/dynatracereceiver/config.yaml --set=service.telemetry.logs.level=debug
