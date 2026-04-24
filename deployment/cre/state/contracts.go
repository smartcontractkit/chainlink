package state

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

type viewContracts struct {
	Forwarder            map[common.Address]*forwarder.KeystoneForwarder
	OCR3                 map[common.Address]*ocr3_capability.OCR3Capability
	WorkflowRegistry     map[common.Address]*workflow_registry_v2.WorkflowRegistry
	CapabilitiesRegistry map[common.Address]*capabilities_registry_v2.CapabilitiesRegistry
}

type contractsPerChain map[uint64]viewContracts

func getContractsPerChain(e deployment.Environment) (contractsPerChain, error) {
	// Cannot do a single `Filter` call because it appears to work as an AND filter.
	ocr3CapabilityContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(crecontracts.OCR3Capability)),
	)
	workflowRegistryContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(crecontracts.WorkflowRegistry)),
	)
	keystoneForwarderContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(crecontracts.KeystoneForwarder)),
	)
	capabilitiesRegistryContracts := e.DataStore.Addresses().Filter(
		datastore.AddressRefByType(datastore.ContractType(crecontracts.CapabilitiesRegistry)),
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
				CapabilitiesRegistry: make(map[common.Address]*capabilities_registry_v2.CapabilitiesRegistry),
				WorkflowRegistry:     make(map[common.Address]*workflow_registry_v2.WorkflowRegistry),
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
		case datastore.ContractType(crecontracts.CapabilitiesRegistry):
			ownedContract, err := crecontracts.GetOwnedContractV2[*capabilities_registry_v2.CapabilitiesRegistry](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve capabilities registry contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.CapabilitiesRegistry[ownedContract.Contract.Address()] = ownedContract.Contract

		case datastore.ContractType(crecontracts.OCR3Capability):
			ownedContract, err := crecontracts.GetOwnedContractV2[*ocr3_capability.OCR3Capability](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve OCR3 capability contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.OCR3[ownedContract.Contract.Address()] = ownedContract.Contract

		case datastore.ContractType(crecontracts.KeystoneForwarder):
			ownedContract, err := crecontracts.GetOwnedContractV2[*forwarder.KeystoneForwarder](
				e.DataStore.Addresses(), chain, contractAddress.Address,
			)
			if err != nil {
				errs = errors.Join(errs, fmt.Errorf("failed to retrieve forwarder contract at %s: %w",
					contractAddress.Address, err))
				continue
			}
			set.Forwarder[ownedContract.Contract.Address()] = ownedContract.Contract

		case datastore.ContractType(crecontracts.WorkflowRegistry):
			ownedContract, err := crecontracts.GetOwnedContractV2[*workflow_registry_v2.WorkflowRegistry](
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
