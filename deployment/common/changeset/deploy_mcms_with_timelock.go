package changeset

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmslib "github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmsevmsdk "github.com/smartcontractkit/mcms/sdk/evm"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"golang.org/x/exp/maps"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"

	evminternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/internal/evm"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

var (
	_ cldf.ChangeSet[map[uint64]types.MCMSWithTimelockConfigV2] = DeployMCMSWithTimelockV2

	// GrantRoleInTimeLock grants proposer, canceller, bypasser, executor, admin roles to the timelock contract with corresponding addresses if the
	// roles are not already set with the same addresses.
	// It creates a proposal if deployer key is not admin of the timelock contract.
	// otherwise it executes the transactions directly.
	// If neither timelock, nor the deployer key is the admin of the timelock contract, it returns an error.
	GrantRoleInTimeLock = cldf.CreateChangeSet(grantRoleLogic, grantRolePreconditions)
)

// DeployMCMSWithTimelockV2 deploys and initializes the MCM and Timelock contracts
func DeployMCMSWithTimelockV2(
	env cldf.Environment, cfgByChain map[uint64]types.MCMSWithTimelockConfigV2,
) (cldf.ChangesetOutput, error) {
	newAddresses := cldf.NewMemoryAddressBook()

	for chainSel, cfg := range cfgByChain {
		family, err := chain_selectors.GetSelectorFamily(chainSel)
		if err != nil {
			return cldf.ChangesetOutput{AddressBook: newAddresses}, err
		}

		switch family {
		case chain_selectors.FamilyEVM:
			// load mcms state
			// we load the state one by one to void early return from MaybeLoadMCMSWithTimelockState
			// due to one of the chain not found
			var chainstate *state.MCMSWithTimelockState
			s, err := state.MaybeLoadMCMSWithTimelockState(env, []uint64{chainSel})
			if err != nil {
				// if the state is not found for chain, we assume it's a fresh deployment
				if !strings.Contains(err.Error(), cldf.ErrChainNotFound.Error()) {
					return cldf.ChangesetOutput{}, err
				}
			}
			if s != nil {
				chainstate = s[chainSel]
			}
			_, err = evminternal.DeployMCMSWithTimelockContractsEVM(env.GetContext(), env.Logger, env.BlockChains.EVMChains()[chainSel], newAddresses, cfg, chainstate)
			if err != nil {
				return cldf.ChangesetOutput{AddressBook: newAddresses}, err
			}

		case chain_selectors.FamilySolana:
			// this is not used in CLD as we need to dynamically resolve the artifacts to deploy these contracts
			// we did not want to add the artifact resolution logic here, so we instead deploy using ccip/changeset/solana/cs_deploy_chain.go
			// for in memory tests, programs and state are pre-loaded, so we use this function via testhelpers.TransferOwnershipSolana
			_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, env.BlockChains.SolanaChains()[chainSel], newAddresses, cfg)
			if err != nil {
				return cldf.ChangesetOutput{AddressBook: newAddresses}, err
			}

		default:
			err = fmt.Errorf("unsupported chain family: %s", family)
			return cldf.ChangesetOutput{AddressBook: newAddresses}, err
		}
	}
	ds, err := deployment.MigrateAddressBook(newAddresses)
	if err != nil {
		return cldf.ChangesetOutput{AddressBook: newAddresses}, fmt.Errorf("failed to migrate address book to data store: %w", err)
	}
	return cldf.ChangesetOutput{AddressBook: newAddresses, DataStore: ds}, nil
}

type GrantRoleInput struct {
	ExistingProposerByChain map[uint64]common.Address // if needed in the future, need to add bypasser and canceller here
	MCMS                    *proposalutils.TimelockConfig
}

func grantRolePreconditions(e cldf.Environment, cfg GrantRoleInput) error {
	mcmsState, err := state.MaybeLoadMCMSWithTimelockState(e, maps.Keys(cfg.ExistingProposerByChain))
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

func grantRoleLogic(e cldf.Environment, cfg GrantRoleInput) (cldf.ChangesetOutput, error) {
	mcmsState, err := state.MaybeLoadMCMSWithTimelockState(e, maps.Keys(cfg.ExistingProposerByChain))
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	timelocks := make(map[uint64]string)
	proposers := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)
	batches := make([]mcmstypes.BatchOperation, 0)
	for chain, existingProposer := range cfg.ExistingProposerByChain {
		stateForChain := mcmsState[chain]
		evmChains := e.BlockChains.EVMChains()
		mcmsTxs, err := evminternal.GrantRolesForTimelock(
			e.GetContext(),
			e.Logger, evmChains[chain], &proposalutils.MCMSWithTimelockContracts{
				CancellerMcm: stateForChain.CancellerMcm,
				BypasserMcm:  stateForChain.BypasserMcm,
				ProposerMcm:  stateForChain.ProposerMcm,
				Timelock:     stateForChain.Timelock,
				CallProxy:    stateForChain.CallProxy,
			}, false)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		if len(mcmsTxs) == 0 {
			continue
		}
		timelocks[chain] = mcmsState[chain].Timelock.Address().Hex()
		proposers[chain] = existingProposer.Hex()
		inspectors[chain] = mcmsevmsdk.NewInspector(evmChains[chain].Client)
		batches = append(batches, mcmstypes.BatchOperation{
			ChainSelector: mcmstypes.ChainSelector(chain),
			Transactions:  mcmsTxs,
		})
	}
	// If there are no batches, it means that deployerkey is the admin of timelock, and it has already performed the role grant
	// as part of the deployment. In this case, we don't need to create a proposal.
	if len(batches) == 0 {
		return cldf.ChangesetOutput{}, nil
	}
	if cfg.MCMS == nil {
		return cldf.ChangesetOutput{}, errors.New("MCMS config is nil, but the deployer key is not the admin of the timelock")
	}
	prop, err := proposalutils.BuildProposalFromBatchesV2(
		e,
		timelocks,
		proposers,
		inspectors,
		batches,
		"Grant roles to timelock contracts",
		*cfg.MCMS,
	)
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcmslib.TimelockProposal{*prop},
	}, err
}

func ValidateOwnership(ctx context.Context, mcms bool, deployerKey, timelock common.Address, contract Ownable) error {
	owner, err := contract.Owner(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("failed to get owner: %w", err)
	}
	if mcms && owner != timelock {
		return fmt.Errorf("%s not owned by timelock", contract.Address())
	} else if !mcms && owner != deployerKey {
		return fmt.Errorf("%s not owned by deployer key", contract.Address())
	}
	return nil
}

func ValidateOwnershipSolanaCommon(mcms bool, deployerKey solana.PublicKey, timelockSignerPDA solana.PublicKey, programOwner solana.PublicKey) error {
	if !mcms {
		if deployerKey.String() != programOwner.String() {
			return fmt.Errorf("deployer key %s does not match owner %s", deployerKey.String(), programOwner.String())
		}
	} else {
		if timelockSignerPDA.String() != programOwner.String() {
			return fmt.Errorf("timelock signer PDA %s does not match owner %s", timelockSignerPDA.String(), programOwner.String())
		}
	}
	return nil
}
