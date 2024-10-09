package keystone

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

type CapabilitiesRegistryDeployer struct {
	lggr     logger.Logger
	contract *capabilities_registry.CapabilitiesRegistry
}

var CapabilityRegistryTypeVersion = deployment.TypeAndVersion{
	Type:    CapabilitiesRegistry,
	Version: deployment.Version1_0_0,
}

func (c *CapabilitiesRegistryDeployer) deploy(req deployRequest) (*deployResponse, error) {
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
	resp := &deployResponse{
		Address: capabilitiesRegistryAddr,
		Tx:      tx.Hash(),
		Tv:      CapabilityRegistryTypeVersion,
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
