package keystone

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
)

type CapabilitiesRegistryDeployer struct {
	lggr     logger.Logger
	contract *capabilities_registry.CapabilitiesRegistry
}

func (c *CapabilitiesRegistryDeployer) deploy(req deployRequest) (*deployResponse, error) {
	capabilitiesRegistryAddr, tx, capabilitiesRegistry, err := capabilities_registry.DeployCapabilitiesRegistry(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy CapabilitiesRegistry: %w", err)
	}

	_, err = req.Chain.Confirm(tx.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save CapabilitiesRegistry: %w", err)
	}
	resp := &deployResponse{
		Address: capabilitiesRegistryAddr,
		Tx:      tx.Hash(),
		Tv: deployment.TypeAndVersion{
			Type:    CapabilitiesRegistry,
			Version: deployment.Version1_0_0,
		},
	}
	c.contract = capabilitiesRegistry
	return resp, nil
}
