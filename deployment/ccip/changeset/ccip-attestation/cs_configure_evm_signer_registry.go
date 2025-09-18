package ccip_attestation

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
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
		_, ok := env.BlockChains.EVMChains()[chainSelector]
		if !ok {
			return fmt.Errorf("chain selector %d not found in environment", chainSelector)
		}

		if chainSelector != BaseMainnetSelector && chainSelector != BaseSepoliaSelector {
			return fmt.Errorf("chain selector %d is not a Base chain", chainSelector)
		}

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
		var signerRegistrySigners []signer_registry.ISignerRegistrySigner
		if chainState.SignerRegistry != nil {
			signerRegistrySigners, err = chainState.SignerRegistry.GetSigners(&bind.CallOpts{Context: env.GetContext()})
		}
		if err != nil {
			return fmt.Errorf("failed to get signers from signer registry on chain selector %d: %w", chainSelector, err)
		}

		if len(signerRegistrySigners) == 0 {
			env.Logger.Infof("No signer registry data found on chain selector %d, skipping", chainSelector)
			continue
		}

		existingSigners := make(map[common.Address]bool)
		for _, signer := range signerRegistrySigners {
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
	// Load onchain state to get MCMS addresses if needed
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	deployerGroup := deployergroup.NewDeployerGroup(env, state, config.MCMS).WithDeploymentContext("configure signer registry with new signer addresses")

	for chainSelector, updates := range config.UpdatesByChain {
		chainState := state.Chains[chainSelector]
		signerRegistry := chainState.SignerRegistry
		if signerRegistry == nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("no signer registry found on chain selector %d", chainSelector)
		}
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain selector %d: %w", chainSelector, err)
		}

		// Prepare arrays for the contract call
		var existingAddresses []common.Address
		var newAddresses []common.Address
		for existing, newAddr := range updates {
			existingAddresses = append(existingAddresses, existing)
			newAddresses = append(newAddresses, newAddr)
		}

		_, err = signerRegistry.SetNewSignerAddresses(opts, existingAddresses, newAddresses)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to set new signer addresses on chain selector %d: %w", chainSelector, err)
		}
	}

	return deployerGroup.Enact()
}
