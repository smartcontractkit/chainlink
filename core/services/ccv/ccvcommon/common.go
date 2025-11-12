package ccvcommon

import (
	"math"
	"math/big"

	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccv/protocol"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func GetLegacyChains(lggr logger.Logger, chainServices []commontypes.ChainService) map[protocol.ChainSelector]legacyevm.Chain {
	chains := make(map[protocol.ChainSelector]legacyevm.Chain)
	for _, c := range chainServices {
		chain, ok := c.(legacyevm.Chain)
		if !ok {
			lggr.Info("CCV: failed to cast legacyevm.Chain")
			continue
		}

		id := chain.ID()
		if id.Cmp(new(big.Int).SetUint64(math.MaxUint64)) > 0 {
			lggr.Info("CCV: chain ID too large")
			continue
		}

		// convert to selector
		chain2, ok := chainselectors.ChainByEvmChainID(id.Uint64())
		if !ok {
			lggr.Infow("CCV: failed to get chain selector")
			continue
		}

		chains[protocol.ChainSelector(chain2.Selector)] = chain
	}
	return chains
}
