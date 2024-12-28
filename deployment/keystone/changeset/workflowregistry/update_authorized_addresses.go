package workflowregistry

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	workflow_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

var _ deployment.ChangeSet[*UpdateAuthorizedAddressesRequest] = UpdateAuthorizedAddresses

type UpdateAuthorizedAddressesRequest struct {
	RegistryChainSel uint64

	Addresses []string
	Allowed   bool

	MCMSConfig *changeset.MCMSConfig
}

func (r *UpdateAuthorizedAddressesRequest) Validate() error {
	if len(r.Addresses) == 0 {
		return errors.New("Must provide at least 1 address")
	}

	return nil
}

func getWorkflowRegistry(env deployment.Environment, chainSel uint64) (*workflow_registry.WorkflowRegistry, error) {
	resp, err := changeset.GetContractSets(env.Logger, &changeset.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

	cs := resp.ContractSets[chainSel]
	if cs.WorkflowRegistry == nil {
		return nil, errors.New("could not find workflow registry")
	}

	return cs.WorkflowRegistry, nil
}

// UpdateAuthorizedAddresses updates the list of DONs that workflows can be sent to.
func UpdateAuthorizedAddresses(env deployment.Environment, req *UpdateAuthorizedAddressesRequest) (deployment.ChangesetOutput, error) {
	if err := req.Validate(); err != nil {
		return deployment.ChangesetOutput{}, err
	}

	resp, err := changeset.GetContractSets(env.Logger, &changeset.GetContractSetsRequest{
		Chains:      env.Chains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	cs := resp.ContractSets[req.RegistryChainSel]
	if cs.WorkflowRegistry == nil {
		return deployment.ChangesetOutput{}, errors.New("could not find workflow registry")
	}
	registry := cs.WorkflowRegistry

	chain, ok := env.Chains[req.RegistryChainSel]
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}

	var addr []common.Address
	for _, a := range req.Addresses {
		addr = append(addr, common.HexToAddress(a))
	}

	var s strategy
	if req.MCMSConfig != nil {
		s = &mcmsTransaction{
			Config:      req.MCMSConfig,
			Description: "proposal to update authorized addresses",
			Address:     registry.Address(),
			ChainSel:    chain.Selector,
			ContractSet: &cs,
		}
	} else {
		s = &simpleTransaction{
			chain: chain,
		}
	}

	return s.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		tx, err := registry.UpdateAuthorizedAddresses(opts, addr, req.Allowed)
		if err != nil {
			err = deployment.DecodeErr(workflow_registry.WorkflowRegistryABI, err)
		}
		return tx, err
	})
}
