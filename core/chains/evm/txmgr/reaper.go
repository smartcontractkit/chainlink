package txmgr

import (
	"math/big"

	txmgrtypes "github.com/smartcontractkit/chainlink/v2/common/txmgr/types"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

// NewEvmReaper instantiates a new EVM-specific reaper object
func NewEvmReaper(lggr logger.Logger, store txmgrtypes.TxHistoryReaper[*big.Int], config EvmReaperConfig, chainID *big.Int) *EvmReaper {
	return txmgr.NewReaper(lggr, store, config, chainID)
}
