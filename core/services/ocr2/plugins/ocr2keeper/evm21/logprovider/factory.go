package logprovider

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func New(lggr logger.Logger, poller logpoller.LogPoller, utilsABI abi.ABI) (LogEventProvider, LogRecoverer) {
	filterStore := NewUpkeepFilterStore()
	packer := NewLogEventsPacker(utilsABI)
	provider := NewLogProvider(lggr, poller, packer, filterStore, nil)
	recoverer := NewLogRecoverer(lggr, poller, DefaultRecoveryInterval)

	return provider, recoverer
}
