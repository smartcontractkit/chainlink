package solana

import (
	"github.com/smartcontractkit/chainlink/core/config/toml"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type TOMLChain struct {
	BalancePollPeriod   *models.Duration
	ConfirmPollPeriod   *models.Duration
	OCR2CachePollPeriod *models.Duration
	OCR2CacheTTL        *models.Duration
	TxTimeout           *models.Duration
	TxRetryTimeout      *models.Duration
	TxConfirmTimeout    *models.Duration
	SkipPreflight       *bool
	Commitment          *string
	MaxRetries          *int
}

type TOMLNode struct {
	Name string
	URL  *toml.URL
}
