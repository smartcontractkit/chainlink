package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
)

// Ownable is an interface for contracts that have an owner.
type Ownable interface {
	Address() common.Address
	Owner(opts *bind.CallOpts) (common.Address, error)
}

// OwnedContract represents a contract and its owned MCMS contracts.
type OwnedContract[T Ownable] struct {
	// The MCMS contracts that the contract might own
	McmsContracts *commonchangeset.MCMSWithTimelockState
	// The actual contract instance
	Contract T
}

// NewOwnable creates an OwnedContract instance.
// It checks if the contract is owned by a timelock contract and loads the MCMS state if necessary.
func NewOwnable[T Ownable](contract T, ab cldf.AddressBook, chain cldf_evm.Chain) (*OwnedContract[T], error) {
	var timelockTV = cldf.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)

	// Look for MCMS contracts that might be owned by the contract
	addresses, err := ab.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	ownerTV, err := GetOwnerTypeAndVersion[T](contract, ab, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner type and version: %w", err)
	}

	// Check if the owner is a timelock contract (owned by MCMS)
	// If the owner is not in the address book (ownerTV = nil and err = nil), we assume it's not owned by MCMS
	if ownerTV != nil && ownerTV.Type == timelockTV.Type && ownerTV.Version.String() == timelockTV.Version.String() {
		// Load MCMS state
		stateMCMS, mcmsErr := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
		if mcmsErr != nil {
			return nil, fmt.Errorf("failed to load MCMS state: %w", mcmsErr)
		}

		return &OwnedContract[T]{
			McmsContracts: stateMCMS,
			Contract:      contract,
		}, nil
	}

	return &OwnedContract[T]{
		McmsContracts: nil,
		Contract:      contract,
	}, nil
}

// NewOwnable creates an OwnedContract instance.
// It checks if the contract is owned by a timelock contract and loads the MCMS state if necessary.
func NewOwnableV2[T Ownable](contract T, ab datastore.AddressRefStore, chain cldf_evm.Chain) (*OwnedContract[T], error) {
	var timelockTV = cldf.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)

	ownerTV, err := GetOwnerTypeAndVersionV2[T](contract, ab, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner type and version: %w", err)
	}

	// Check if the owner is a timelock contract (owned by MCMS)
	// If the owner is not in the address book (ownerTV = nil and err = nil), we assume it's not owned by MCMS
	if ownerTV != nil && ownerTV.Type == timelockTV.Type && ownerTV.Version.String() == timelockTV.Version.String() {
		addressesMap := matchLabels(ab, *ownerTV, chain.Selector)
		stateMCMS, mcmsErr := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addressesMap)
		if mcmsErr != nil {
			return nil, fmt.Errorf("failed to load MCMS state: %w", mcmsErr)
		}

		return &OwnedContract[T]{
			McmsContracts: stateMCMS,
			Contract:      contract,
		}, nil
	}

	return &OwnedContract[T]{
		McmsContracts: nil,
		Contract:      contract,
	}, nil
}

func matchLabels(ab datastore.AddressRefStore, tv cldf.TypeAndVersion, chainSelector uint64) map[string]cldf.TypeAndVersion {
	addresses := ab.Filter(datastore.AddressRefByChainSelector(chainSelector))
	addressesMap := make(map[string]cldf.TypeAndVersion)
	for _, addr := range addresses {
		if !tv.Labels.Equal(cldf.NewLabelSet(addr.Labels.List()...)) {
			continue
		}
		addressesMap[addr.Address] = cldf.TypeAndVersion{
			Type:    cldf.ContractType(addr.Type),
			Version: *addr.Version,
			Labels:  cldf.NewLabelSet(addr.Labels.List()...),
		}
	}
	return addressesMap
}

// GetOwnerTypeAndVersion retrieves the owner type and version of a contract.
func GetOwnerTypeAndVersion[T Ownable](contract T, ab cldf.AddressBook, chain cldf_evm.Chain) (*cldf.TypeAndVersion, error) {
	// Get the contract owner
	owner, err := contract.Owner(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract owner: %w", err)
	}

	// Look for owner in address book
	addresses, err := ab.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chain.Selector, err)
	}

	ownerStr := owner.Hex()
	// Check if owner matches any address in the address book
	if tv, exists := addresses[ownerStr]; exists {
		return &tv, nil
	}

	// Handle case where owner is not in address book
	// Check for case-insensitive match since some addresses might be stored with different casing
	for addr, tv := range addresses {
		if common.HexToAddress(addr) == owner {
			return &tv, nil
		}
	}

	// Owner not found, assume it's non-MCMS so no error is returned
	return nil, nil
}

// GetOwnerTypeAndVersionV2 retrieves the owner type and version of a contract using the datastore instead of the address book.
func GetOwnerTypeAndVersionV2[T Ownable](contract T, ab datastore.AddressRefStore, chain cldf_evm.Chain) (*cldf.TypeAndVersion, error) {
	// Get the contract owner
	owner, err := contract.Owner(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract owner: %w", err)
	}

	// Look for owner in address book
	addresses := ab.Filter(datastore.AddressRefByChainSelector(chain.Selector))

	// Handle case where owner is not in address book
	// Check for case-insensitive match since some addresses might be stored with different casing
	for _, addr := range addresses {
		if common.HexToAddress(addr.Address) == owner {
			return &cldf.TypeAndVersion{
				Type:    cldf.ContractType(addr.Type),
				Version: *addr.Version,
				Labels:  cldf.NewLabelSet(addr.Labels.List()...),
			}, nil
		}
	}

	// Owner not found, assume it's non-MCMS so no error is returned
	return nil, nil
}

// GetOwnableContract retrieves a contract instance of type T from the address book.
// If `targetAddr` is provided, it will look for that specific address.
// If not, it will default to looking one contract of type T, and if it doesn't find exactly one, it will error.
func GetOwnableContract[T Ownable](ab cldf.AddressBook, chain cldf_evm.Chain, targetAddr *string) (*T, error) {
	var contractType cldf.ContractType
	// Determine contract type based on T
	switch any(*new(T)).(type) {
	case *forwarder.KeystoneForwarder:
		contractType = KeystoneForwarder
	case *capabilities_registry.CapabilitiesRegistry:
		contractType = CapabilitiesRegistry
	case *ocr3_capability.OCR3Capability:
		contractType = OCR3Capability
	case *workflow_registry.WorkflowRegistry:
		contractType = WorkflowRegistry
	default:
		return nil, fmt.Errorf("unsupported contract type %T", *new(T))
	}

	addresses, err := ab.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chain.Selector, err)
	}

	// If addr is provided, look for that specific address
	if targetAddr != nil {
		addr := *targetAddr
		tv, exists := addresses[addr]
		if !exists {
			return nil, fmt.Errorf("address %s not found in address book", addr)
		}

		if tv.Type != contractType {
			return nil, fmt.Errorf("address %s is not a %s, got %s", addr, contractType, tv.Type)
		}

		return createContractInstance[T](addr, chain)
	}

	var foundAddr string
	var contractCount int

	for addr, tv := range addresses {
		if tv.Type == contractType {
			contractCount++
			foundAddr = addr

			if contractCount > 1 {
				return nil, fmt.Errorf("multiple contracts of type %s found, must provide a `targetAddr`", contractType)
			}
		}
	}

	if contractCount == 0 {
		return nil, fmt.Errorf("no contract of type %s found", contractType)
	}

	return createContractInstance[T](foundAddr, chain)
}

// GetOwnableContractV2 retrieves a contract instance of type T from the datastore.
// If `targetAddr` is provided, it will look for that specific address.
// If not, it will default to looking one contract of type T, and if it doesn't find exactly one, it will error.
func GetOwnableContractV2[T Ownable](addrs datastore.AddressRefStore, chain cldf_evm.Chain, targetAddr string) (*T, error) {
	// Determine contract type based on T
	switch any(*new(T)).(type) {
	case *forwarder.KeystoneForwarder:
	case *capabilities_registry.CapabilitiesRegistry:
	case *ocr3_capability.OCR3Capability:
	case *workflow_registry.WorkflowRegistry:
	default:
		return nil, fmt.Errorf("unsupported contract type %T", *new(T))
	}

	addresses := addrs.Filter(datastore.AddressRefByChainSelector(chain.Selector))

	var foundAddr bool
	for _, a := range addresses {
		if targetAddr == a.Address {
			foundAddr = true
			break
		}
	}
	if !foundAddr {
		return nil, fmt.Errorf("address %s not found in address book", targetAddr)
	}

	return createContractInstance[T](targetAddr, chain)
}

// createContractInstance is a helper function to create contract instances
func createContractInstance[T Ownable](addr string, chain cldf_evm.Chain) (*T, error) {
	var instance T
	var err error

	switch any(*new(T)).(type) {
	case *forwarder.KeystoneForwarder:
		c, e := forwarder.NewKeystoneForwarder(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *capabilities_registry.CapabilitiesRegistry:
		c, e := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *ocr3_capability.OCR3Capability:
		c, e := ocr3_capability.NewOCR3Capability(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *workflow_registry.WorkflowRegistry:
		c, e := workflow_registry.NewWorkflowRegistry(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	default:
		return nil, errors.New("unsupported contract type for instance creation")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}

	return &instance, nil
}

// GetOwnedContract is a helper function that gets a contract and wraps it in OwnedContract
func GetOwnedContract[T Ownable](addressBook cldf.AddressBook, chain cldf_evm.Chain, addr string) (*OwnedContract[T], error) {
	contract, err := GetOwnableContract[T](addressBook, chain, &addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract at %s: %w", addr, err)
	}

	ownedContract, err := NewOwnable(*contract, addressBook, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to create owned contract for %s: %w", addr, err)
	}

	return ownedContract, nil
}

func GetOwnedContractV2[T Ownable](addrs datastore.AddressRefStore, chain cldf_evm.Chain, addr string) (*OwnedContract[T], error) {
	addresses := addrs.Filter(datastore.AddressRefByChainSelector(chain.Selector))

	var foundAddr bool
	for _, a := range addresses {
		if addr == a.Address {
			foundAddr = true
			break
		}
	}
	if !foundAddr {
		return nil, fmt.Errorf("address %s not found in address book", addr)
	}
	contract, err := GetOwnableContractV2[T](addrs, chain, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract at %s: %w", addr, err)
	}

	ownedContract, err := NewOwnableV2(*contract, addrs, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to create owned contract for %s: %w", addr, err)
	}

	return ownedContract, nil
}

// loadCapabilityRegistry loads the CapabilitiesRegistry contract from the address book or datastore.
func loadCapabilityRegistry(registryChain cldf_evm.Chain, env cldf.Environment, ref datastore.AddressRefKey) (*OwnedContract[*capabilities_registry.CapabilitiesRegistry], error) {
	err := shouldUseDatastore(env, ref)
	if err != nil {
		return nil, fmt.Errorf("failed to check registry ref: %w", err)
	}

	var cr *OwnedContract[*capabilities_registry.CapabilitiesRegistry]

	// `shouldUseDatastore` is already checking for the nil ref, no need to `ref == nil` here
	a, err := env.DataStore.Addresses().Get(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	cr, err = GetOwnedContractV2[*capabilities_registry.CapabilitiesRegistry](env.DataStore.Addresses(), registryChain, a.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned contract: %w", err)
	}

	return cr, nil
}

func getTransferableContracts(addressStore datastore.AddressRefStore, chainSelector uint64) []common.Address {
	var transferableContracts []common.Address

	addresses := addressStore.Filter(datastore.AddressRefByChainSelector(chainSelector))
	for _, addr := range addresses {
		isOCR3Capability := addr.Type == datastore.ContractType(OCR3Capability)
		isWorkflowRegistry := addr.Type == datastore.ContractType(WorkflowRegistry)
		isKeystoneForwarder := addr.Type == datastore.ContractType(KeystoneForwarder)
		isCapabilityRegistry := addr.Type == datastore.ContractType(CapabilitiesRegistry)

		if isCapabilityRegistry || isWorkflowRegistry || isKeystoneForwarder || isOCR3Capability {
			transferableContracts = append(transferableContracts, common.HexToAddress(addr.Address))
		}
	}

	return transferableContracts
}
