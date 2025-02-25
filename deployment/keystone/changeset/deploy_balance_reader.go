package changeset

import (
	"fmt"
	"maps"
	"slices"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[DeployBalanceReaderRequest] = DeployBalanceReader

type DeployBalanceReaderRequest struct {
	ChainSelectors []uint64 // filter to only deploy to these chains; if empty, deploy to all chains
}

// DeployBalanceReader deploys the BalanceReader contract to all chains in the environment
// callers must merge the output addressbook with the existing one
func DeployBalanceReader(env deployment.Environment, cfg DeployBalanceReaderRequest) (deployment.ChangesetOutput, error) {
	lggr := env.Logger
	ab := deployment.NewMemoryAddressBook()
	selectors := cfg.ChainSelectors
	if len(selectors) == 0 {
		selectors = slices.Collect(maps.Keys(env.Chains))
	}
	for _, sel := range selectors {
		chain, ok := env.Chains[sel]
		if !ok {
			return deployment.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", sel)
		}
		lggr.Infow("deploying balancereader", "chainSelector", chain.Selector)
		balanceReaderResp, err := internal.DeployBalanceReader(chain, ab)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to deploy BalanceReader to chain selector %d: %w", chain.Selector, err)
		}
		lggr.Infof("Deployed %s chain selector %d addr %s", balanceReaderResp.Tv.String(), chain.Selector, balanceReaderResp.Address.String())
	}

	return deployment.ChangesetOutput{AddressBook: ab}, nil
}
