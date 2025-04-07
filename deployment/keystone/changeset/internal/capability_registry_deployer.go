package internal

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	"github.com/smartcontractkit/chainlink/deployment"
)

type CapabilitiesRegistryDeployer struct {
	lggr     logger.Logger
	contract *capabilities_registry.CapabilitiesRegistry
}

func NewCapabilitiesRegistryDeployer() (*CapabilitiesRegistryDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &CapabilitiesRegistryDeployer{lggr: lggr}, nil
}

func (c *CapabilitiesRegistryDeployer) Contract() *capabilities_registry.CapabilitiesRegistry {
	return c.contract
}

func (c *CapabilitiesRegistryDeployer) Deploy(req DeployRequest) (*DeployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, capabilities_registry.CapabilitiesRegistryABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Debugf("Capability registry estimated gas: %d", est)

	capabilitiesRegistryAddr, tx, capabilitiesRegistry, err := capabilities_registry.DeployCapabilitiesRegistry(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save CapabilitiesRegistry: %w", err)
	}
	tvStr, err := capabilitiesRegistry.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}

	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	resp := &DeployResponse{
		Address: capabilitiesRegistryAddr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = capabilitiesRegistry
	return resp, nil
}

func estimateDeploymentGas(client deployment.OnchainClient, bytecode string) (uint64, error) {
	// fake contract address required for gas estimation, otherwise it will fail
	contractAddress := common.HexToAddress("0x0000000000000000000000000000000000000000")

	msg := ethereum.CallMsg{
		To:   &contractAddress, // nil ok for
		Gas:  0,                // initial gas estimate (will be updated)
		Data: []byte(bytecode),
	}
	gasEstimate, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, fmt.Errorf("failed to estimate gas: %w", err)
	}
	return gasEstimate, nil
}
