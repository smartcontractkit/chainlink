package config

import "time"

type JobSpecReporter interface {
	Enabled() bool
	PollingInterval() time.Duration
	// EnabledOCR2PluginTypes is the allowlist of OCR2 plugin types to emit for
	// (e.g. "median"). An empty slice disables all. Use ["all"] to enable all types.
	EnabledOCR2PluginTypes() []string
}
