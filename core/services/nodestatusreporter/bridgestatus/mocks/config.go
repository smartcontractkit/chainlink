package mocks

import (
	"time"
)

// TestBridgeStatusReporterConfig implements config.BridgeStatusReporter for testing
type TestBridgeStatusReporterConfig struct {
	enabled              bool
	statusPath           string
	pollingInterval      time.Duration
	ignoreInvalidBridges bool
	ignoreJoblessBridges bool
}

func NewTestBridgeStatusReporterConfig(enabled bool, statusPath string, pollingInterval time.Duration) *TestBridgeStatusReporterConfig {
	return &TestBridgeStatusReporterConfig{
		enabled:              enabled,
		statusPath:           statusPath,
		pollingInterval:      pollingInterval,
		ignoreInvalidBridges: true,
		ignoreJoblessBridges: false,
	}
}

func NewTestBridgeStatusReporterConfigWithSkip(enabled bool, statusPath string, pollingInterval time.Duration, ignoreInvalidBridges bool, ignoreJoblessBridges bool) *TestBridgeStatusReporterConfig {
	return &TestBridgeStatusReporterConfig{
		enabled:              enabled,
		statusPath:           statusPath,
		pollingInterval:      pollingInterval,
		ignoreInvalidBridges: ignoreInvalidBridges,
		ignoreJoblessBridges: ignoreJoblessBridges,
	}
}

func (e *TestBridgeStatusReporterConfig) Enabled() bool {
	return e.enabled
}

func (e *TestBridgeStatusReporterConfig) StatusPath() string {
	return e.statusPath
}

func (e *TestBridgeStatusReporterConfig) PollingInterval() time.Duration {
	return e.pollingInterval
}

func (e *TestBridgeStatusReporterConfig) IgnoreInvalidBridges() bool {
	return e.ignoreInvalidBridges
}

func (e *TestBridgeStatusReporterConfig) IgnoreJoblessBridges() bool {
	return e.ignoreJoblessBridges
}
