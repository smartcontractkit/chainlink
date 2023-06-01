package txmgr

import (
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// NewEvmBroadcaster returns a new concrete EvmBroadcaster
func NewEvmBroadcaster(
	txStore EvmTxStore,
	evmClient EvmTxmClient,
	config txmgrtypes.BroadcasterConfig[*assets.Wei],
	keystore EvmKeyStore,
	eventBroadcaster pg.EventBroadcaster,
	txAttemptBuilder EvmTxAttemptBuilder,
	nonceSyncer EvmNonceSyncer,
	logger logger.Logger,
	checkerFactory EvmTransmitCheckerFactory,
	autoSyncNonce bool,
) *EvmBroadcaster {
	return txmgr.NewBroadcaster(txStore, evmClient, config, keystore, eventBroadcaster, txAttemptBuilder, nonceSyncer, logger, checkerFactory, autoSyncNonce, stringToGethAddress)
}
