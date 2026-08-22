package ccvcommon

import (
	"context"
	"fmt"
	"slices"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccv/protocol"
	common "github.com/smartcontractkit/chainlink-common/pkg/logger"
	commontypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"
)

// GetLegacyChains maps the chain services of the node to the chain selectors in the job
// configuration. A chain service that the node cannot use is skipped, not an error, so that one
// bad chain does not stop the whole job.
//
// The second return value holds the chain selectors from chainsInConfig that have no chain object.
// Give this list to NewMissingChainsMonitor so that operators get a repeated alarm.
//
// An error occurs only when no chain in the configuration has a chain object, because a job with
// no chains cannot do work.
func GetLegacyChains(
	ctx context.Context,
	lggr common.Logger,
	chainServices []commontypes.ChainService,
	chainsInConfig []protocol.ChainSelector,
) (map[protocol.ChainSelector]legacyevm.Chain, []protocol.ChainSelector, error) {
	chains := make(map[protocol.ChainSelector]legacyevm.Chain)
	for _, c := range chainServices {
		chainInfo, err := c.GetChainInfo(ctx)
		if err != nil {
			lggr.Warnw("skipping chain service: failed to get chain info", "chain", c.Name(), "err", err)
			continue
		}

		chain, ok := c.(legacyevm.Chain)
		if !ok {
			lggr.Warnw("skipping chain service: failed to cast to legacyevm.Chain, LOOPP mode is currently not supported",
				"chain", c.Name(), "chainInfo", chainInfo)
			continue
		}

		id := chain.ID()

		// convert to selector
		chain2, ok := chainselectors.ChainByEvmChainID(id.Uint64())
		if !ok {
			lggr.Warnw("skipping chain service: failed to get chain selector", "chain", c.Name(), "chainID", id.String())
			continue
		}

		if !slices.Contains(chainsInConfig, protocol.ChainSelector(chain2.Selector)) {
			lggr.Infow("skipping chain not in config", "chain", chain2.Selector, "chainID", id.String())
			continue
		}

		chains[protocol.ChainSelector(chain2.Selector)] = chain
	}

	missing := MissingChains(chainsInConfig, chains)

	if len(chains) == 0 && len(chainsInConfig) > 0 {
		return nil, nil, fmt.Errorf("no chain object is available for any of the %d chains in the configuration: %v",
			len(chainsInConfig), chainsInConfig)
	}

	return chains, missing, nil
}

// MissingChains returns the chain selectors in chainsInConfig that are not keys of chains.
func MissingChains[V any](chainsInConfig []protocol.ChainSelector, chains map[protocol.ChainSelector]V) []protocol.ChainSelector {
	var missing []protocol.ChainSelector
	for _, sel := range chainsInConfig {
		if _, ok := chains[sel]; !ok {
			missing = append(missing, sel)
		}
	}
	slices.Sort(missing)
	return slices.Compact(missing)
}
