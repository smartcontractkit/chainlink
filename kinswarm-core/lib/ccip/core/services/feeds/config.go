package feeds

import (
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	coreconfig "github.com/smartcontractkit/chainlink/v2/core/config"
)

type GeneralConfig interface {
	OCR() coreconfig.OCR
	Insecure() coreconfig.Insecure
}

type FeatureConfig interface {
	MultiFeedsManagers() bool
}

type JobConfig interface {
	DefaultHTTPTimeout() commonconfig.Duration
}

type InsecureConfig interface {
	OCRDevelopmentMode() bool
}

type OCRConfig interface {
	Enabled() bool
}

type OCR2Config interface {
	Enabled() bool
	BlockchainTimeout() time.Duration
	ContractConfirmations() uint16
	ContractPollInterval() time.Duration
	ContractTransmitterTransmitTimeout() time.Duration
	DatabaseTimeout() time.Duration
	DefaultTransactionQueueDepth() uint32
	SimulateTransactions() bool
	TraceLogging() bool
}
