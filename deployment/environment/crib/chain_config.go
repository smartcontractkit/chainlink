package crib

import (
	"fmt"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

type TierConfigs struct {
	NumChains int
	NumNodes  int
	NumLanes  int
}

type ChainTiers struct {
	Tiers []TierConfigs
}

func (ct *ChainTiers) Validate(e *deployment.Environment) error {
	sumChainsInTiers := 0
	for _, tier := range ct.Tiers {
		if tier.NumChains <= 0 || tier.NumNodes <= 0 || tier.NumLanes <= 0 {
			return fmt.Errorf("each tier must have positive values for NumChains, NumNodes, and NumLanes, got: %v", tier)
		}
		sumChainsInTiers += tier.NumChains
	}

	allSelectors := e.BlockChains.ListChainSelectors()
	if sumChainsInTiers != len(allSelectors) {
		return fmt.Errorf("the sum of chains in all tiers (%d) must match the number of chains in the environment (%d)", sumChainsInTiers, len(allSelectors))
	}
	return nil
}
