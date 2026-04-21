package config

import "time"

type JobSpecReporter interface {
	Enabled() bool
	PollingInterval() time.Duration
	// EnabledOCR2PluginTypes is the allowlist of OCR2 plugin types to emit for (e.g. "median", "ocr2keeper").
	// An empty slice means all OCR2 plugin types are allowed.
	EnabledOCR2PluginTypes() []string
	// EmitNonOCR2Jobs controls whether non-OCR2 jobs (VRF, Keeper, Functions, CCIP, Workflow, …)
	// emit the generic envelope. Default false for the initial rollout.
	EmitNonOCR2Jobs() bool
}
