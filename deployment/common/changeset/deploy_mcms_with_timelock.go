package changeset

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	xerrgroup "golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gagliardetto/solana-go"
	"github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	evminternal "github.com/smartcontractkit/chainlink/deployment/common/changeset/evm/mcms"
	solanaMCMS "github.com/smartcontractkit/chainlink/deployment/common/changeset/solana/mcms"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/opsutils"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// migrateAddressBookWithQualifiers migrates an address book to a data store,
// applying custom qualifiers from MCMS configs when available
func migrateAddressBookWithQualifiers(ab cldf.AddressBook, cfgByChain map[uint64]types.MCMSWithTimelockConfigV2) (datastore.MutableDataStore, error) {
	addrs, err := ab.Addresses()
	if err != nil {
		return nil, err
	}

	ds := datastore.NewMemoryDataStore()

	for chainSelector, chainAddresses := range addrs {
		// Get the qualifier for this chain from the config
		qualifier := ""
		if cfg, exists := cfgByChain[chainSelector]; exists && cfg.Qualifier != nil && *cfg.Qualifier != "" {
			qualifier = *cfg.Qualifier
		}

		for addr, typever := range chainAddresses {
			ref := datastore.AddressRef{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          datastore.ContractType(typever.Type),
				Version:       &typever.Version,
			}

			// If we have a custom qualifier for this chain, use it for MCMS contracts
			if qualifier != "" && isMCMSContract(string(typever.Type)) {
				ref.Qualifier = qualifier
			} else {
				// Use the original auto-generated qualifier for other contracts
				ref.Qualifier = fmt.Sprintf("%s-%s", addr, typever.Type)
			}

			// If the address book has labels, we need to add them to the addressRef
			if !typever.Labels.IsEmpty() {
				ref.Labels = datastore.NewLabelSet(typever.Labels.List()...)
			}

			if err = ds.Addresses().Add(ref); err != nil {
				return nil, fmt.Errorf("failed to add address %s: %w", addr, err)
			}
		}
	}
	return ds, nil
}

// isMCMSContract checks if a contract type is part of the MCMS system
func isMCMSContract(contractType string) bool {
	mcmsTypes := []string{
		string(types.RBACTimelock),
		string(types.ManyChainMultisig),
		string(types.ProposerManyChainMultisig),
		string(types.BypasserManyChainMultisig),
		string(types.CancellerManyChainMultisig),
		string(types.CallProxy),
	}

	return slices.Contains(mcmsTypes, contractType)
}

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

	eg := xerrgroup.Group{}
	mu := sync.Mutex{}
	allReports := make([]operations.Report[any, any], 0)
	for chainSel, cfg := range cfgByChain {
		eg.Go(func() error {
			family, err := chain_selectors.GetSelectorFamily(chainSel)
			if err != nil {
				return err
			}

			switch family {
			case chain_selectors.FamilyEVM:
				// Extract qualifier from config for this chain
				qualifier := ""
				if cfg.Qualifier != nil {
					qualifier = *cfg.Qualifier
				}

				// load mcms state with qualifier awareness
				// we load the state one by one to avoid early return from MaybeLoadMCMSWithTimelockStateWithQualifier
				// due to one of the chain not found
				var chainstate *state.MCMSWithTimelockState
				s, err := state.MaybeLoadMCMSWithTimelockStateWithQualifier(env, []uint64{chainSel}, qualifier)
				if err != nil {
					// if the state is not found for chain, we assume it's a fresh deployment
					// this includes "no addresses found" which is expected for new qualifiers
					if !strings.Contains(err.Error(), cldf.ErrChainNotFound.Error()) &&
						!strings.Contains(err.Error(), "no addresses found") {
						return err
					}
				}
				if s != nil {
					chainstate = s[chainSel]
				}
				reports, err := evminternal.DeployMCMSWithTimelockContractsEVM(env, env.BlockChains.EVMChains()[chainSel], newAddresses, cfg, chainstate)
				mu.Lock()
				allReports = append(allReports, reports...)
				mu.Unlock()

				return err

			case chain_selectors.FamilySolana:
				// this is not used in CLD as we need to dynamically resolve the artifacts to deploy these contracts
				// we did not want to add the artifact resolution logic here, so we instead deploy using ccip/changeset/solana/cs_deploy_chain.go
				// for in memory tests, programs and state are pre-loaded, so we use this function via testhelpers.TransferOwnershipSolana
				_, err := solanaMCMS.DeployMCMSWithTimelockProgramsSolana(env, env.BlockChains.SolanaChains()[chainSel], newAddresses, cfg)
				return err

			default:
				return fmt.Errorf("unsupported chain family: %s", family)
			}
		})
	}
	err := eg.Wait()
	if err != nil {
		return cldf.ChangesetOutput{Reports: allReports, AddressBook: newAddresses}, err
	}
	ds, err := migrateAddressBookWithQualifiers(newAddresses, cfgByChain)
	if err != nil {
		return cldf.ChangesetOutput{Reports: allReports, AddressBook: newAddresses}, fmt.Errorf("failed to migrate address book to data store: %w", err)
	}
	return cldf.ChangesetOutput{Reports: allReports, AddressBook: newAddresses, DataStore: ds}, nil
}

type GrantRoleInput struct {
	ExistingProposerByChain map[uint64]common.Address // if needed in the future, need to add bypasser and canceller here
	MCMS                    *proposalutils.TimelockConfig
	GasBoostConfigPerChain  map[uint64]types.GasBoostConfig
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
	mcmsStateForProposal := make(map[uint64]state.MCMSWithTimelockState)
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
			mcmsStateForProposal[k] = state.MCMSWithTimelockState{
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
			e, evmChains[chain], &proposalutils.MCMSWithTimelockContracts{
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
