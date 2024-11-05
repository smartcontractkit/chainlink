package keystone

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"
	common_v1_0 "github.com/smartcontractkit/chainlink/deployment/common/view/v1_0"
	"github.com/smartcontractkit/chainlink/deployment/keystone/view"
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

func (cs ContractSet) View() (view.KeystoneChainView, error) {
	out := view.NewKeystoneChainView()
	if cs.CapabilitiesRegistry != nil {
		capRegView, err := common_v1_0.GenerateCapabilityRegistryView(cs.CapabilitiesRegistry)
		if err != nil {
			return view.KeystoneChainView{}, err
		}
		out.CapabilityRegistry[cs.CapabilitiesRegistry.Address().String()] = capRegView
	}
	return out, nil
}

func GetContractSets(lggr logger.Logger, req *GetContractSetsRequest) (*GetContractSetsResponse, error) {
	resp := &GetContractSetsResponse{
		ContractSets: make(map[uint64]ContractSet),
	}
	for id, chain := range req.Chains {
		addrs, err := req.AddressBook.AddressesForChain(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for chain %d: %w", id, err)
		}
		cs, err := loadContractSet(lggr, chain, addrs)
		if err != nil {
			return nil, fmt.Errorf("failed to load contract set for chain %d: %w", id, err)
		}
		resp.ContractSets[id] = *cs
	}
	return resp, nil
}

func loadContractSet(lggr logger.Logger, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*ContractSet, error) {
	var out ContractSet

	for addr, tv := range addresses {
		// todo handle versions
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
			lggr.Warnw("unknown contract type", "type", tv.Type)
			// ignore unknown contract types
		}
	}
	return &out, nil
}
