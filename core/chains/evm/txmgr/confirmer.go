package txmgr

import (
	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/txmgr"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewEvmConfirmer instantiates a new EVM confirmer
func NewEvmConfirmer(
	txStore EvmTxStore,
	evmClient EvmTxmClient,
	config txmgrtypes.ConfirmerConfig[*assets.Wei],
	keystore EvmKeyStore,
	txAttemptBuilder EvmTxAttemptBuilder,
	lggr logger.Logger,
) *EvmConfirmer {
	return txmgr.NewConfirmer(txStore, evmClient, config, keystore, txAttemptBuilder, lggr, func(r *evmtypes.Receipt) bool { return r == nil })
}
