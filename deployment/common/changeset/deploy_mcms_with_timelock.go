package changeset

import (
	"fmt"

	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	evminternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (

	// GrantRoleInTimeLock grants proposer, canceller, bypasser, executor, admin roles to the timelock contract with corresponding addresses if the
	// roles are not already set with the same addresses.
	// It creates a proposal if deployer key is not admin of the timelock contract.
	// otherwise it executes the transactions directly.
	// If neither timelock, nor the deployer key is the admin of the timelock contract, it returns an error.
	GrantRoleInTimeLock = cldf.CreateChangeSet(grantRoleLogic, grantRolePreconditions)
)

type GrantRoleInput struct {
	ExistingProposerByChain map[uint64]common.Address // if needed in the future, need to add bypasser and canceller here
	MCMS                    *proposalutils.TimelockConfig
	GasBoostConfigPerChain  map[uint64]cldfproposalutils.GasBoostConfig
}

func grantRolePreconditions(e cldf.Environment, cfg GrantRoleInput) error {
	mcmsState, err := loadMCMSStatePerChainWithQualifier(e, cfg)
	if err != nil {
		return err
	}
	for selector, proposer := range cfg.ExistingProposerByChain {
		if proposer == (common.Address{}) {
			return fmt.Errorf("proposer address not found for chain %d", selector)
		}
		chain, ok := e.BlockChains.EVMChains()[selector]
		if !ok {
			return fmt.Errorf("chain not found for chain %d", selector)
		}
		timelockContracts, ok := mcmsState[selector]
		if !ok {
			return fmt.Errorf("timelock state not found for chain %d", selector)
		}
		if timelockContracts.Timelock == nil {
			return fmt.Errorf("timelock contract not found for chain %s", chain.String())
		}
		if timelockContracts.ProposerMcm == nil {
			return fmt.Errorf("proposerMcm contract not found for chain %s", chain.String())
		}
		if timelockContracts.CancellerMcm == nil {
			return fmt.Errorf("cancellerMcm contract not found for chain %s", chain.String())
		}
		if timelockContracts.BypasserMcm == nil {
			return fmt.Errorf("bypasserMcm contract not found for chain %s", chain.String())
		}
		if timelockContracts.CallProxy == nil {
			return fmt.Errorf("callProxy contract not found for chain %s", chain.String())
		}
	}
	return nil
}

// loads MCMS state for each chain using per-chain qualifiers from cfg.MCMS.TimelockQualifierPerChain when available
func loadMCMSStatePerChainWithQualifier(e cldf.Environment, cfg GrantRoleInput) (map[uint64]*evmstate.MCMSWithTimelockState, error) {
	result := make(map[uint64]*evmstate.MCMSWithTimelockState)
	for selector := range cfg.ExistingProposerByChain {
		qualifier := ""
		if cfg.MCMS != nil && cfg.MCMS.TimelockQualifierPerChain != nil {
			qualifier = cfg.MCMS.TimelockQualifierPerChain[selector]
		}
		chainState, err := evmstate.MaybeLoadMCMSWithTimelockStateWithQualifier(e, []uint64{selector}, qualifier)
		if err != nil {
			return nil, err
		}
		result[selector] = chainState[selector]
	}
	return result, nil
}

func grantRoleLogic(e cldf.Environment, cfg GrantRoleInput) (cldf.ChangesetOutput, error) {
	mcmsState, err := loadMCMSStatePerChainWithQualifier(e, cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	mcmsStateForProposal := make(map[uint64]evmstate.MCMSWithTimelockState)
	for k, v := range mcmsState {
		if v != nil {
			// Replace the proposer MCM in state with the existing proposer.
			// This is to ensure that we are using an MCM contract that already has the proposer role.
			existingProposerMcm, err := gethwrappers.NewManyChainMultiSig(
				cfg.ExistingProposerByChain[k],
				e.BlockChains.EVMChains()[k].Client,
			)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create ManyChainMultiSig for existing proposer %s on chain %d: %w",
					cfg.ExistingProposerByChain[k].Hex(), k, err)
			}
			mcmsStateForProposal[k] = evmstate.MCMSWithTimelockState{
				CancellerMcm: v.CancellerMcm,
				BypasserMcm:  v.BypasserMcm,
				ProposerMcm:  existingProposerMcm,
				Timelock:     v.Timelock,
				CallProxy:    v.CallProxy,
			}
		}
	}

	out := cldf.ChangesetOutput{}
	gasBoostConfigs := opsutils.GasBoostConfigsForChainMap(cfg.ExistingProposerByChain, cfg.GasBoostConfigPerChain)
	for chain := range cfg.ExistingProposerByChain {
		stateForChain := mcmsState[chain]
		evmChains := e.BlockChains.EVMChains()
		seqReport, err := evminternal.GrantRolesForTimelock(
			e, evmChains[chain], &cldfproposalutils.MCMSWithTimelockContracts{
				CancellerMcm: stateForChain.CancellerMcm,
				BypasserMcm:  stateForChain.BypasserMcm,
				ProposerMcm:  stateForChain.ProposerMcm,
				Timelock:     stateForChain.Timelock,
				CallProxy:    stateForChain.CallProxy,
			}, false, gasBoostConfigs[chain])
		out, err = opsutils.AddEVMCallSequenceToCSOutput(e, out, seqReport, err, mcmsStateForProposal, cfg.MCMS, fmt.Sprintf("GrantRolesForTimelock on %s", evmChains[chain]))
		if err != nil {
			return out, fmt.Errorf("failed to grant roles for timelock on chain %d: %w", chain, err)
		}
	}

	return out, nil
}
