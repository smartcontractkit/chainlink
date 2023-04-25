package config

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
)

// OCR1Config is a subset of global config relevant to OCR v1.
type OCR1Config interface {
	// OCR1 config, can override in jobs, only ethereum.
	OCRBlockchainTimeout() time.Duration
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRKeyBundleID() (string, error)
	OCRObservationTimeout() time.Duration
	OCRSimulateTransactions() bool
	OCRTransmitterAddress() (ethkey.EIP55Address, error) // OCR2 can support non-evm changes
	// OCR1 config, cannot override in jobs
	OCRTraceLogging() bool
	OCRDefaultTransactionQueueDepth() uint32
	OCRCaptureEATelemetry() bool
}
