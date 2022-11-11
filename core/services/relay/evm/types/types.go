package types

import (
	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
	"gopkg.in/guregu/null.v2"
)

type MercuryConfig struct {
	FeedID string      `json:"feedID"`
	URL    *models.URL `json:"url"`
}

type RelayConfig struct {
	MercuryConfig               *MercuryConfig
	ChainID                     *utils.Big     `json:"chainID"`
	FromBlock                   uint64         `json:"fromBlock"`
	EffectiveTransmitterAddress null.String    `json:"effectiveTransmitterAddress"`
	SendingKeys                 pq.StringArray `json:"sendingKeys"`
}
