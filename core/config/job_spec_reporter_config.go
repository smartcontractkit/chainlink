package config

import "time"

type JobSpecReporter interface {
	Enabled() bool
	PollingInterval() time.Duration
	// EnabledOCR2PluginTypes is the allowlist of OCR2 plugin types to emit for
	// (e.g. "median", "ocr2keeper"). An empty slice means emit for all types.
	EnabledOCR2PluginTypes() []string
	// EmitNonOCR2Jobs toggles emission for non-OCR2 job types (VRF, Keeper,
	// Functions, CCIP, Workflow, …). Defaults to false.
	EmitNonOCR2Jobs() bool
}
