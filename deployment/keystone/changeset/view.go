package changeset

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	commonview "github.com/smartcontractkit/chainlink/deployment/common/view"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ViewStateV2 = ViewKeystone

type contractsPerChain map[uint64]viewContracts

type viewContracts struct {
	OCR3                 map[common.Address]*ocr3_capability.OCR3Capability
	Forwarder            map[common.Address]*forwarder.KeystoneForwarder
	CapabilitiesRegistry map[common.Address]*capabilities_registry.CapabilitiesRegistry
	WorkflowRegistry     map[common.Address]*workflow_registry.WorkflowRegistry
}

func ViewKeystone(e deployment.Environment, previousView json.Marshaler) (json.Marshaler, error) {
	lggr := e.Logger
	contractsMap, err := getContractsPerChain(e)
	// This is an unrecoverable error
	if err != nil {
		return nil, fmt.Errorf("failed to get contract sets: %w", err)
	}

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
	for chainSel, contracts := range contractsMap {
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
		v, err := GenerateKeystoneChainView(e.GetContext(), e.Logger, prevView.Chains[chainName], contracts)
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

func getContractsPerChain(e deployment.Environment) (contractsPerChain, error) {
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

	contracts := make(contractsPerChain)
	var errs error

	// Initialize all contract sets first
	for _, addr := range contractAddresses {
		if _, ok := contracts[addr.ChainSelector]; !ok {
			contracts[addr.ChainSelector] = viewContracts{
				OCR3:                 make(map[common.Address]*ocr3_capability.OCR3Capability),
				Forwarder:            make(map[common.Address]*forwarder.KeystoneForwarder),
				CapabilitiesRegistry: make(map[common.Address]*capabilities_registry.CapabilitiesRegistry),
				WorkflowRegistry:     make(map[common.Address]*workflow_registry.WorkflowRegistry),
			}
		}
	}

	for _, contractAddress := range contractAddresses {
		chain, ok := e.BlockChains.EVMChains()[contractAddress.ChainSelector]
		if !ok {
			// the chain might not be present in the environment if it was removed due to RPC instability
			e.Logger.Warnf("chain with selector %d not found, skipping contract address %s", contractAddress.ChainSelector, contractAddress.Address)
			continue
		}

		// Get a mutable copy of the ContractSet
		set := contracts[contractAddress.ChainSelector]

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
			set.CapabilitiesRegistry[ownedContract.Contract.Address()] = ownedContract.Contract

		case datastore.ContractType(internal.OCR3Capability):
			ownedContract, err := GetOwnedContractV2[*ocr3_capability.OCR3Capability](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve OCR3 capability contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.OCR3[ownedContract.Contract.Address()] = ownedContract.Contract

		case datastore.ContractType(internal.KeystoneForwarder):
			ownedContract, err := GetOwnedContractV2[*forwarder.KeystoneForwarder](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve forwarder contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.Forwarder[ownedContract.Contract.Address()] = ownedContract.Contract

		case datastore.ContractType(internal.WorkflowRegistry):
			ownedContract, err := GetOwnedContractV2[*workflow_registry.WorkflowRegistry](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve workflow registry contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.WorkflowRegistry[ownedContract.Contract.Address()] = ownedContract.Contract
		}

		// Store the updated `contractSet` back in the map
		contracts[contractAddress.ChainSelector] = set
	}

	return contracts, errs
}
