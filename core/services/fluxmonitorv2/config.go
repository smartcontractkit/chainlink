package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Config defines the Flux Monitor configuration.
type Config interface {
	FlagsContractAddress() string
	MinimumContractPayment() *assets.Link
	EvmGasLimitDefault() uint32
	EvmGasLimitFMJobType() *uint32
	EvmMaxQueuedTransactions() uint64
	FMDefaultTransactionQueueDepth() uint32
	pg.QConfig
}

type JobPipelineConfig interface {
	DefaultHTTPTimeout() models.Duration
}

// MinimumPollingInterval returns the minimum duration between polling ticks
func MinimumPollingInterval(c JobPipelineConfig) time.Duration {
	return c.DefaultHTTPTimeout().Duration()
}
