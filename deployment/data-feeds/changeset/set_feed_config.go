package changeset

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset/types"
	"math/big"
)

var _ deployment.ChangeSet[types.SetFeedDecimalConfig] = SetFeedConfigChangeset

func SetFeedConfigChangeset(env deployment.Environment, c types.SetFeedDecimalConfig) (deployment.ChangesetOutput, error) {
	state, err := LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load on chain state %w", err)
	}
	chain, ok := env.Chains[c.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in environment")
	}
	chainState, ok := state.Chains[c.ChainSelector]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("chain not found in on chain state")
	}
	if chainState.DataFeedsCache == nil {
		return deployment.ChangesetOutput{}, errors.New("DataFeedsCache not found in on chain state")
	}
	contract, ok := chainState.DataFeedsCache[c.CacheAddress]
	if !ok {
		return deployment.ChangesetOutput{}, errors.New("contract not found in on chain state")
	}

	txOpt := chain.DeployerKey
	if c.UseMCMS {
		txOpt = deployment.SimTransactOpts()
	}

	tx, err := contract.SetDecimalFeedConfigs(txOpt, c.DataIds, c.Descriptions, c.WorkflowMetadata)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to set feed config %w", err)
	}

	if !c.UseMCMS {
		_, err = chain.Confirm(tx)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to confirm transaction: %s, %w", tx.Hash().String(), err)
		}
	} else {
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
			"proposal to set feed config on cache",
			0,
		)
		if err != nil {
			return deployment.ChangesetOutput{}, fmt.Errorf("failed to build proposal: %w", err)
		}
		return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{*proposal}}, nil
	}

	return deployment.ChangesetOutput{}, nil
}
