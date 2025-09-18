package ccip_attestation

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	mcmssdk "github.com/smartcontractkit/mcms/sdk"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	signer_registry "github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	EVMSignerRegistrySetNewSignerAddressesChangeset = cldf.CreateChangeSet(signerRegistrySetNewSignerAddressesLogic, signerRegistrySetNewSignerAddressesPrecondition)
)

type SetNewSignerAddressesConfig struct {
	// MCMS defines the delay to use for Timelock (if absent, the changeset will attempt to use the deployer key).
	MCMS *proposalutils.TimelockConfig
	// UpdatesByChain maps chain selector -> (existing signer -> new signer) for per-chain updates.
	UpdatesByChain map[uint64]map[common.Address]common.Address
}

func signerRegistrySetNewSignerAddressesPrecondition(env cldf.Environment, config SetNewSignerAddressesConfig) error {
	if len(config.UpdatesByChain) == 0 {
		return errors.New("no signer updates provided")
	}

	// Per-chain basic validation and duplicate checks
	for chainSelector, updates := range config.UpdatesByChain {
		if len(updates) == 0 {
			return fmt.Errorf("no signer updates provided for chain selector %d", chainSelector)
		}
		seenNew := make(map[common.Address]common.Address)
		for existingAddr, newAddr := range updates {
			if existingAddr == utils.ZeroAddress {
				return errors.New("existing signer address cannot be zero address")
			}
			if newAddr == utils.ZeroAddress {
				return fmt.Errorf("new signer address for %s cannot be zero address", existingAddr.Hex())
			}
			if existingAddr == newAddr {
				return fmt.Errorf("existing address %s and new address are the same", existingAddr.Hex())
			}
			if prevExisting, exists := seenNew[newAddr]; exists {
				return fmt.Errorf("duplicate new address %s for existing signers %s and %s",
					newAddr.Hex(), prevExisting.Hex(), existingAddr.Hex())
			}
			seenNew[newAddr] = existingAddr
		}
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Validate signers exist on each provided chain using the loaded state
	for chainSelector, updates := range config.UpdatesByChain {
		chainState, exists := state.Chains[chainSelector]
		if !exists {
			continue
		}

		if len(chainState.SignerRegistrySigners) == 0 {
			env.Logger.Infof("No signer registry data found on chain selector %d, skipping", chainSelector)
			continue
		}

		existingSigners := make(map[common.Address]bool)
		for _, signer := range chainState.SignerRegistrySigners {
			existingSigners[signer.EvmAddress] = true

			if signer.NewEVMAddress != utils.ZeroAddress {
				existingSigners[signer.NewEVMAddress] = true
			}
		}

		// Check each address we want to update exists in the registry for this chain
		for existingAddr, newAddr := range updates {
			if !existingSigners[existingAddr] {
				return fmt.Errorf("address %s is not a registered signer on chain selector %d", existingAddr.Hex(), chainSelector)
			}
			if newAddr != utils.ZeroAddress {
				if existingSigners[newAddr] {
					return fmt.Errorf("new address %s is already a signer or pending new address on chain selector %d", newAddr.Hex(), chainSelector)
				}
			}
		}
	}

	return nil
}

func signerRegistrySetNewSignerAddressesLogic(env cldf.Environment, config SetNewSignerAddressesConfig) (cldf.ChangesetOutput, error) {
	addressBook := cldf.NewMemoryAddressBook()

	// Load onchain state to get MCMS addresses if needed
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// If using MCMS, we need to collect transactions for the proposal
	var batches []mcmstypes.BatchOperation
	timelocks := make(map[uint64]string)
	inspectors := make(map[uint64]mcmssdk.Inspector)

	for chainSelector, updates := range config.UpdatesByChain {
		chain, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			continue
		}
		// Get addresses for this chain
		addresses, err := env.ExistingAddresses.AddressesForChain(chain.ChainSelector())
		if err != nil {
			env.Logger.Infof("Failed to get addresses for chain %s: %v", chain.String(), err)
			continue
		}

		// Find signer registry address
		var signerRegistryAddress common.Address
		found := false
		for addr, tv := range addresses {
			if tv.Type == shared.EVMSignerRegistry && tv.Version == deployment.Version1_0_0 {
				signerRegistryAddress = common.HexToAddress(addr)
				found = true
				break
			}
		}
		if !found {
			env.Logger.Infof("Signer registry not found on chain %s, skipping", chain.String())
			continue
		}

		signerRegistry, err := signer_registry.NewSignerRegistry(signerRegistryAddress, chain.Client)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to create signer registry instance on %s: %w", chain.String(), err)
		}

		// Prepare arrays for the contract call
		var existingAddresses []common.Address
		var newAddresses []common.Address
		for existing, newAddr := range updates {
			existingAddresses = append(existingAddresses, existing)
			newAddresses = append(newAddresses, newAddr)
		}

		// Execute or prepare the transaction based on MCMS configuration
		txOpts := chain.DeployerKey
		if config.MCMS != nil {
			// Use simulated backend for MCMS to get tx data without sending
			txOpts = cldf.SimTransactOpts()
		}

		tx, err := signerRegistry.SetNewSignerAddresses(txOpts, existingAddresses, newAddresses)

		// Handle based on MCMS configuration
		switch {
		case err != nil:
			// Error preparing transaction
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to prepare transaction for %s: %w", chain.String(), err)
		case config.MCMS == nil:
			// Direct execution - confirm transaction
			_, err = cldf.ConfirmIfNoErrorWithABI(chain, tx, signer_registry.SignerRegistryABI, err)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to set new signer addresses on %s: %w", chain.String(), err)
			}
			env.Logger.Infof("Successfully set new signer addresses on %s (tx: %s)", chain.String(), tx.Hash().Hex())
		default:
			// MCMS mode - prepare batch operation
			if err := stateview.ValidateChain(env, state, chain.ChainSelector(), config.MCMS); err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to validate chain %s for MCMS: %w", chain.String(), err)
			}
			chainState := state.MustGetEVMChainState(chain.ChainSelector())
			if chainState.Timelock == nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("timelock not found on chain %s", chain.String())
			}
			timelocks[chainSelector] = chainState.Timelock.Address().Hex()

			inspector, err := proposalutils.McmsInspectorForChain(env, chain.ChainSelector())
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get inspector for chain %s: %w", chain.String(), err)
			}
			inspectors[chainSelector] = inspector

			batchOperation, err := proposalutils.BatchOperationForChain(
				chainSelector,
				signerRegistryAddress.Hex(),
				tx.Data(),
				big.NewInt(0),
				string(shared.EVMSignerRegistry),
				[]string{},
			)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to create batch operation for chain %s: %w", chain.String(), err)
			}

			batches = append(batches, batchOperation)
			env.Logger.Infof("Prepared transaction for MCMS proposal on %s", chain.String())
		}
	}

	// If using MCMS, build and return the proposal
	if config.MCMS != nil {
		mcmsContractByChain, err := deployergroup.BuildMcmAddressesPerChainByAction(env, state, config.MCMS)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build mcm addresses per chain: %w", err)
		}

		proposal, err := proposalutils.BuildProposalFromBatchesV2(
			env,
			timelocks,
			mcmsContractByChain,
			inspectors,
			batches,
			"Set new signer addresses in SignerRegistry",
			*config.MCMS,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
		}

		env.Logger.Infof("MCMS proposal created with %d operations", len(batches))
		return cldf.ChangesetOutput{
			AddressBook:           addressBook,
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		}, nil
	}

	return cldf.ChangesetOutput{AddressBook: addressBook}, nil
}
