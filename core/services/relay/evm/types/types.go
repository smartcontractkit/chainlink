package types

import (
	"github.com/lib/pq"

	"gopkg.in/guregu/null.v2"

	"github.com/smartcontractkit/chainlink/core/utils"
)

type RelayConfig struct {
	ChainID   *utils.Big `json:"chainID"`
	FromBlock uint64     `json:"fromBlock"`
	// TODO: The key-specific stuff ought to be moved into plugin config since its not relevant for mercury
	EffectiveTransmitterAddress null.String    `json:"effectiveTransmitterAddress"`
	SendingKeys                 pq.StringArray `json:"sendingKeys"`
}
