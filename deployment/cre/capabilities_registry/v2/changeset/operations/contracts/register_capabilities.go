package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/capabilities_registry/v2/changeset/pkg"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type RegisterCapabilitiesDeps struct {
	Env      *cldf.Environment
	Strategy strategies.TransactionStrategy
}

type RegisterCapabilitiesInput struct {
	Address       string
	ChainSelector uint64
	Capabilities  []RegisterableCapability
	MCMSConfig    *contracts.MCMSConfig
}

type RegisterCapabilitiesOutput struct {
	Capabilities []RegisterableCapability
	Operation    *mcmstypes.BatchOperation
}

type RegisterableCapability struct {
	CapabilityID          string
	ConfigurationContract common.Address
	Metadata              pkg.CapabilityConfig
}

// RegisterCapabilities is an operation that registers nodes in the V2 Capabilities Registry contract.
var RegisterCapabilities = operations.NewOperation[RegisterCapabilitiesInput, RegisterCapabilitiesOutput, RegisterCapabilitiesDeps](
	"register-capabilities-op",
	semver.MustParse("1.0.0"),
	"Register Capabilities in Capabilities Registry",
	func(b operations.Bundle, deps RegisterCapabilitiesDeps, input RegisterCapabilitiesInput) (RegisterCapabilitiesOutput, error) {
		if len(input.Capabilities) == 0 {
			b.Logger.Info("no capabilities provided, skipping operation")
			return RegisterCapabilitiesOutput{
				Capabilities: []RegisterableCapability{},
			}, nil
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		capabilitiesRegistry, err := capabilities_registry_v2.NewCapabilitiesRegistry(
			common.HexToAddress(input.Address), chain.Client,
		)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistry: %w", err)
		}

		b.Logger.Debugw("registering capabilities", "address", input.Address, "newCapabilities", input.Capabilities, "chainSelector", input.ChainSelector)

		// We have to make sure the capabilities are not already in the contract, to avoid reverting the transaction.
		// This is also important when we use MCMS, so the whole batch doesn't get reverted.
		capabilities, err := dedupCapabilities(capabilitiesRegistry, input.Capabilities)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to deduplicate capabilities: %w", err)
		}

		if len(capabilities) == 0 {
			b.Logger.Info("no new capabilities to register after deduplication, skipping operation")

			return RegisterCapabilitiesOutput{
				Capabilities: []RegisterableCapability{},
			}, nil
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return capabilitiesRegistry.AddCapabilities(opts, capabilities)
		})
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to execute AddCapabilities: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for RegisterCapabilities on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully registered %d capabilities on chain %d", len(capabilities), input.ChainSelector)
		}

		return RegisterCapabilitiesOutput{
			Operation:    operation,
			Capabilities: input.Capabilities,
		}, nil
	},
)

// dedupCapabilities deduplicates the capabilities with respect to the registry
// The contract reverts on adding the same capability twice and that would cause the whole transaction to revert.
func dedupCapabilities(
	capReg *capabilities_registry_v2.CapabilitiesRegistry,
	capabilities []RegisterableCapability,
) ([]capabilities_registry_v2.CapabilitiesRegistryCapability, error) {
	if capReg == nil {
		return nil, errors.New("capabilities registry is nil")
	}
	if len(capabilities) == 0 {
		return nil, errors.New("capabilities list is empty")
	}

	// Fetch all capabilities via generic pagination helper
	caps, err := pkg.GetCapabilities(nil, capReg)
	if err != nil {
		return nil, fmt.Errorf("failed to call GetCapabilities: %w", err)
	}

	existingByID := make(map[string]struct{})
	for _, existingCap := range caps {
		existingByID[existingCap.CapabilityId] = struct{}{}
	}

	var out []capabilities_registry_v2.CapabilitiesRegistryCapability

	// Deduplicate capabilities by their ID
	seen := make(map[string]struct{}, len(capabilities))
	for _, candidate := range capabilities {
		// Process a capability only once in terms of the input list, to avoid duplicates in the output
		if _, exists := seen[candidate.CapabilityID]; exists {
			continue
		}
		seen[candidate.CapabilityID] = struct{}{}

		// Skip capabilities that already exist in the registry
		if _, exists := existingByID[candidate.CapabilityID]; !exists {
			metadataBytes, metadataErr := candidate.Metadata.MarshalJSON()
			if metadataErr != nil {
				return nil, fmt.Errorf("failed to marshal capability metadata for capability %s: %w", candidate.CapabilityID, metadataErr)
			}

			out = append(out, capabilities_registry_v2.CapabilitiesRegistryCapability{
				Metadata:              metadataBytes,
				CapabilityId:          candidate.CapabilityID,
				ConfigurationContract: candidate.ConfigurationContract,
			})
		}
	}

	return out, nil
}
