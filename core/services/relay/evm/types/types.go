package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"

	"gopkg.in/guregu/null.v2"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type RelayConfig struct {
	ChainID                *utils.Big  `json:"chainID"`
	FromBlock              uint64      `json:"fromBlock"`
	EffectiveTransmitterID null.String `json:"effectiveTransmitterID"`

	// Contract-specific
	SendingKeys pq.StringArray `json:"sendingKeys"`

	// Mercury-specific
	FeedID *common.Hash `json:"feedID"`
}
