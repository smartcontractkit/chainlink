package contracts

import (
	"context"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
)

type DeployCapabilitiesRegistryOpDeps struct {
	Env *cldf.Environment
}

type DeployCapabilitiesRegistryOpInput struct {
	ChainSelector uint64
	Qualifier     string
}

type DeployCapabilitiesRegistryOpOutput struct {
	Address       string
	ChainSelector uint64
	Qualifier     string
	Type          string
	Version       string
	Labels        []string
}

// DeployCapabilitiesRegistryOp is an operation that deploys the V2 Capabilities Registry contract.
// This atomic operation performs the single side effect of deploying and registering the contract.
var DeployCapabilitiesRegistryOp = operations.NewOperation[DeployCapabilitiesRegistryOpInput, DeployCapabilitiesRegistryOpOutput, DeployCapabilitiesRegistryOpDeps](
	"deploy-capabilities-registry-v2-op",
	semver.MustParse("1.0.0"),
	"Deploy CapabilitiesRegistry V2 Contract",
	func(b operations.Bundle, deps DeployCapabilitiesRegistryOpDeps, input DeployCapabilitiesRegistryOpInput) (DeployCapabilitiesRegistryOpOutput, error) {
		lggr := deps.Env.Logger

		// Get the target chain
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return DeployCapabilitiesRegistryOpOutput{}, fmt.Errorf("chain not found for selector %d", input.ChainSelector)
		}

		// Estimate gas for deployment
		est, err := estimateDeploymentGas(chain.Client, capabilities_registry_v2.CapabilitiesRegistryABI)
		if err != nil {
			return DeployCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to estimate gas: %w", err)
		}
		lggr.Debugf("Capabilities Registry V2 estimated gas: %d", est)

		// Deploy the V2 CapabilitiesRegistry contract
		capabilitiesRegistryAddr, tx, capabilitiesRegistry, err := capabilities_registry_v2.DeployCapabilitiesRegistry(
			chain.DeployerKey,
			chain.Client,
			capabilities_registry_v2.CapabilitiesRegistryConstructorParams{},
		)
		if err != nil {
			return DeployCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to deploy CapabilitiesRegistry V2: %w", err)
		}

		// Wait for deployment confirmation
		_, err = chain.Confirm(tx)
		if err != nil {
			return DeployCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to confirm CapabilitiesRegistry V2 deployment: %w", err)
		}

		// Get type and version from the deployed contract
		tvStr, err := capabilitiesRegistry.TypeAndVersion(&bind.CallOpts{})
		if err != nil {
			return DeployCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to get type and version: %w", err)
		}

		tv, err := cldf.TypeAndVersionFromString(tvStr)
		if err != nil {
			return DeployCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
		}

		lggr.Infof("Deployed %s on chain selector %d at address %s", tv.String(), chain.Selector, capabilitiesRegistryAddr.String())

		return DeployCapabilitiesRegistryOpOutput{
			Address:       capabilitiesRegistryAddr.String(),
			ChainSelector: input.ChainSelector,
			Qualifier:     input.Qualifier,
			Type:          string(tv.Type),
			Version:       tv.Version.String(),
			Labels:        tv.Labels.List(),
		}, nil
	},
)

func estimateDeploymentGas(client cldf_evm.OnchainClient, bytecode string) (uint64, error) {
	// fake contract address required for gas estimation, otherwise it will fail
	contractAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Gas:  0,
		Data: []byte(bytecode),
	}
	gasEstimate, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}
	return gasEstimate, nil
}
