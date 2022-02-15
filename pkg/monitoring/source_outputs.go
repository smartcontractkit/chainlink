package monitoring

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/libocr/offchainreporting2/types"
)

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
	BlockNumber uint64
	Transmitter types.Account
	LinkBalance *big.Int

	// The "fee coin" is different for each chain.
	JuelsPerFeeCoin   *big.Int
	AggregatorRoundID uint32
}

type TxResults struct {
	NumSucceeded, NumFailed uint64
}
