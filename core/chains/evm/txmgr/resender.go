package txmgr

import (
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewEvnResender creates a new concrete EvmResender
func NewEvmResender(
	lggr logger.Logger,
	txStore EvmTxStore,
	evmClient EvmTxmClient,
	ks EvmKeyStore,
	pollInterval time.Duration,
	config EvmResenderConfig,
) *EvmResender {
	return txmgr.NewResender(lggr, txStore, evmClient, ks, pollInterval, config)
}
