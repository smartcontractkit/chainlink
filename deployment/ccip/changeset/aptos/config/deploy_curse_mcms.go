package config

import (
	"fmt"

	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// DeployCurseMCMSConfig holds per-chain configuration for deploying and
// configuring a CurseMCMS contract on Aptos chains.
type DeployCurseMCMSConfig struct {
	CurseMCMSConfigPerChain    map[uint64]types.MCMSWithTimelockConfigV2
	MCMSTimelockConfigPerChain map[uint64]proposalutils.TimelockConfig
}

func (c DeployCurseMCMSConfig) Validate() error {
	for cs, cfg := range c.CurseMCMSConfigPerChain {
		if err := cldf.IsValidChainSelector(cs); err != nil {
			return fmt.Errorf("invalid chain selector: %d - %w", cs, err)
		}
		for _, mcmsCfg := range []mcmstypes.Config{cfg.Bypasser, cfg.Canceller, cfg.Proposer} {
			if err := mcmsCfg.Validate(); err != nil {
				return fmt.Errorf("invalid MCMS config for chain %d: %w", cs, err)
			}
		}
	}
	return nil
}
