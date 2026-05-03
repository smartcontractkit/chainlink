package chainlink

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/config"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
)

var _ config.JobSpecReporter = (*jobSpecReporterConfig)(nil)

type jobSpecReporterConfig struct {
	c toml.JobSpecReporter
}

func (e *jobSpecReporterConfig) Enabled() bool {
	if e.c.Enabled == nil {
		return false
	}
	return *e.c.Enabled
}

func (e *jobSpecReporterConfig) PollingInterval() time.Duration {
	if e.c.PollingInterval == nil {
		return time.Hour
	}
	return e.c.PollingInterval.Duration()
}

func (e *jobSpecReporterConfig) EnabledOCR2PluginTypes() []string {
	if e.c.EnabledOCR2PluginTypes == nil {
		return []string{"median"}
	}
	return *e.c.EnabledOCR2PluginTypes
}
