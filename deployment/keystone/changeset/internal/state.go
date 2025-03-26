package internal

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

type GetContractSetsRequest struct {
	Chains      map[uint64]deployment.Chain
	AddressBook deployment.AddressBook

	// Labels indicates the label set that a contract must include to be considered as a member
	// of the returned contract set.  By default, an empty label set implies that only contracts without
	// labels will be considered.  Otherwise, all labels must be on the contract (e.g., "label1" AND "label2").
	Labels []string
}

type GetContractSetsResponse struct {
	ContractSets map[uint64]ContractSet
}

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

func (cs ContractSet) getOCR3Contract(addr *common.Address) (*ocr3_capability.OCR3Capability, error) {
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
		forwarderAddrs := make(map[string]deployment.TypeAndVersion)
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

// loadContractSet loads the MCMS state and then sets the Keystone contract state.
func loadContractSet(
	lggr logger.Logger,
	chain deployment.Chain,
	addresses map[string]deployment.TypeAndVersion,
) (*ContractSet, error) {
	var out ContractSet
	mcmsWithTimelock, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load mcms contract: %w", err)
	}
	out.MCMSWithTimelockState = *mcmsWithTimelock

	if err := setContracts(lggr, addresses, chain.Client, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

// setContracts sets the Keystone contract state.  Non-Keystone contracts (e.g., MCMS contracts) are
// ignored.
func setContracts(
	lggr logger.Logger,
	addresses map[string]deployment.TypeAndVersion,
	client deployment.OnchainClient,
	set *ContractSet,
) error {
	for addr, tv := range addresses {
		// todo handle versions
		switch tv.Type {
		case CapabilitiesRegistry:
			c, err := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(addr), client)
			if err != nil {
				return fmt.Errorf("failed to create capability registry contract from address %s: %w", addr, err)
			}
			set.CapabilitiesRegistry = c
		case KeystoneForwarder:
			c, err := forwarder.NewKeystoneForwarder(common.HexToAddress(addr), client)
			if err != nil {
				return fmt.Errorf("failed to create forwarder contract from address %s: %w", addr, err)
			}
			set.Forwarder = c
		case OCR3Capability:
			c, err := ocr3_capability.NewOCR3Capability(common.HexToAddress(addr), client)
			if err != nil {
				return fmt.Errorf("failed to create OCR3Capability contract from address %s: %w", addr, err)
			}
			if set.OCR3 == nil {
				set.OCR3 = make(map[common.Address]*ocr3_capability.OCR3Capability)
			}
			set.OCR3[common.HexToAddress(addr)] = c
		case WorkflowRegistry:
			c, err := workflow_registry.NewWorkflowRegistry(common.HexToAddress(addr), client)
			if err != nil {
				return fmt.Errorf("failed to create OCR3Capability contract from address %s: %w", addr, err)
			}
			set.WorkflowRegistry = c
		default:
			// do nothing, non-exhaustive
			lggr.Warnf("skipping contract of type : %s", tv.Type)
		}
	}
	return nil
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
