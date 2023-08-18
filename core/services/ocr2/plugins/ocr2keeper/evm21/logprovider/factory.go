package logprovider

import (
	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21/core"
)

func New(lggr logger.Logger, poller logpoller.LogPoller, c client.Client, utilsABI abi.ABI, stateStore core.UpkeepStateReader) (LogEventProvider, LogRecoverer) {
	filterStore := NewUpkeepFilterStore()
	packer := NewLogEventsPacker(utilsABI)
	provider := NewLogProvider(lggr, poller, packer, filterStore, nil)
	recoverer := NewLogRecoverer(lggr, poller, c, stateStore, packer, filterStore, 0, provider.opts.LookbackBlocks)

	return provider, recoverer
}
