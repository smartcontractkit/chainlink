package ccvcommon

import (
	"context"
	"fmt"
	"maps"
	"slices"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func GetLegacyChains(ctx context.Context, lggr logger.Logger, chainServices []commontypes.ChainService, chainsInConfig []protocol.ChainSelector) (map[protocol.ChainSelector]legacyevm.Chain, error) {
	chains := make(map[protocol.ChainSelector]legacyevm.Chain)
	for _, c := range chainServices {
		chainInfo, err := c.GetChainInfo(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get chain info for chain %s: %w", c.Name(), err)
		}

		chain, ok := c.(legacyevm.Chain)
		if !ok {
			return nil, fmt.Errorf("failed to cast chain service %s to legacyevm.Chain (info: %+v), LOOPP mode is currently not supported", c.Name(), chainInfo)
		}

		id := chain.ID()

		// convert to selector
		chain2, ok := chainselectors.ChainByEvmChainID(id.Uint64())
		if !ok {
			return nil, fmt.Errorf("failed to get chain selector for chain %s", id.String())
		}

		if !slices.Contains(chainsInConfig, protocol.ChainSelector(chain2.Selector)) {
			lggr.Infow("skipping chain not in config", "chain", chain2.Selector, "chainID", id.String())
			continue
		}

		chains[protocol.ChainSelector(chain2.Selector)] = chain
	}
	// check if we got all the chains in the configuration
	if len(chains) != len(chainsInConfig) {
		// Don't error out here because we still wanna run with whatever we got.
		lggr.Warnw("did not get all the chains in the configuration", "want", chainsInConfig, "got", slices.Collect(maps.Keys(chains)))
	}
	return chains, nil
}
