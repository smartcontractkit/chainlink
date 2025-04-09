package workflowregistry

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type workflowRegistryDeployer struct {
	lggr     logger.Logger
	contract *workflow_registry.WorkflowRegistry
}

func newWorkflowRegistryDeployer() (*workflowRegistryDeployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &workflowRegistryDeployer{lggr: lggr}, nil
}

func (c *workflowRegistryDeployer) Contract() *workflow_registry.WorkflowRegistry {
	return c.contract
}

func (c *workflowRegistryDeployer) Deploy(req changeset.DeployRequest) (*changeset.DeployResponse, error) {
	addr, tx, wr, err := workflow_registry.DeployWorkflowRegistry(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, deployment.DecodeErr(workflow_registry.WorkflowRegistryABI, err)
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
	resp := &changeset.DeployResponse{
		Address: addr,
		Tx:      tx.Hash(),
		Tv:      tv,
	}
	c.contract = wr
	return resp, nil
}

// deployWorkflowRegistry deploys the WorkflowRegistry contract to the chain
// and saves the address in the address book. This mutates the address book.
func deployWorkflowRegistry(chain deployment.Chain, ab deployment.AddressBook) (*changeset.DeployResponse, error) {
	deployer, err := newWorkflowRegistryDeployer()
	resp, err := deployer.Deploy(changeset.DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy WorkflowRegistry: %w", err)
	}
	err = ab.Save(chain.Selector, resp.Address.String(), resp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save WorkflowRegistry: %w", err)
	}
	return resp, nil
}
