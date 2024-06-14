package feeds

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

type JobConfig interface {
	DefaultHTTPTimeout() models.Duration
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
