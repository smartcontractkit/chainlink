package config

import (
	"time"
)

// OCR2Config is a subset of global config relevant to OCR v2.
type OCR2Config interface {
	// OCR2 config, can override in jobs, all chains
	OCR2ContractConfirmations() uint16
	OCR2ContractTransmitterTransmitTimeout() time.Duration
	OCR2BlockchainTimeout() time.Duration
	OCR2DatabaseTimeout() time.Duration
	OCR2ContractPollInterval() time.Duration
	OCR2ContractSubscribeInterval() time.Duration
	OCR2KeyBundleID() (string, error)
	// OCR2 config, cannot override in jobs
	OCR2TraceLogging() bool
	OCR2CaptureEATelemetry() bool
}
