package workflowregistry

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var _ cldf.ChangeSet[*UpdateAuthorizedAddressesRequest] = UpdateAuthorizedAddresses

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

func getWorkflowRegistry(env cldf.Environment, chainSel uint64) (*workflow_registry.WorkflowRegistry, error) {
	resp, err := changeset.GetContractSets(env.Logger, &changeset.GetContractSetsRequest{
		Chains:      env.BlockChains.EVMChains(),
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
func UpdateAuthorizedAddresses(env cldf.Environment, req *UpdateAuthorizedAddressesRequest) (cldf.ChangesetOutput, error) {
	if err := req.Validate(); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	evmChains := env.BlockChains.EVMChains()
	resp, err := changeset.GetContractSets(env.Logger, &changeset.GetContractSetsRequest{
		Chains:      evmChains,
		AddressBook: env.ExistingAddresses,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get contract sets: %w", err)
	}

	cs := resp.ContractSets[req.RegistryChainSel]
	if cs.WorkflowRegistry == nil {
		return cldf.ChangesetOutput{}, errors.New("could not find workflow registry")
	}
	registry := cs.WorkflowRegistry

	chain, ok := evmChains[req.RegistryChainSel]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
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
			Env:         env,
		}
	} else {
		s = &simpleTransaction{
			chain: chain,
		}
	}

	return s.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		tx, err := registry.UpdateAuthorizedAddresses(opts, addr, req.Allowed)
		if err != nil {
			err = cldf.DecodeErr(workflow_registry.WorkflowRegistryABI, err)
		}
		return tx, err
	})
}
