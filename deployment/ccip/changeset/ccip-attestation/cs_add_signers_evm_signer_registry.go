package ccip_attestation

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/bindings/signer_registry"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/deployergroup"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

var (
	EVMSignerRegistryAddSignersChangeset = cldf.CreateChangeSet(addSignersLogic, addSignersPrecondition)
)

// AddSignersConfig drives EVMSignerRegistryAddSignersChangeset. The changeset
// appends signers to an already-deployed SignerRegistry via addSigners(). It
// does not touch existing signers; use EVMSignerRegistrySetNewSignerAddressesChangeset
// for key rotation.
type AddSignersConfig struct {
	// MCMS, if non-nil, routes the tx through the Timelock via MCMS proposal.
	// Leave nil to execute directly with the deployer key (only works if the
	// SignerRegistry owner is the deployer).
	MCMS *proposalutils.TimelockConfig
	// SignersByChain maps a Base chain selector to the signers to append.
	// NewEVMAddress must be the zero address; use the rotation changeset to
	// set pending signer addresses.
	SignersByChain map[uint64][]signer_registry.ISignerRegistrySigner
}

func addSignersPrecondition(env cldf.Environment, config AddSignersConfig) error {
	if len(config.SignersByChain) == 0 {
		return errors.New("no signers provided")
	}

	for chainSelector, signers := range config.SignersByChain {
		if _, ok := env.BlockChains.EVMChains()[chainSelector]; !ok {
			return fmt.Errorf("chain selector %d not found in environment", chainSelector)
		}
		if chainSelector != BaseMainnetSelector && chainSelector != BaseSepoliaSelector {
			return fmt.Errorf("chain selector %d is not a Base chain", chainSelector)
		}
		if len(signers) == 0 {
			return fmt.Errorf("no signers provided for chain selector %d", chainSelector)
		}

		seen := make(map[common.Address]struct{})
		for _, s := range signers {
			if s.EvmAddress == utils.ZeroAddress {
				return fmt.Errorf("signer has zero evm address on chain selector %d", chainSelector)
			}
			if s.NewEVMAddress != utils.ZeroAddress {
				return fmt.Errorf("signer %s has non-zero new evm address on chain selector %d", s.EvmAddress.Hex(), chainSelector)
			}
			if _, dup := seen[s.EvmAddress]; dup {
				return fmt.Errorf("duplicate signer evm address %s on chain selector %d", s.EvmAddress.Hex(), chainSelector)
			}
			seen[s.EvmAddress] = struct{}{}
		}
	}

	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	for chainSelector, signers := range config.SignersByChain {
		chainState, ok := state.Chains[chainSelector]
		if !ok {
			return fmt.Errorf("no chain state for chain selector %d", chainSelector)
		}
		registry := chainState.SignerRegistry
		if registry == nil {
			return fmt.Errorf("no signer registry found on chain selector %d", chainSelector)
		}

		callOpts := &bind.CallOpts{Context: env.GetContext()}

		maxSigners, err := registry.GetMaxSigners(callOpts)
		if err != nil {
			return fmt.Errorf("failed to get max signers on chain selector %d: %w", chainSelector, err)
		}
		currentCount, err := registry.GetSignerCount(callOpts)
		if err != nil {
			return fmt.Errorf("failed to get signer count on chain selector %d: %w", chainSelector, err)
		}
		if currentCount.Uint64()+uint64(len(signers)) > maxSigners.Uint64() {
			return fmt.Errorf("chain selector %d: adding %d signers would exceed max signers (current=%d, max=%d)",
				chainSelector, len(signers), currentCount.Uint64(), maxSigners.Uint64())
		}

		existing, err := registry.GetSigners(callOpts)
		if err != nil {
			return fmt.Errorf("failed to get signers on chain selector %d: %w", chainSelector, err)
		}
		existingSet := make(map[common.Address]struct{}, len(existing)*2)
		for _, s := range existing {
			existingSet[s.EvmAddress] = struct{}{}
			if s.NewEVMAddress != utils.ZeroAddress {
				existingSet[s.NewEVMAddress] = struct{}{}
			}
		}
		for _, s := range signers {
			if _, dup := existingSet[s.EvmAddress]; dup {
				return fmt.Errorf("signer %s is already registered on chain selector %d", s.EvmAddress.Hex(), chainSelector)
			}
		}
	}

	return nil
}

func addSignersLogic(env cldf.Environment, config AddSignersConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(env)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	deployerGroup := deployergroup.NewDeployerGroup(env, state, config.MCMS).
		WithDeploymentContext("add signers to EVM SignerRegistry")

	for chainSelector, signers := range config.SignersByChain {
		chainState := state.Chains[chainSelector]
		registry := chainState.SignerRegistry
		if registry == nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("no signer registry found on chain selector %d", chainSelector)
		}
		opts, err := deployerGroup.GetDeployer(chainSelector)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer for chain selector %d: %w", chainSelector, err)
		}
		if _, err := registry.AddSigners(opts, signers); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to add signers on chain selector %d: %w", chainSelector, err)
		}
	}

	return deployerGroup.Enact()
}
