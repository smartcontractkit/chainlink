package config

import (
	"time"
)

// OCR2 is a subset of global config relevant to OCR v2.
type OCR2 interface {
	Enabled() bool
	// OCR2 config, can override in jobs, all chains
	ContractConfirmations() uint16
	ContractTransmitterTransmitTimeout() time.Duration
	BlockchainTimeout() time.Duration
	DatabaseTimeout() time.Duration
	ContractPollInterval() time.Duration
	ContractSubscribeInterval() time.Duration
	KeyBundleID() (string, error)
	// OCR2 config, cannot override in jobs
	TraceLogging() bool
	CaptureEATelemetry() bool
	DefaultTransactionQueueDepth() uint32
	SimulateTransactions() bool
	CaptureAutomationCustomTelemetry() bool
}
