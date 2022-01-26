package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

// Config defines the Flux Monitor configuration.
type Config interface {
	DefaultHTTPTimeout() models.Duration
	FlagsContractAddress() string
	MinimumContractPayment() *assets.Link
	EvmGasLimitDefault() uint64
	EvmMaxQueuedTransactions() uint64
	FMDefaultTransactionQueueDepth() uint32
	LogSQL() bool
}

// MinimumPollingInterval returns the minimum duration between polling ticks
func MinimumPollingInterval(c Config) time.Duration {
	return c.DefaultHTTPTimeout().Duration()
}
