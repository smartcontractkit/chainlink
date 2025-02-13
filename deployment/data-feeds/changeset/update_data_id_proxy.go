package changeset

import (
	"encoding/json"
	"errors"
	"fmt"

	mcmslib "github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/sdk"
	"github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
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
		ops := &mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(c.ChainSelector),
			Transactions: []mcmstypes.Transaction{
				{
					To:               contract.Address().Hex(),
					Data:             tx.Data(),
					AdditionalFields: json.RawMessage(`{"value": 0}`),
				},
			},
		}

		timelocksPerChain := map[uint64]string{
			c.ChainSelector: chainState.Timelock.Address().Hex(),
		}
		proposerMCMSes := map[uint64]string{
			c.ChainSelector: chainState.ProposerMcm.Address().Hex(),
		}

		inspectorPerChain := map[uint64]sdk.Inspector{}
		inspectorPerChain[c.ChainSelector] = evm.NewInspector(chain.Client)

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env.GetContext(),
			timelocksPerChain,
			proposerMCMSes,
			inspectorPerChain,
			[]mcmstypes.BatchOperation{*ops},
			"proposal to update proxy-dataId mapping on a cache",
			c.McmsConfig.MinDelay,
		)
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
