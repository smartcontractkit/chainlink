package v1_5

import (
	"context"
	"fmt"
	"math/big"
	"slices"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/rmn_contract"

	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	commoncs "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var _ cldf.ChangeSet[PermaBlessCommitStoreConfig] = PermaBlessCommitStoreChangeset

type PermaBlessConfigPerSourceChain struct {
	SourceChainSelector uint64
	PermaBless          bool // if true, the commit store will be included in adds and if false it will be included in removes list,
	// https://github.com/smartcontractkit/ccip/blob/ccip-develop/contracts/src/v0.8/ccip/RMN.sol#L699C30-L699C54
}

func (p PermaBlessConfigPerSourceChain) Validate(destChain uint64, state stateview.CCIPOnChainState, permaBlessedCommitStores []common.Address) error {
	if err := cldf.IsValidChainSelector(p.SourceChainSelector); err != nil {
		return fmt.Errorf("invalid SourceChainSelector: %w", err)
	}
	_, ok := state.EVMChainState(p.SourceChainSelector)
	if !ok {
		return fmt.Errorf("source chain state not found for chain selector %d", p.SourceChainSelector)
	}
	destState := state.MustGetEVMChainState(destChain)
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
	MCMSConfig *proposalutils.TimelockConfig
}

func (c PermaBlessCommitStoreConfig) Validate(env cldf.Environment) error {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}
	for destChain, pCfg := range c.Configs {
		if err := cldf.IsValidChainSelector(destChain); err != nil {
			return fmt.Errorf("invalid DestChainSelector: %w", err)
		}
		destState, ok := state.EVMChainState(destChain)
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

		if err := commoncs.ValidateOwnership(context.Background(), c.MCMSConfig != nil, env.BlockChains.EVMChains()[destChain].DeployerKey.From, destState.Timelock.Address(), destState.RMN); err != nil {
			return fmt.Errorf("failed to validate ownership: %w", err)
		}
	}
	return nil
}

// PermaBlessCommitStoreChangeset permablesses the commit stores on the RMN contract
// If commit store addresses are added to the permaBlessed list, those will be considered automatically blessed.
// This changeset can add to or remove from the existing permaBlessed list.
func PermaBlessCommitStoreChangeset(env cldf.Environment, c PermaBlessCommitStoreConfig) (cldf.ChangesetOutput, error) {
	if err := c.Validate(env); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("invalid PermaBlessCommitStoreConfig: %w", err)
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	ops := make([]mcmstypes.BatchOperation, 0)
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for destChain, cfg := range c.Configs {
		destState := state.MustGetEVMChainState(destChain)
		RMN := destState.RMN

		var removes, adds []common.Address
		for _, sourceCfg := range cfg.Sources {
			commitStore := destState.CommitStore[sourceCfg.SourceChainSelector]
			if sourceCfg.PermaBless {
				adds = append(adds, commitStore.Address())
			} else {
				removes = append(removes, commitStore.Address())
			}
		}

		txOpts := env.BlockChains.EVMChains()[destChain].DeployerKey
		if c.MCMSConfig != nil {
			txOpts = cldf.SimTransactOpts()
		}
		tx, err := RMN.OwnerRemoveThenAddPermaBlessedCommitStores(txOpts, removes, adds)

		// note: error check is handled below
		if c.MCMSConfig == nil {
			_, err = cldf.ConfirmIfNoErrorWithABI(env.BlockChains.EVMChains()[destChain], tx, rmn_contract.RMNContractABI, err)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
			env.Logger.Infof("PermaBlessed commit stores on chain %d removed %v, added %v", destChain, removes, adds)
			continue
		} else if err != nil {
			return cldf.ChangesetOutput{}, err
		}

		timelocks[destChain] = destState.Timelock.Address().Hex()

		inspectors[destChain], err = proposalutils.McmsInspectorForChain(env, destChain)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %d: %w", destChain, err)
		}

		batchOperation, err := proposalutils.BatchOperationForChain(destChain, RMN.Address().Hex(), tx.Data(), big.NewInt(0),
			string(shared.RMN), []string{})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create batch operation for chain %d: %w", destChain, err)
		}

		ops = append(ops, batchOperation)
	}
	if c.MCMSConfig == nil {
		return cldf.ChangesetOutput{}, nil
	}

	mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(env, state, c.MCMSConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build mcm addresses per chain: %w", err)
	}
	timelockProposal, err := proposalutils.BuildProposalFromBatchesV2(
		env,
		timelocks,
		mcmsContractByChain,
		inspectors,
		ops,
		"PermaBless commit stores on RMN",
		*c.MCMSConfig,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	env.Logger.Infof("perma bless commit stores proposal created with %d operations", len(ops))
	return cldf.ChangesetOutput{MCMSTimelockProposals: []mcmslib.TimelockProposal{
		*timelockProposal,
	}}, nil
}
