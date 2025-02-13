package v1_5

import (
	"context"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/mcms"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/proposal/timelock"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_5_0/rmn_contract"
)

var _ deployment.ChangeSet[PermaBlessCommitStoreConfig] = PermaBlessCommitStoreChangeset

type PermaBlessConfigPerSourceChain struct {
	SourceChainSelector uint64
	PermaBless          bool // if true, the commit store will be included in adds and if false it will be included in removes list,
	// https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/RMN.sol#L699C30-L699C54
}

func (p PermaBlessConfigPerSourceChain) Validate(destChain uint64, state changeset.CCIPOnChainState, permaBlessedCommitStores []common.Address) error {
	if err := deployment.IsValidChainSelector(p.SourceChainSelector); err != nil {
		return fmt.Errorf("invalid SourceChainSelector: %w", err)
	}
	_, ok := state.Chains[p.SourceChainSelector]
	if !ok {
		return fmt.Errorf("source chain state not found for chain selector %d", p.SourceChainSelector)
	}
	destState := state.Chains[destChain]
	if destState.CommitStore[p.SourceChainSelector] == nil {
		return fmt.Errorf("dest chain %d does not have a commit store for source chain %d", destChain, p.SourceChainSelector)
	}

	if p.PermaBless {
		if slices.Contains(permaBlessedCommitStores, destState.CommitStore[p.SourceChainSelector].Address()) {
			return fmt.Errorf("commit store for source chain %d is already permablessed", p.SourceChainSelector)
		}
	} else {
		if !slices.Contains(permaBlessedCommitStores, destState.CommitStore[p.SourceChainSelector].Address()) {
			return fmt.Errorf("commit store for source chain %d is not permablessed, cannot be removed", p.SourceChainSelector)
		}
	}
	return nil
}

type PermaBlessCommitStoreConfigPerDest struct {
	Sources []PermaBlessConfigPerSourceChain
}

type PermaBlessCommitStoreConfig struct {
	Configs    map[uint64]PermaBlessCommitStoreConfigPerDest
	MCMSConfig *changeset.MCMSConfig
}

func (c PermaBlessCommitStoreConfig) Validate(env deployment.Environment) error {
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for destChain, pCfg := range c.Configs {
		if err := deployment.IsValidChainSelector(destChain); err != nil {
			return fmt.Errorf("invalid DestChainSelector: %w", err)
		}
		destState, ok := state.Chains[destChain]
		if !ok {
			return fmt.Errorf("dest chain state not found for chain selector %d", destChain)
		}
		if destState.RMN == nil {
			return fmt.Errorf("dest chain %d does not have an RMN", destChain)
		}
		if destState.CommitStore == nil {
			return fmt.Errorf("dest chain %d does not have any commit store", destChain)
		}
		// get all permablessed commit stores
		permaBlessedCommitStores, err := destState.RMN.GetPermaBlessedCommitStores(nil)
		if err != nil {
			return fmt.Errorf("failed to get perma blessed commit stores: %w", err)
		}
		for _, sourceCfg := range pCfg.Sources {
			if err := sourceCfg.Validate(destChain, state, permaBlessedCommitStores); err != nil {
				return fmt.Errorf("invalid PermaBlessConfig for source chain %d and dest chain %d : %w", sourceCfg.SourceChainSelector, destChain, err)
			}
		}

		if err := commoncs.ValidateOwnership(context.Background(), c.MCMSConfig != nil, env.Chains[destChain].DeployerKey.From, destState.Timelock.Address(), destState.RMN); err != nil {
			return fmt.Errorf("failed to validate ownership: %w", err)
		}
	}
	return nil
}

// PermaBlessCommitStoreChangeset permablesses the commit stores on the RMN contract
// If commit store addresses are added to the permaBlessed list, those will be considered automatically blessed.
// This changeset can add to or remove from the existing permaBlessed list.
func PermaBlessCommitStoreChangeset(env deployment.Environment, c PermaBlessCommitStoreConfig) (deployment.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid PermaBlessCommitStoreConfig: %w", err)
	}
	state, err := changeset.LoadOnchainState(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	ops := make([]timelock.BatchChainOperation, 0)
	timelocks := make(map[uint64]common.Address)
	proposerMCM := make(map[uint64]*gethwrappers.ManyChainMultiSig)
	for destChain, cfg := range c.Configs {
		destState := state.Chains[destChain]
		RMN := destState.RMN
		var removes []common.Address
		var adds []common.Address
		for _, sourceCfg := range cfg.Sources {
			commitStore := destState.CommitStore[sourceCfg.SourceChainSelector]
			if sourceCfg.PermaBless {
				adds = append(adds, commitStore.Address())
			} else {
				removes = append(removes, commitStore.Address())
			}
		}
		txOpts := env.Chains[destChain].DeployerKey
		if c.MCMSConfig != nil {
			txOpts = deployment.SimTransactOpts()
		}
		tx, err := RMN.OwnerRemoveThenAddPermaBlessedCommitStores(txOpts, removes, adds)
		if c.MCMSConfig == nil {
			_, err = deployment.ConfirmIfNoErrorWithABI(env.Chains[destChain], tx, rmn_contract.RMNContractABI, err)
			if err != nil {
				return deployment.ChangesetOutput{}, err
			}
			env.Logger.Infof("PermaBlessed commit stores on chain %d removed %v, added %v", destChain, removes, adds)
			continue
		}
		timelocks[destChain] = destState.Timelock.Address()
		proposerMCM[destChain] = destState.ProposerMcm
		ops = append(ops, timelock.BatchChainOperation{
			ChainIdentifier: mcms.ChainIdentifier(destChain),
			Batch: []mcms.Operation{
				{
					To:    RMN.Address(),
					Data:  tx.Data(),
					Value: big.NewInt(0),
				},
			},
		})
	}
	if c.MCMSConfig == nil {
		return deployment.ChangesetOutput{}, nil
	}
	p, err := proposalutils.BuildProposalFromBatches(
		timelocks,
		proposerMCM,
		ops,
		"PermaBless commit stores on RMN",
		c.MCMSConfig.MinDelay,
	)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	env.Logger.Infof("perma bless commit stores proposal created with %d operations", len(ops))
	return deployment.ChangesetOutput{Proposals: []timelock.MCMSWithTimelockProposal{
		*p,
	}}, nil
}
