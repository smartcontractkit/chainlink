package terra

import (
	"github.com/shopspring/decimal"

	"github.com/smartcontractkit/chainlink/core/config/toml"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

type TOMLChain struct {
	BlockRate             *models.Duration
	BlocksUntilTxTimeout  int
	ConfirmPollPeriod     *models.Duration
	FallbackGasPriceULuna *decimal.Decimal
	FCDURL                *toml.URL
	GasLimitMultiplier    *decimal.Decimal
	MaxMsgsPerBatch       int64
	OCR2CachePollPeriod   *models.Duration
	OCR2CacheTTL          *models.Duration
	TxMsgTimeout          *models.Duration
}

type TOMLNode struct {
	Name          string
	TendermintURL *toml.URL
}
