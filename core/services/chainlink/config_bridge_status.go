package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.BridgeStatusReporter = (*bridgeStatusReporterConfig)(nil)

type bridgeStatusReporterConfig struct {
	c toml.BridgeStatusReporter
}

func (e *bridgeStatusReporterConfig) Enabled() bool {
	if e.c.Enabled == nil {
		return false
	}
	return *e.c.Enabled
}

func (e *bridgeStatusReporterConfig) StatusPath() string {
	if e.c.StatusPath == nil {
		return "/status"
	}
	return *e.c.StatusPath
}

func (e *bridgeStatusReporterConfig) PollingInterval() time.Duration {
	if e.c.PollingInterval == nil {
		return 5 * time.Minute
	}
	return e.c.PollingInterval.Duration()
}

func (e *bridgeStatusReporterConfig) IgnoreInvalidBridges() bool {
	if e.c.IgnoreInvalidBridges == nil {
		return true
	}
	return *e.c.IgnoreInvalidBridges
}

func (e *bridgeStatusReporterConfig) IgnoreJoblessBridges() bool {
	if e.c.IgnoreJoblessBridges == nil {
		return false
	}
	return *e.c.IgnoreJoblessBridges
}
