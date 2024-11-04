package changeset

import (
	"encoding/json"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet = UpdateNodeCapabilities

type P2PSignerEnc = internal.P2PSignerEnc

type UpdateNodeCapabilitiesRequest struct {
	AddressBook      deployment.AddressBook
	RegistryChainSel uint64

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*P2PSignerEnc
}

func (req *UpdateNodeCapabilitiesRequest) Validate() error {
	if req.AddressBook == nil {
		return fmt.Errorf("address book is nil")
	}
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if len(req.NopToNodes) == 0 {
		return fmt.Errorf("nopToNodes is empty")
	}
	_, exists := chainsel.ChainBySelector(req.RegistryChainSel)
	if !exists {
		return fmt.Errorf("registry chain selector %d does not exist", req.RegistryChainSel)
	}
	return nil
}

type UpdateNodeCapabilitiesImplRequest struct {
	Chain    deployment.Chain
	Registry *kcr.CapabilitiesRegistry

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
	NopToNodes        map[kcr.CapabilitiesRegistryNodeOperator][]*internal.P2PSignerEnc
}

func (req *UpdateNodeCapabilitiesImplRequest) Validate() error {
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	if len(req.NopToNodes) == 0 {
		return fmt.Errorf("nopToNodes is empty")
	}
	if req.Registry == nil {
		return fmt.Errorf("registry is nil")
	}

	return nil
}

func (req *UpdateNodeCapabilitiesRequest) updateNodeCapabilitiesImplRequest(e deployment.Environment) (*internal.UpdateNodeCapabilitiesImplRequest, error) {
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
	return &internal.UpdateNodeCapabilitiesImplRequest{
		Chain:             registryChain,
		Registry:          registry,
		P2pToCapabilities: req.P2pToCapabilities,
		NopToNodes:        req.NopToNodes,
	}, nil
}

// UpdateNodeCapabilities updates the capabilities of nodes in the registry
func UpdateNodeCapabilities(env deployment.Environment, config any) (deployment.ChangesetOutput, error) {
	req, ok := config.(*UpdateNodeCapabilitiesRequest)
	if !ok {
		return deployment.ChangesetOutput{}, fmt.Errorf("invalid config type. want %T, got %T", &UpdateNodeCapabilitiesRequest{}, config)
	}
	c, err := req.updateNodeCapabilitiesImplRequest(env)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to convert request: %w", err)
	}

	r, err := internal.UpdateNodeCapabilitiesImpl(env.Logger, c)
	if err == nil {
		b, err2 := json.Marshal(r)
		if err2 != nil {
			env.Logger.Debugf("Updated node capabilities '%s'", b)
		}
	}
	return deployment.ChangesetOutput{}, err
}
