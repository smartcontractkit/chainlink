package feeds

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Config --output ./mocks/ --case=underscore

type Config interface {
	Dev() bool
	FeatureOffchainReporting() bool
	DefaultHTTPTimeout() models.Duration
	OCRBlockchainTimeout() time.Duration
	OCRContractConfirmations() uint16
	OCRContractPollInterval() time.Duration
	OCRContractSubscribeInterval() time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRObservationTimeout() time.Duration
	OCRObservationGracePeriod() time.Duration
	LogSQL() bool
}
