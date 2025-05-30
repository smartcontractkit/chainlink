package workflowregistry

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"

	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
)

var _ cldf.ChangeSet[*UpdateAllowedDonsRequest] = UpdateAllowedDons

type UpdateAllowedDonsRequest struct {
	RegistryChainSel uint64
	DonIDs           []uint32
	Allowed          bool

	MCMSConfig *changeset.MCMSConfig
}

func (r *UpdateAllowedDonsRequest) Validate() error {
	if len(r.DonIDs) == 0 {
		return errors.New("Must provide at least one DonID")
	}
	return nil
}

// UpdateAllowedDons updates the list of DONs that workflows can be sent to.
func UpdateAllowedDons(env cldf.Environment, req *UpdateAllowedDonsRequest) (cldf.ChangesetOutput, error) {
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

	var s strategy
	if req.MCMSConfig != nil {
		s = &mcmsTransaction{
			Config:      req.MCMSConfig,
			Description: "proposal to update allowed dons",
			Address:     registry.Address(),
			ChainSel:    req.RegistryChainSel,
			ContractSet: &cs,
			Env:         env,
		}
	} else {
		s = &simpleTransaction{
			chain: chain,
		}
	}

	return s.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
		tx, err := registry.UpdateAllowedDONs(opts, req.DonIDs, req.Allowed)
		if err != nil {
			err = cldf.DecodeErr(workflow_registry.WorkflowRegistryABI, err)
		}
		return tx, err
	})
}
