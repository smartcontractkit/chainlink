package keystone

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
)

type KeystoneForwarderDeployer struct {
	lggr     logger.Logger
	contract *forwarder.KeystoneForwarder
}

var ForwarderTypeVersion = deployment.TypeAndVersion{
	Type:    KeystoneForwarder,
	Version: deployment.Version1_0_0,
}

func (c *KeystoneForwarderDeployer) deploy(req deployRequest) (*deployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, forwarder.KeystoneForwarderABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Debugf("Forwarder estimated gas: %d", est)

	forwarderAddr, tx, forwarder, err := forwarder.DeployKeystoneForwarder(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy KeystoneForwarder: %w", err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save KeystoneForwarder: %w", err)
	}
	resp := &deployResponse{
		Address: forwarderAddr,
		Tx:      tx.Hash(),
		Tv:      ForwarderTypeVersion,
	}
	c.contract = forwarder
	return resp, nil
}
