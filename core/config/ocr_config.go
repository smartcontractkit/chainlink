package config

import (
	"time"

	"github.com/smartcontractkit/chainlink-evm/pkg/types"
)

// OCR is a subset of global config relevant to OCR v1.
type OCR interface {
	Enabled() bool
	// OCR1 config, can override in jobs, only ethereum.
	BlockchainTimeout() time.Duration
	ContractPollInterval() time.Duration
	ContractSubscribeInterval() time.Duration
	KeyBundleID() (string, error)
	ObservationTimeout() time.Duration
	SimulateTransactions() bool
	TransmitterAddress() (types.EIP55Address, error) // OCR2 can support non-evm changes
	// OCR1 config, cannot override in jobs
	TraceLogging() bool
	DefaultTransactionQueueDepth() uint32
	CaptureEATelemetry() bool
}
