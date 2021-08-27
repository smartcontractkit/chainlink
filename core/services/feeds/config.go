package feeds

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/chains"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

//go:generate mockery --name Config --output ./mocks/ --case=underscore

type Config interface {
	Chain() *chains.Chain
	ChainID() *big.Int
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
}
