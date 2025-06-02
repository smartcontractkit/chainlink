package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
)

type GetContractSetsRequest struct {
	Chains      map[uint64]cldf_evm.Chain
	AddressBook cldf.AddressBook

	// Labels indicates the label set that a contract must include to be considered as a member
	// of the returned contract set.  By default, an empty label set implies that only contracts without
	// labels will be considered.  Otherwise, all labels must be on the contract (e.g., "label1" AND "label2").
	Labels []string
}

type GetContractSetsResponse struct {
	ContractSets map[uint64]ContractSet
}

type ContractSet struct {
	commonchangeset.MCMSWithTimelockState
	OCR3                 map[common.Address]*ocr3_capability.OCR3Capability
	Forwarder            *forwarder.KeystoneForwarder
	CapabilitiesRegistry *capabilities_registry.CapabilitiesRegistry
	WorkflowRegistry     *workflow_registry.WorkflowRegistry
}

func (cs ContractSet) Convert() internal.ContractSet {
	return internal.ContractSet{
		MCMSWithTimelockState: commonchangeset.MCMSWithTimelockState{
			MCMSWithTimelockContracts: cs.MCMSWithTimelockContracts,
		},
		Forwarder:            cs.Forwarder,
		WorkflowRegistry:     cs.WorkflowRegistry,
		OCR3:                 cs.OCR3,
		CapabilitiesRegistry: cs.CapabilitiesRegistry,
	}
}

func (cs ContractSet) TransferableContracts() []common.Address {
	var out []common.Address
	if cs.OCR3 != nil {
		for _, ocr := range cs.OCR3 {
			out = append(out, ocr.Address())
		}
	}
	if cs.Forwarder != nil {
		out = append(out, cs.Forwarder.Address())
	}
	if cs.CapabilitiesRegistry != nil {
		out = append(out, cs.CapabilitiesRegistry.Address())
	}
	if cs.WorkflowRegistry != nil {
		out = append(out, cs.WorkflowRegistry.Address())
	}
	return out
}

func (cs ContractSet) GetOCR3Contract(addr *common.Address) (*ocr3_capability.OCR3Capability, error) {
	return getOCR3Contract(cs.OCR3, addr)
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

		// Forwarder addresses now have informative labels, but we don't want them to be ignored if no labels are provided for filtering.
		// If labels are provided, just filter by those.
		forwarderAddrs := make(map[string]cldf.TypeAndVersion)
		if len(req.Labels) == 0 {
			for addr, tv := range addrs {
				if tv.Type == KeystoneForwarder {
					forwarderAddrs[addr] = tv
				}
			}
		}

		// TODO: we need to expand/refactor the way labeled addresses are filtered
		// see: https://smartcontract-it.atlassian.net/browse/CRE-363
		filtered := deployment.LabeledAddresses(addrs).And(req.Labels...)

		for addr, tv := range forwarderAddrs {
			filtered[addr] = tv
		}

		cs, err := loadContractSet(lggr, chain, filtered)
		if err != nil {
			return nil, fmt.Errorf("failed to load contract set for chain %d: %w", id, err)
		}
		resp.ContractSets[id] = *cs
	}
	return resp, nil
}

func loadContractSet(lggr logger.Logger, chain cldf_evm.Chain, addresses map[string]cldf.TypeAndVersion) (*ContractSet, error) {
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
