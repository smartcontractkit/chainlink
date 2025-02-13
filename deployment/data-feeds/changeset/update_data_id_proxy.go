package changeset

import (
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var _ deployment.ChangeSet[types.UpdateDataIDProxyConfig] = UpdateDataIDProxyChangeset

func UpdateDataIDProxyChangeset(env deployment.Environment, c types.UpdateDataIDProxyConfig) (deployment.ChangesetOutput, error) {
	if len(c.DataIDs) != len(c.Proxies) {
		return deployment.ChangesetOutput{}, errors.New("dataIds and proxies length mismatch")
	}
	err := ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate cache for chain %w", err)
	}

	state, _ := LoadOnchainState(env)
	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.McmsConfig != nil {
		txOpt = deployment.SimTransactOpts()
	}

	tx, err := contract.UpdateDataIdMappingsForProxies(txOpt, c.Proxies, c.DataIDs)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set proxy-dataId mapping %w", err)
	}

	if c.McmsConfig != nil {
		proposal, err := BuildMCMProposal(env, "proposal to update proxy-dataId mapping on a cache", c.ChainSelector, contract.Address().Hex(), tx, c.McmsConfig.MinDelay)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{*proposal}}, nil
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{}, nil
}
