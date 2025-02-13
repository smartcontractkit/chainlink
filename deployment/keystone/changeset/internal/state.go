package internal

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

// ContractSet is a set of contracts for a single chain
// It is a mirror of changeset.ContractSet, and acts an an adapter to the internal package
//
// TODO: remove after CRE-227
type ContractSet struct {
	commonchangeset.MCMSWithTimelockState
	OCR3                 map[common.Address]*ocr3_capability.OCR3Capability
	Forwarder            *forwarder.KeystoneForwarder
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
	WorkflowRegistry     *workflow_registry.WorkflowRegistry
}

type getContractSetsRequest struct {
	Chains      map[uint64]deployment.Chain
	AddressBook deployment.AddressBook
}

type getContractSetsResponse struct {
	ContractSets map[uint64]ContractSet
}

func (cs ContractSet) getOCR3Contract(addr *common.Address) (*ocr3_capability.OCR3Capability, error) {
	return getOCR3Contract(cs.OCR3, addr)
}

func getContractSets(lggr logger.Logger, req *getContractSetsRequest) (*getContractSetsResponse, error) {
	resp := &getContractSetsResponse{
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
	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	out.MCMSWithTimelockState = *mcmsWithTimelock

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
			if out.OCR3 == nil {
				out.OCR3 = make(map[common.Address]*ocr3_capability.OCR3Capability)
			}
			out.OCR3[common.HexToAddress(addr)] = c
		case WorkflowRegistry:
			c, err := workflow_registry.NewWorkflowRegistry(common.HexToAddress(addr), chain.Client)
			if err != nil {
				return nil, fmt.Errorf("failed to create OCR3Capability contract from address %s: %w", addr, err)
			}
			out.WorkflowRegistry = c
		default:
			lggr.Warnw("unknown contract type", "type", tv.Type)
			// ignore unknown contract types
		}
	}
	return &out, nil
}

// getOCR3Contract returns the OCR3 contract from the contract set.  By default, it returns the only
// contract in the set if there is no address specified.  If an address is specified, it returns the
// contract with that address.  If the address is specified but not found in the contract set, it returns
// an error.
func getOCR3Contract(contracts map[common.Address]*ocr3_capability.OCR3Capability, addr *common.Address) (*ocr3_capability.OCR3Capability, error) {
	// Fail if the OCR3 contract address is unspecified and there are multiple OCR3 contracts
	if addr == nil && len(contracts) > 1 {
		return nil, errors.New("OCR contract address is unspecified")
	}

	// Use the first OCR3 contract if the address is unspecified
	if addr == nil && len(contracts) == 1 {
		// use the first OCR3 contract
		for _, c := range contracts {
			return c, nil
		}
	}

	// Select the OCR3 contract by address
	if contract, ok := contracts[*addr]; ok {
		return contract, nil
	}

	addrSet := make([]string, 0, len(contracts))
	for a := range contracts {
		addrSet = append(addrSet, a.String())
	}

	// Fail if the OCR3 contract address is specified but not found in the contract set
	return nil, fmt.Errorf("OCR3 contract address %s not found in contract set %v", *addr, addrSet)
}
