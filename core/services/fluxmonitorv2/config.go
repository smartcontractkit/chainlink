package fluxmonitorv2

import (
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
)

// Config defines the Flux Monitor configuration.
type Config struct {
	DefaultHTTPTimeout             time.Duration
	FlagsContractAddress           string
	MinContractPayment             *assets.Link
	EthGasLimit                    uint64
	EthMaxQueuedTransactions       uint64
	FMDefaultTransactionQueueDepth uint32
}

// MinimumPollingInterval returns the minimum duration between polling ticks
func (c *Config) MinimumPollingInterval() time.Duration {
	return c.DefaultHTTPTimeout
}
