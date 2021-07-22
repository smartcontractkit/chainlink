package feeds

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Config --output ./mocks/ --case=underscore

type Config interface {
	ChainID() *big.Int
	Dev() bool
	FeatureOffchainReporting() bool
	DefaultHTTPTimeout() models.Duration
	OCRBlockchainTimeout(override time.Duration) time.Duration
	OCRContractConfirmations(override uint16) uint16
	OCRContractPollInterval(override time.Duration) time.Duration
	OCRContractSubscribeInterval(override time.Duration) time.Duration
	OCRContractTransmitterTransmitTimeout() time.Duration
	OCRDatabaseTimeout() time.Duration
	OCRObservationTimeout(override time.Duration) time.Duration
	OCRObservationGracePeriod() time.Duration
}
