package mcmsnew

import (
	"math/big"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	mcmsTypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
)

// MCMSWithTimelockConfig holds the configuration for an MCMS with timelock.
// Note that this type already exists in types.go, but this one is using the new lib version.
type MCMSWithTimelockConfig struct {
	Canceller        mcmsTypes.Config
	Bypasser         mcmsTypes.Config
	Proposer         mcmsTypes.Config
	TimelockMinDelay *big.Int
}

func DeployMCMSWithTimelockContractsBatch(
	lggr logger.Logger,
	chains deployment.MultiFamilyChains,
	ab deployment.AddressBook,
	cfgByChain map[uint64]MCMSWithTimelockConfig,
) error {
	for chainSel, cfg := range cfgByChain {
		family, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return err
		}
		switch family {
		case chain_selectors.FamilyEVM:
			_, err := deployMCMSWithTimelockContractsEVM(lggr, chains.EVMChains[chainSel], ab, cfg)
			if err != nil {
				return err
			}
		case chain_selectors.FamilySolana:
			_, err := deployMCMSWithTimelockContractsSolana(lggr, chains.SolChains[chainSel], ab, cfg)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
