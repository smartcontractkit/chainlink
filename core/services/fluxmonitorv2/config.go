package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

// Config defines the Flux Monitor configuration.
type Config interface {
	FlagsContractAddress() string           // Evm
	MinimumContractPayment() *assets.Link   // Evm
	EvmGasLimitDefault() uint32             // Evm
	EvmGasLimitFMJobType() *uint32          // Evm
	EvmMaxQueuedTransactions() uint64       // Evm
	FMDefaultTransactionQueueDepth() uint32 // FluxMonitor
}

type JobPipelineConfig interface {
	DefaultHTTPTimeout() models.Duration
}

// MinimumPollingInterval returns the minimum duration between polling ticks
func MinimumPollingInterval(c JobPipelineConfig) time.Duration {
	return c.DefaultHTTPTimeout().Duration()
}
