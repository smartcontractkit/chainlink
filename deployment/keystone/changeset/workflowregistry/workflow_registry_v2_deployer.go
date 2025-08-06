package workflowregistry

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

type workflowRegistryV2Deployer struct {
	lggr     logger.Logger
	contract *workflow_registry_v2.WorkflowRegistry
}

func newWorkflowRegistryV2Deployer() (*workflowRegistryV2Deployer, error) {
	lggr, err := logger.New()
	if err != nil {
		return nil, err
	}
	return &workflowRegistryV2Deployer{lggr: lggr}, nil
}

func (c *workflowRegistryV2Deployer) Contract() *workflow_registry_v2.WorkflowRegistry {
	return c.contract
}

func (c *workflowRegistryV2Deployer) Deploy(req changeset.DeployRequest) (*changeset.DeployResponse, error) {
	addr, tx, wr, err := workflow_registry_v2.DeployWorkflowRegistry(
		req.Chain.DeployerKey,
		req.Chain.Client)
	if err != nil {
		return nil, cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
	}

	_, err = req.Chain.Confirm(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm and save WorkflowRegistry v2: %w", err)
	}
	tvStr, err := wr.TypeAndVersion(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("failed to get type and version: %w", err)
	}

	tv, err := cldf.TypeAndVersionFromString(tvStr)
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

// DeployWorkflowRegistryV2 deploys the WorkflowRegistry v2 contract to the chain
// and saves the address in the address book. This mutates the address book.
func DeployWorkflowRegistryV2(ctx context.Context, chain cldf_evm.Chain, ab cldf.AddressBook) (*changeset.DeployResponse, error) {
	deployer, err := newWorkflowRegistryV2Deployer()
	if err != nil {
		return nil, fmt.Errorf("failed to create WorkflowRegistry v2 deployer: %w", err)
	}
	resp, err := deployer.Deploy(changeset.DeployRequest{Chain: chain})
	if err != nil {
		return nil, fmt.Errorf("failed to deploy WorkflowRegistry v2: %w", err)
	}
	err = ab.Save(chain.Selector, resp.Address.String(), resp.Tv)
	if err != nil {
		return nil, fmt.Errorf("failed to save WorkflowRegistry v2: %w", err)
	}
	return resp, nil
}
