package config

import "time"

type JobSpecReporter interface {
	Enabled() bool
	PollingInterval() time.Duration
	// EnabledOCR2PluginTypes is the allowlist of OCR2 plugin types to emit for
	// (e.g. "median", "functions"). An empty slice means emit for all types.
	EnabledOCR2PluginTypes() []string
	// EmitNonOCR2Jobs toggles emission for supported non-OCR2 job types.
	// Defaults to false.
	EmitNonOCR2Jobs() bool
}
