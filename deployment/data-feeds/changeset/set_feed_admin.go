package changeset

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
)

var _ deployment.ChangeSet[types.SetFeedAdminConfig] = SetFeedAdminChangeset

func SetFeedAdminChangeset(env deployment.Environment, c types.SetFeedAdminConfig) (deployment.ChangesetOutput, error) {
	err := ValidateCacheForChain(env, c.ChainSelector, c.CacheAddress)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to validate cache for chain %w", err)
	}

	state, _ := LoadOnchainState(env)
	chain := env.Chains[c.ChainSelector]
	chainState := state.Chains[c.ChainSelector]
	contract := chainState.DataFeedsCache[c.CacheAddress]

	txOpt := chain.DeployerKey
	if c.UseMCMS {
		txOpt = deployment.SimTransactOpts()
	}

	tx, err := contract.SetFeedAdmin(txOpt, c.AdminAddress, c.IsAdmin)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set feed admin %w", err)
	}

	if c.UseMCMS {
		ops := &timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(c.ChainSelector),
			Batch: []mcms.Operation{
				{
					To:    contract.Address(),
					Data:  tx.Data(),
					Value: big.NewInt(0),
				},
			},
		}

		timelocksPerChain := map[uint64]common.Address{
			c.ChainSelector: chainState.Timelock.Address(),
		}
		proposerMCMSes := map[uint64]*gethwrappers.ManyChainMultiSig{
			c.ChainSelector: chainState.ProposerMcm,
		}

		proposal, err := proposalutils.BuildProposalFromBatches(
			timelocksPerChain,
			proposerMCMSes,
			[]timelock.BatchChainOperation{*ops},
			"proposal to set feed admin on cache",
			0,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{*proposal}}, nil
	}
	_, err = chain.Confirm(tx)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
	}

	return deployment.ChangesetOutput{}, nil
}
