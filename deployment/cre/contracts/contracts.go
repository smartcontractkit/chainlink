package contracts

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	shard_config "github.com/smartcontractkit/chainlink-evm/contracts/cre/gobindings/shardconfig/generated/v1_0_0/shard_config"
	capabilities_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	capabilities_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/capabilities_registry_wrapper_v2"
	workflow_registry "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v1"
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

var (
	CapabilitiesRegistry      cldf.ContractType = "CapabilitiesRegistry"      // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/CapabilitiesRegistry.sol#L392
	WorkflowRegistry          cldf.ContractType = "WorkflowRegistry"          // https://github.com/smartcontractkit/chainlink/blob/develop/contracts/src/v0.8/workflow/WorkflowRegistry.sol
	KeystoneForwarder         cldf.ContractType = "KeystoneForwarder"         // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/KeystoneForwarder.sol#L90
	OCR3Capability            cldf.ContractType = "OCR3Capability"            // https://github.com/smartcontractkit/chainlink/blob/50c1b3dbf31bd145b312739b08967600a5c67f30/contracts/src/v0.8/keystone/OCR3Capability.sol#L12
	BalanceReader             cldf.ContractType = "BalanceReader"             // https://github.com/smartcontractkit/chainlink-evm/blob/2724ef8937488de77b320e4e9692ed0dcb3a165a/contracts/src/v0.8/keystone/BalanceReader.sol
	FeedConsumer              cldf.ContractType = "FeedConsumer"              // no type and a version in contract https://github.com/smartcontractkit/chainlink/blob/89183a8a5d22b1aeca0ade3b76d16aa84067aa57/contracts/src/v0.8/keystone/KeystoneFeedsConsumer.sol#L1
	RBACTimelock              cldf.ContractType = "RBACTimelock"              // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/RBACTimelock.sol
	ProposerManyChainMultiSig cldf.ContractType = "ProposerManyChainMultiSig" // no type and a version in contract https://github.com/smartcontractkit/ccip-owner-contracts/blob/main/src/ManyChainMultiSig.sol
	ShardConfig               cldf.ContractType = "ShardConfig"               // manages desired shard count configuration
)

type MCMSConfig = proposalutils.TimelockConfig

// Ownable is an interface for contracts that have an owner.
type Ownable interface {
	Address() common.Address
	Owner(opts *bind.CallOpts) (common.Address, error)
}

func isOwnedByMCMSV2[T Ownable](contract T, store datastore.AddressRefStore, chain cldf_evm.Chain) (bool, error) {
	var timelockTV = cldf.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)

	r, err := getOwnerReference(contract, store, chain)
	if err != nil {
		return false, fmt.Errorf("failed to get owner reference: %w", err)
	}

	if r != nil && cldf.ContractType(r.Type) == timelockTV.Type && r.Version.String() == timelockTV.Version.String() {
		return true, nil
	}

	return false, nil
}

// OwnedContract represents a contract and its owned MCMS contracts.
type OwnedContract[T Ownable] struct {
	// The MCMS contracts that the contract might own
	McmsContracts *state.MCMSWithTimelockState
	// The actual contract instance
	Contract T
}

// NewOwnable creates an OwnedContract instance.
// It checks if the contract is owned by a timelock contract and loads the MCMS state if necessary.
func NewOwnableV2[T Ownable](contract T, store datastore.AddressRefStore, chain cldf_evm.Chain) (*OwnedContract[T], error) {
	isOwnedByMCMSV2, err := isOwnedByMCMSV2[T](contract, store, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to check if contract is owned by MCMS: %w", err)
	}
	if !isOwnedByMCMSV2 {
		return &OwnedContract[T]{
			McmsContracts: nil,
			Contract:      contract,
		}, nil
	}
	// find all the addresses by the chain and qualifier that match the timelock
	// make sure they constitute a valid MCMS with timelock
	r, err := getOwnerReference(contract, store, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner reference: %w", err)
	}
	// in the latest versions, the qualifier should be the same for all the mcms contracts
	// which enables multiple MCMS deployments on a single chain
	stateMCMS, err := state.GetMCMSWithTimelockState(store, chain, r.Qualifier)
	if err != nil {
		return nil, fmt.Errorf("failed to get MCMS with timelock state: %w", err)
	}

	if err := stateMCMS.Validate(); err != nil {
		// older versions had adhoc qualifiers, so we have to try labels sets
		// TODO CRE-1360: remove this after we complete migration to consistent qualifiers
		m := matchLabels(store, *r, chain.Selector)
		var err2 error
		stateMCMS, err2 = state.MaybeLoadMCMSWithTimelockChainState(chain, m)
		if err2 != nil {
			return nil, fmt.Errorf("failed to get MCMS with timelock state by labels: %w", err2)
		}
		err2 = stateMCMS.Validate()
		if err2 != nil {
			return nil, fmt.Errorf("failed to validate MCMS with timelock state by labels: %w", err2)
		}
	}

	return &OwnedContract[T]{
		McmsContracts: stateMCMS,
		Contract:      contract,
	}, nil
}

func matchLabels(ab datastore.AddressRefStore, ref datastore.AddressRef, chainSelector uint64) map[string]cldf.TypeAndVersion {
	addresses := ab.Filter(datastore.AddressRefByChainSelector(chainSelector))
	addressesMap := make(map[string]cldf.TypeAndVersion)
	for _, addr := range addresses {
		if !ref.Labels.Equal(addr.Labels) {
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

// GetOwnerTypeAndVersionV2 retrieves the owner type and version of a contract using the datastore instead of the address book.
func GetOwnerTypeAndVersionV2[T Ownable](contract T, ab datastore.AddressRefStore, chain cldf_evm.Chain) (*cldf.TypeAndVersion, error) {
	// Get the contract owner
	owner, err := contract.Owner(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract owner: %w", err)
	}

	// Look for owner in datastore
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

// getOwnerReference retrieves the owner reference of a contract using the datastore.
// If the owner is not found, it returns nil without an error.
func getOwnerReference[T Ownable](contract T, store datastore.AddressRefStore, chain cldf_evm.Chain) (*datastore.AddressRef, error) {
	// Get the contract owner
	owner, err := contract.Owner(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract owner: %w", err)
	}

	// Look for owner in address book
	addresses := store.Filter(datastore.AddressRefByChainSelector(chain.Selector), datastore.AddressRefByAddress(owner.Hex()))
	if len(addresses) == 1 {
		return &addresses[0], nil
	}
	// If the owner is not found, then we assume it's non-MCMS so return nil
	return nil, nil
}

// GetOwnableContractV2 retrieves a contract instance of type T from the datastore.
// If `targetAddr` is provided, it will look for that specific address.
// If not, it will default to looking one contract of type T, and if it doesn't find exactly one, it will error.
func GetOwnableContractV2[T Ownable](addrs datastore.AddressRefStore, chain cldf_evm.Chain, targetAddr string) (*T, error) {
	// Determine contract type based on T
	switch any(*new(T)).(type) {
	case *forwarder.KeystoneForwarder:
	case *capabilities_registry.CapabilitiesRegistry:
	case *capabilities_registry_v2.CapabilitiesRegistry:
	case *ocr3_capability.OCR3Capability:
	case *workflow_registry.WorkflowRegistry:
	case *workflow_registry_v2.WorkflowRegistry:
	case *shard_config.ShardConfig:
	default:
		return nil, fmt.Errorf("unsupported contract type %T", *new(T))
	}

	addresses := addrs.Filter(datastore.AddressRefByChainSelector(chain.Selector), datastore.AddressRefByAddress(targetAddr))
	if len(addresses) != 1 {
		return nil, fmt.Errorf("expected exactly one address for contract at %s on chain %d, found %d", targetAddr, chain.Selector, len(addresses))
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
	case *capabilities_registry_v2.CapabilitiesRegistry:
		c, e := capabilities_registry_v2.NewCapabilitiesRegistry(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *ocr3_capability.OCR3Capability:
		c, e := ocr3_capability.NewOCR3Capability(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *workflow_registry.WorkflowRegistry:
		c, e := workflow_registry.NewWorkflowRegistry(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *workflow_registry_v2.WorkflowRegistry:
		c, e := workflow_registry_v2.NewWorkflowRegistry(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	case *shard_config.ShardConfig:
		c, e := shard_config.NewShardConfig(common.HexToAddress(addr), chain.Client)
		instance, err = any(c).(T), e
	default:
		return nil, errors.New("unsupported contract type for instance creation")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}

	return &instance, nil
}

// GetOwnedContractV2 retrieves an OwnedContract instance of type T from the datastore for a specific address.
func GetOwnedContractV2[T Ownable](store datastore.AddressRefStore, chain cldf_evm.Chain, addr string) (*OwnedContract[T], error) {
	addresses := store.Filter(datastore.AddressRefByChainSelector(chain.Selector), datastore.AddressRefByAddress(addr))

	if len(addresses) == 0 {
		return nil, fmt.Errorf("address %s not found in address ref store for chain %d", addr, chain.Selector)
	}
	// should not happen since address is unique
	if len(addresses) > 1 {
		return nil, fmt.Errorf("multiple addresses found for %s in address ref store for chain %d", addr, chain.Selector)
	}
	contract, err := GetOwnableContractV2[T](store, chain, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract at %s: %w", addr, err)
	}

	ownedContract, err := NewOwnableV2(*contract, store, chain)
	if err != nil {
		return nil, fmt.Errorf("failed to create owned contract for %s: %w", addr, err)
	}

	return ownedContract, nil
}
