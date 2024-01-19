package feeds

import (
	"time"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

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
	TraceLogging() bool
}
