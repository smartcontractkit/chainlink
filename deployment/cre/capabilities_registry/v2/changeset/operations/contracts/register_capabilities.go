package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type RegisterCapabilitiesDeps struct {
	Env *cldf.Environment
}

type RegisterCapabilitiesInput struct {
	Address       string
	ChainSelector uint64
	Capabilities  []capabilities_registry_v2.CapabilitiesRegistryCapability
}

type RegisterCapabilitiesOutput struct {
	Capabilities []*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfigured
}

// RegisterCapabilities is an operation that registers nodes in the V2 Capabilities Registry contract.
var RegisterCapabilities = operations.NewOperation[RegisterCapabilitiesInput, RegisterCapabilitiesOutput, RegisterCapabilitiesDeps](
	"register-capabilities-op",
	semver.MustParse("1.0.0"),
	"Register Capabilities in Capabilities Registry",
	func(b operations.Bundle, deps RegisterCapabilitiesDeps, input RegisterCapabilitiesInput) (RegisterCapabilitiesOutput, error) {
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

		// We have to make sure the capabilities are not already in the contract, to avoid reverting the transaction.
		// This is also important when we use MCMS, so the whole batch doesn't get reverted.
		capabilities, err := dedupCapabilities(capabilitiesRegistry, input.Capabilities)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to deduplicate capabilities: %w", err)
		}

		tx, err := capabilitiesRegistry.AddCapabilities(chain.DeployerKey, capabilities)
		if err != nil {
			err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to call AddCapabilities: %w", err)
		}

		_, err = chain.Confirm(tx)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to confirm AddCapabilities transaction %s: %w", tx.Hash().String(), err)
		}

		ctx := b.GetContext()
		receipt, err := bind.WaitMined(ctx, chain.Client, tx)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to mine AddCapabilities confirm transaction %s: %w", tx.Hash().String(), err)
		}

		resp := RegisterCapabilitiesOutput{
			Capabilities: make([]*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfigured, 0, len(receipt.Logs)),
		}
		// Parse the logs to get the added capabilities
		for i, log := range receipt.Logs {
			if log == nil {
				continue
			}

			o, err := capabilitiesRegistry.ParseCapabilityConfigured(*log)
			if err != nil {
				return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to parse log %d for capability added: %w", i, err)
			}
			resp.Capabilities = append(resp.Capabilities, o)
		}

		return resp, nil
	},
)

// dedupCapabilities deduplicates the capabilities with respect to the registry
// The contract reverts on adding the same capability twice and that would cause the whole transaction to revert.
func dedupCapabilities(
	capReg *capabilities_registry_v2.CapabilitiesRegistry,
	capabilities []capabilities_registry_v2.CapabilitiesRegistryCapability,
) ([]capabilities_registry_v2.CapabilitiesRegistryCapability, error) {
	if capReg == nil {
		return nil, errors.New("capabilities registry is nil")
	}
	if len(capabilities) == 0 {
		return nil, errors.New("capabilities list is empty")
	}

	caps, err := capReg.GetCapabilities(nil)
	if err != nil {
		err = cldf.DecodeErr(capabilities_registry_v2.CapabilitiesRegistryABI, err)
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
		if _, exists := seen[candidate.CapabilityId]; exists {
			continue
		}
		seen[candidate.CapabilityId] = struct{}{}

		// Skip capabilities that already exist in the registry
		if _, exists := existingByID[candidate.CapabilityId]; !exists {
			out = append(out, candidate)
		}
	}

	return out, nil
}
