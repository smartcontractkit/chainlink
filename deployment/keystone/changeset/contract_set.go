package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink/deployment"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

// contractSetV2 represents a set of contracts for a specific chain.
// TODO: remove this once migration to DataStore is complete.
type contractSetV2 struct {
	OCR3                 map[common.Address]OwnedContract[*ocr3_capability.OCR3Capability]
	Forwarder            *OwnedContract[*forwarder.KeystoneForwarder]
	CapabilitiesRegistry *OwnedContract[*capabilities_registry.CapabilitiesRegistry]
	WorkflowRegistry     *OwnedContract[*workflow_registry.WorkflowRegistry]
}

// getContractSetsRequestV2 is the request structure for getting contract sets.
type getContractSetsRequestV2 struct {
	AddressBook deployment.AddressBook
	Chains      map[uint64]deployment.Chain
	Labels      []string
}

// getContractSetsResponseV2 is the response structure for getting contract sets.
type getContractSetsResponseV2 struct {
	ContractSets map[uint64]contractSetV2
}

func (cs contractSetV2) toContractSet() ContractSet {
	// We don't set the `ContractSet.MCMSWithTimelockState` due to the nature of `ContractSetV2`.
	// Now each contract has its own MCMS state, and that format is not compatible with `ContractSet`.
	out := ContractSet{}

	ocr3Contracts := make(map[common.Address]*ocr3_capability.OCR3Capability)
	for addr, ocr := range cs.OCR3 {
		ocr3Contracts[addr] = ocr.Contract
	}

	out.OCR3 = ocr3Contracts
	if cs.Forwarder != nil {
		out.Forwarder = cs.Forwarder.Contract
	}
	if cs.CapabilitiesRegistry != nil {
		out.CapabilitiesRegistry = cs.CapabilitiesRegistry.Contract
	}
	if cs.WorkflowRegistry != nil {
		out.WorkflowRegistry = cs.WorkflowRegistry.Contract
	}

	return out
}

// transferableContracts returns a list of addresses of contracts that are transferable.
func (cs contractSetV2) transferableContracts() []common.Address {
	return cs.toContractSet().TransferableContracts()
}

// getContractSetsV2 retrieves the contract sets for the given chains and labels.
func getContractSetsV2(lggr logger.Logger, req getContractSetsRequestV2) (*getContractSetsResponseV2, error) {
	out := &getContractSetsResponseV2{
		ContractSets: make(map[uint64]contractSetV2),
	}

	for id, chain := range req.Chains {
		addresses, err := req.AddressBook.AddressesForChain(id)
		if err != nil {
			return nil, fmt.Errorf("failed to get addresses for chain %d: %w", id, err)
		}

		// Forwarder addresses now have informative labels, but we don't want them to be ignored if no labels are provided for filtering.
		// If labels are provided, just filter by those.
		forwarderAddrs := make(map[string]deployment.TypeAndVersion)
		if len(req.Labels) == 0 {
			for addr, tv := range addresses {
				if tv.Type == KeystoneForwarder {
					forwarderAddrs[addr] = tv
				}
			}
		}

		// TODO: we need to expand/refactor the way labeled addresses are filtered
		// see: https://smartcontract-it.atlassian.net/browse/CRE-363
		filtered := deployment.LabeledAddresses(addresses).And(req.Labels...)
		for addr, tv := range forwarderAddrs {
			filtered[addr] = tv
		}

		cs, err := loadContractSetV2(lggr, req.AddressBook, chain, filtered)
		if err != nil {
			return nil, fmt.Errorf("error loading contract set for chain %s: %w", chain.Name(), err)
		}

		out.ContractSets[id] = *cs
	}

	return out, nil
}

func loadContractSetV2(lggr logger.Logger, addressBook deployment.AddressBook, chain deployment.Chain, addresses map[string]deployment.TypeAndVersion) (*contractSetV2, error) {
	var out contractSetV2

	handlers := map[deployment.ContractType]func(string) error{
		OCR3Capability: func(addr string) error {
			contract, err := GetOwnedContract[*ocr3_capability.OCR3Capability](addressBook, chain, addr)
			if err != nil {
				return err
			}
			if out.OCR3 == nil {
				out.OCR3 = make(map[common.Address]OwnedContract[*ocr3_capability.OCR3Capability])
			}
			out.OCR3[common.HexToAddress(addr)] = *contract
			return nil
		},
		KeystoneForwarder: func(addr string) error {
			contract, err := GetOwnedContract[*forwarder.KeystoneForwarder](addressBook, chain, addr)
			if err != nil {
				return err
			}
			out.Forwarder = contract
			return nil
		},
		CapabilitiesRegistry: func(addr string) error {
			contract, err := GetOwnedContract[*capabilities_registry.CapabilitiesRegistry](addressBook, chain, addr)
			if err != nil {
				return err
			}
			out.CapabilitiesRegistry = contract
			return nil
		},
		WorkflowRegistry: func(addr string) error {
			contract, err := GetOwnedContract[*workflow_registry.WorkflowRegistry](addressBook, chain, addr)
			if err != nil {
				return err
			}
			out.WorkflowRegistry = contract
			return nil
		},
	}

	for addr, tv := range addresses {
		handler, exists := handlers[tv.Type]
		if !exists {
			// ignore unknown contract types
			lggr.Warnf("Unknown contract type %s for address %s", tv.Type, addr)
			continue
		}

		if err := handler(addr); err != nil {
			return nil, fmt.Errorf("error trying to load contract type `%s` for address %s: %w", tv.Type, addr, err)
		}
	}

	return &out, nil
}
