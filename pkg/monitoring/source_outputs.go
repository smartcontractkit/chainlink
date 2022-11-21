package monitoring

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

// Envelope contains data that is required from all the chain integrations.
// Integrators usually create an EnvelopeSource to produce Envelope instances.
type Envelope struct {
	// latest transmission details
	ConfigDigest    types.ConfigDigest
	Epoch           uint32
	Round           uint8
	LatestAnswer    *big.Int
	LatestTimestamp time.Time

	// latest contract config
	ContractConfig types.ContractConfig

	// extra
	BlockNumber             uint64
	Transmitter             types.Account
	LinkBalance             *big.Int
	LinkAvailableForPayment *big.Int

	// The "fee coin" is different for each chain.
	JuelsPerFeeCoin   *big.Int
	AggregatorRoundID uint32
}

// TxResults counts the number of successful and failed transactions in a predetermined window of time.
// Integrators usually create an TxResultsSource to produce TxResults instances.
type TxResults struct {
	NumSucceeded, NumFailed uint64
}
