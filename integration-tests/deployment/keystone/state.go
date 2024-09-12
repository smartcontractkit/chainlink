package keystone

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/integration-tests/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability"
)

type GetContractSetsRequest struct {
	Chains      map[uint64]deployment.Chain
	AddressBook deployment.AddressBook
}

type GetContractSetsResponse struct {
	ContractSets map[uint64]ContractSet
}

type ContractSet struct {
	OCR3                 *ocr3_capability.OCR3Capability
	Forwarder            *forwarder.KeystoneForwarder
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
}

func GetContractSets(req *GetContractSetsRequest) (*GetContractSetsResponse, error) {
	resp := &GetContractSetsResponse{
		ContractSets: make(map[uint64]ContractSet),
	}
	for id, chain := range req.Chains {
		addrs, err := req.AddressBook.AddressesForChain(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for chain %d: %w", id, err)
		}
		cs, err := loadContractSet(chain, addrs)
		if err != nil {
			return nil, fmt.Errorf("failed to load contract set for chain %d: %w", id, err)
		}
		resp.ContractSets[id] = *cs
	}
	return resp, nil
}

func loadContractSet(chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*ContractSet, error) {
	var out ContractSet

	for addr, tv := range addresses {
		// todo handle versions
		if !tv.Version.Equal(&deployment.Version1_0_0) {
			return nil, fmt.Errorf("unsupported version %s", tv.Version.String())
		}
		switch tv.Type {
		case CapabilitiesRegistry:
			c, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create capability registry contract from address %s: %w", addr, err)
			}
			out.CapabilitiesRegistry = c
		case KeystoneForwarder:
			c, err := forwarder.NewKeystoneForwarder(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create forwarder contract from address %s: %w", addr, err)
			}
			out.Forwarder = c
		case OCR3Capability:
			c, err := ocr3_capability.NewOCR3Capability(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create OCR3Capability contract from address %s: %w", addr, err)
			}
			out.OCR3 = c
		default:
			return nil, fmt.Errorf("unknown contract type %s", tv.Type)
		}
	}
	return &out, nil
}
