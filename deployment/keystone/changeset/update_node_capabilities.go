package changeset

import (
	"encoding/json"
	"fmt"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone"
	kslib "github.com/smartcontractkit/chainlink/deployment/keystone"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
)

var _ deployment.ChangeSet[*MutateNodeCapabilitiesRequest] = UpdateNodeCapabilities

type P2PSignerEnc = internal.P2PSignerEnc

func NewP2PSignerEnc(n *keystone.Node, registryChainSel uint64) (*P2PSignerEnc, error) {
	p2p, signer, enc, err := kslib.ExtractKeys(n, registryChainSel)
	if err != nil {
		return nil, fmt.Errorf("failed to extract keys: %w", err)
	}
	return &P2PSignerEnc{
		Signer:              signer,
		P2PKey:              p2p,
		EncryptionPublicKey: enc,
	}, nil
}

// UpdateNodeCapabilitiesRequest is a request to set the capabilities of nodes in the registry
type UpdateNodeCapabilitiesRequest = MutateNodeCapabilitiesRequest

// MutateNodeCapabilitiesRequest is a request to change the capabilities of nodes in the registry
type MutateNodeCapabilitiesRequest struct {
	AddressBook      deployment.AddressBook
	RegistryChainSel uint64

	P2pToCapabilities map[p2pkey.PeerID][]kcr.CapabilitiesRegistryCapability
}

func (req *MutateNodeCapabilitiesRequest) Validate() error {
	if req.AddressBook == nil {
		return fmt.Errorf("address book is nil")
	}
	if len(req.P2pToCapabilities) == 0 {
		return fmt.Errorf("p2pToCapabilities is empty")
	}
	_, exists := chainsel.ChainBySelector(req.RegistryChainSel)
	if !exists {
		return fmt.Errorf("registry chain selector %d does not exist", req.RegistryChainSel)
	}

	return nil
}

func (req *MutateNodeCapabilitiesRequest) updateNodeCapabilitiesImplRequest(e deployment.Environment) (*internal.UpdateNodeCapabilitiesImplRequest, error) {
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

	return &internal.UpdateNodeCapabilitiesImplRequest{
		Chain:             registryChain,
		Registry:          registry,
		P2pToCapabilities: req.P2pToCapabilities,
	}, nil
}

// UpdateNodeCapabilities updates the capabilities of nodes in the registry
func UpdateNodeCapabilities(env deployment.Environment, req *MutateNodeCapabilitiesRequest) (deployment.ChangesetOutput, error) {
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
