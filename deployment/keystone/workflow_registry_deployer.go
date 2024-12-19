package keystone

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	workflow_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

type WorkflowRegistryDeployer struct {
	lggr     logger.Logger
	contract *workflow_registry.WorkflowRegistry
}

func NewWorkflowRegistryDeployer() (*WorkflowRegistryDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &WorkflowRegistryDeployer{lggr: lggr}, nil
}

func (c *WorkflowRegistryDeployer) Contract() *workflow_registry.WorkflowRegistry {
	return c.contract
}

func (c *WorkflowRegistryDeployer) Deploy(req DeployRequest) (*DeployResponse, error) {
	est, err := estimateDeploymentGas(req.Chain.Client, workflow_registry.WorkflowRegistryABI)
	if err != nil {
		return nil, fmt.Errorf("failed to estimate gas: %w", err)
	}
	c.lggr.Debugf("WorkflowRegistry estimated gas: %d", est)

	addr, tx, wr, err := workflow_registry.DeployWorkflowRegistry(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, DecodeErr(workflow_registry.WorkflowRegistryABI, err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save WorkflowRegistry: %w", err)
	}
	tvStr, err := wr.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}

	tv, err := deployment.TypeAndVersionFromString(tvStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse type and version from %s: %w", tvStr, err)
	}
	resp := &DeployResponse{
		Address: addr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = wr
	return resp, nil
}
