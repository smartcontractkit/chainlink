package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[*AppendNodeCapabilitiesRequest] = AppendNodeCapabilities

// AppendNodeCapabilitiesRequest is a request to add capabilities to the existing capabilities of nodes in the registry
type AppendNodeCapabilitiesRequest = MutateNodeCapabilitiesRequest

// AppendNodeCapabilities adds any new capabilities to the registry, merges the new capabilities with the existing capabilities
// of the node, and updates the nodes in the registry host the union of the new and existing capabilities.
func AppendNodeCapabilities(env deployment.Environment, req *AppendNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	cfg, err := req.convert(env)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	_, err = internal.AppendNodeCapabilitiesImpl(env.Logger, cfg)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}

func (req *AppendNodeCapabilitiesRequest) convert(e deployment.Environment) (*internal.AppendNodeCapabilitiesRequest, error) {
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate UpdateNodeCapabilitiesRequest: %w", err)
	}
	registryChain, ok := e.Chains[req.RegistryChainSel]
	if !ok {
		return nil, fmt.Errorf("registry chain selector %d does not exist in environment", req.RegistryChainSel)
	}
	contracts, err := kslib.GetContractSets(e.Logger, &kslib.GetContractSetsRequest{
		Chains:      map[uint64]deployment.Chain{req.RegistryChainSel: registryChain},
		AddressBook: req.AddressBook,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	registry := contracts.ContractSets[req.RegistryChainSel].CapabilitiesRegistry
	if registry == nil {
		return nil, fmt.Errorf("capabilities registry not found for chain %d", req.RegistryChainSel)
	}

	return &internal.AppendNodeCapabilitiesRequest{
		Chain:             registryChain,
		Registry:          registry,
		P2pToCapabilities: req.P2pToCapabilities,
	}, nil
}
