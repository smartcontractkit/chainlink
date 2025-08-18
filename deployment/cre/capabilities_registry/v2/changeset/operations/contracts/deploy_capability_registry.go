package contracts

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type DeployCapabilitiesRegistryDeps struct {
	Env *cldf.Environment
}

type DeployCapabilitiesRegistryInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployCapabilitiesRegistryOutput struct {
	Address       string
	ChainSelector uint64
	Qualifier     string
	Type          string
	Version       string
	Labels        []string
}

// DeployCapabilitiesRegistry is an operation that deploys the V2 Capabilities Registry contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployCapabilitiesRegistry = operations.NewOperation[DeployCapabilitiesRegistryInput, DeployCapabilitiesRegistryOutput, DeployCapabilitiesRegistryDeps](
	"deploy-capabilities-registry-v2-op",
	semver.MustParse("1.0.0"),
	"Deploy CapabilitiesRegistry V2 Contract",
	func(b operations.Bundle, deps DeployCapabilitiesRegistryDeps, input DeployCapabilitiesRegistryInput) (DeployCapabilitiesRegistryOutput, error) {
		lggr := deps.Env.Logger

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployCapabilitiesRegistryOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Deploy the V2 CapabilitiesRegistry contract
		capabilitiesRegistryAddr, tx, capabilitiesRegistry, err := capabilities_registry_v2.DeployCapabilitiesRegistry(
			chain.DeployerKey,
			chain.Client,
			capabilities_registry_v2.CapabilitiesRegistryConstructorParams{},
		)
		if err != nil {
			return DeployCapabilitiesRegistryOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry V2: %w", err)
		}

		// Wait for deployment confirmation
		_, err = chain.Confirm(tx)
		if err != nil {
			return DeployCapabilitiesRegistryOutput{}, fmt.Errorf("failed to confirm CapabilitiesRegistry V2 deployment: %w", err)
		}

		// Get type and version from the deployed contract
		tvStr, err := capabilitiesRegistry.TypeAndVersion(&bind.CallOpts{})
		if err != nil {
			return DeployCapabilitiesRegistryOutput{}, fmt.Errorf("failed to get type and version: %w", err)
		}

		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return DeployCapabilitiesRegistryOutput{}, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}

		lggr.Infof("Deployed %s on chain selector %d at address %s", tv.String(), chain.Selector, capabilitiesRegistryAddr.String())

		return DeployCapabilitiesRegistryOutput{
			Address:       capabilitiesRegistryAddr.String(),
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
			Type:          string(tv.Type),
			Version:       tv.Version.String(),
			Labels:        tv.Labels.List(),
		}, nil
	},
)
