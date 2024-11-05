package changeset

import (
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet = AppendNodeCapabilities

type AppendNodeCapabilitiesRequest struct {
	AddressBook      deployment.AddressBook
	RegistryChainSel uint64

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*P2PSignerEnc
}

func (req *AppendNodeCapabilitiesRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if len(req.NopToNodes) == 0 {
		return fmt.Errorf("nopToNodes is empty")
	}
	if req.AddressBook == nil {
		return fmt.Errorf("registry is nil")
	}
	_, exists := chainsel.ChainBySelector(req.RegistryChainSel)
	if !exists {
		return fmt.Errorf("registry chain selector %d does not exist", req.RegistryChainSel)
	}

	return nil
}

/*
// AppendNodeCapabilibity adds any new capabilities to the registry, merges the new capabilities with the existing capabilities
// of the node, and updates the nodes in the registry host the union of the new and existing capabilities.
func AppendNodeCapabilities(lggr logger.Logger, req *AppendNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
	_, err := appendNodeCapabilitiesImpl(lggr, req)
	if err != nil {
		return deployment.ChangesetOutput{}, err
	}
	return deployment.ChangesetOutput{}, nil
}
*/

// AppendNodeCapabilibity adds any new capabilities to the registry, merges the new capabilities with the existing capabilities
// of the node, and updates the nodes in the registry host the union of the new and existing capabilities.
func AppendNodeCapabilities(env deployment.Environment, config any) (deployment.ChangesetOutput, error) {
	req, ok := config.(*AppendNodeCapabilitiesRequest)
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid config type")
	}

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
	contracts, err := kslib.GetContractSets(&kslib.GetContractSetsRequest{
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
		NopToNodes:        req.NopToNodes,
	}, nil
}
