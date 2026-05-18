package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/assets"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config"
)

// Config defines the Flux Monitor configuration.
type Config interface {
	FlagsContractAddress() string     // Evm
	MinContractPayment() *assets.Link // Evm
}

type EvmFeeConfig interface {
	LimitDefault() uint64 // Evm
	LimitJobType() config.LimitJobType
}

type EvmTransactionsConfig interface {
	MaxQueued() uint64 // Evm
}

type FluxMonitorConfig interface {
	DefaultTransactionQueueDepth() uint32
}

type JobPipelineConfig interface {
	DefaultHTTPTimeout() commonconfig.Duration
}

// MinimumPollingInterval returns the minimum duration between polling ticks
func MinimumPollingInterval(c JobPipelineConfig) time.Duration {
	return c.DefaultHTTPTimeout().Duration()
}
