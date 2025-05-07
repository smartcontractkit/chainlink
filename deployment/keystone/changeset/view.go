package changeset

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper"

	"github.com/smartcontractkit/chainlink/deployment"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ViewStateV2 = ViewKeystone

func ViewKeystone(e deployment.Environment, previousView json.Marshaler) (json.Marshaler, error) {
	lggr := e.Logger
	contractSets, err := getContractSets(e)
	// This is an unrecoverable error
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}
	fmt.Printf("contract sets: %+v\n", contractSets)

	prevViewBytes, err := previousView.MarshalJSON()
	if err != nil {
		// just log the error, we don't need to stop the execution since the previous view is optional
		lggr.Warnf("failed to marshal previous keystone view: %v", err)
	}
	var prevView KeystoneView
	if len(prevViewBytes) == 0 {
		prevView.Chains = make(map[string]KeystoneChainView)
	} else if err = json.Unmarshal(prevViewBytes, &prevView); err != nil {
		lggr.Warnf("failed to unmarshal previous keystone view: %v", err)
		prevView.Chains = make(map[string]KeystoneChainView)
	}

	var viewErrs error
	chainViews := make(map[string]KeystoneChainView)
	for chainSel, contracts := range contractSets {
		chainid, err := chainsel.ChainIdFromSelector(chainSel)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain id for selector %d: %w", chainSel, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		chainName, err := chainsel.NameFromChainId(chainid)
		if err != nil {
			err2 := fmt.Errorf("failed to resolve chain name for chain id %d: %w", chainid, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			continue
		}
		v, err := contracts.View(e.GetContext(), prevView.Chains[chainName], e.Logger)
		if err != nil {
			err2 := fmt.Errorf("failed to view chain %s: %w", chainName, err)
			lggr.Error(err2)
			viewErrs = errors.Join(viewErrs, err2)
			// don't continue; add the partial view
		}
		chainViews[chainName] = v
	}
	nopsView, err := commonview.GenerateNopsView(e.Logger, e.NodeIDs, e.Offchain)
	if err != nil {
		err2 := fmt.Errorf("failed to view nops: %w", err)
		lggr.Error(err2)
		viewErrs = errors.Join(viewErrs, err2)
	}
	return &KeystoneView{
		Chains: chainViews,
		Nops:   nopsView,
	}, viewErrs
}

func getContractSets(e deployment.Environment) (map[uint64]ContractSet, error) {
	// Cannot do a single `Filter` call because it appears to work as an AND filter.
	ocr3CapabilityContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.OCR3Capability)),
	)
	workflowRegistryContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.WorkflowRegistry)),
	)
	keystoneForwarderContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.KeystoneForwarder)),
	)
	capabilitiesRegistryContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(internal.CapabilitiesRegistry)),
	)
	contractAddresses := make([]datastore.AddressRef, 0, len(ocr3CapabilityContracts)+len(workflowRegistryContracts)+
		len(keystoneForwarderContracts)+len(capabilitiesRegistryContracts))
	contractAddresses = append(contractAddresses, ocr3CapabilityContracts...)
	contractAddresses = append(contractAddresses, workflowRegistryContracts...)
	contractAddresses = append(contractAddresses, keystoneForwarderContracts...)
	contractAddresses = append(contractAddresses, capabilitiesRegistryContracts...)

	contractSets := make(map[uint64]ContractSet)
	var errs error

	// Initialize all contract sets first
	for _, addr := range contractAddresses {
		if _, ok := contractSets[addr.ChainSelector]; !ok {
			contractSets[addr.ChainSelector] = ContractSet{
				OCR3: make(map[common.Address]*ocr3_capability.OCR3Capability),
			}
		}
	}

	// TODO: what happens when there are multiple contracts of the same type deployed on the same chain?
	// As of now, it assigns the last one found to the contract set.
	for _, contractAddress := range contractAddresses {
		chain, ok := e.Chains[contractAddress.ChainSelector]
		if !ok {
			errs = errors.Join(errs, fmt.Errorf("chain with selector %d not found", contractAddress.ChainSelector))
			continue
		}

		// Get a mutable copy of the ContractSet
		set := contractSets[contractAddress.ChainSelector]

		switch contractAddress.Type {
		case datastore.ContractType(internal.CapabilitiesRegistry):
			ownedContract, err := GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve capabilities registry contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.CapabilitiesRegistry = ownedContract.Contract

		case datastore.ContractType(internal.OCR3Capability):
			ownedContract, err := GetOwnedContractV2[*ocr3_capability.OCR3Capability](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve OCR3 capability contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.OCR3[common.HexToAddress(contractAddress.Address)] = ownedContract.Contract

		case datastore.ContractType(internal.KeystoneForwarder):
			ownedContract, err := GetOwnedContractV2[*forwarder.KeystoneForwarder](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve forwarder contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.Forwarder = ownedContract.Contract

		case datastore.ContractType(internal.WorkflowRegistry):
			ownedContract, err := GetOwnedContractV2[*workflow_registry.WorkflowRegistry](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve workflow registry contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.WorkflowRegistry = ownedContract.Contract
		}

		// Store the updated `contractSet` back in the map
		contractSets[contractAddress.ChainSelector] = set
	}

	return contractSets, errs
}
