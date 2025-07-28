package crib

import (
	"fmt"
)

type TierConfigs struct {
	NumChains int `json:"NumChains"`
	NumNodes  int `json:"NumNodes"`
}

type ChainTiers struct {
	Tiers []TierConfigs `json:"tiers"`
}

func (ct *ChainTiers) Validate(numChains int) error {
	sumChainsInTiers := 0
	for _, tier := range ct.Tiers {
		if tier.NumChains <= 0 || tier.NumNodes <= 0 {
			return fmt.Errorf("each tier must have positive values for NumChains, NumNodes got: %v", tier)
		}
		sumChainsInTiers += tier.NumChains
	}

	if sumChainsInTiers != numChains {
		return fmt.Errorf("the number of chains in all tiers (%d) must match the number of chains (%d)", sumChainsInTiers, numChains)
	}
	return nil
}
