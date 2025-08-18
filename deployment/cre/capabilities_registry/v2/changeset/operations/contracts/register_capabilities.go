package contracts

import (
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
		// Validate input

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Get the CapabilitiesRegistryTransactor contract
		capabilityRegistryTransactor, err := capabilities_registry_v2.NewCapabilitiesRegistryTransactor(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryTransactor: %w", err)
		}

		tx, err := capabilityRegistryTransactor.AddCapabilities(chain.DeployerKey, input.Capabilities)
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

		// Get the CapabilitiesRegistryFilterer contract
		capabilityRegistryFilterer, err := capabilities_registry_v2.NewCapabilitiesRegistryFilterer(
			common.HexToAddress(input.Address),
			chain.Client,
		)
		if err != nil {
			return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to create CapabilitiesRegistryFilterer: %w", err)
		}

		resp := RegisterCapabilitiesOutput{
			Capabilities: make([]*capabilities_registry_v2.CapabilitiesRegistryCapabilityConfigured, 0, len(receipt.Logs)),
		}
		// Parse the logs to get the added capabilities
		for i, log := range receipt.Logs {
			if log == nil {
				continue
			}

			o, err := capabilityRegistryFilterer.ParseCapabilityConfigured(*log)
			if err != nil {
				return RegisterCapabilitiesOutput{}, fmt.Errorf("failed to parse log %d for capability added: %w", i, err)
			}
			resp.Capabilities = append(resp.Capabilities, o)
		}

		return resp, nil
	},
)
