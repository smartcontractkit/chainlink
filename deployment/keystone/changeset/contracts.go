package changeset

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/smartcontractkit/chainlink/deployment"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"

	capabilities_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry_1_1_0"
	forwarder "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/forwarder_1_0_0"
	ocr3_capability "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/ocr3_capability_1_0_0"
	workflow_registry "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/workflow/generated/workflow_registry_wrapper"
)

type Ownable interface {
	Address() common.Address
	Owner(opts *bind.CallOpts) (common.Address, error)
}

// OwnedContract represents a contract and its owned MCMS contracts.
type OwnedContract[T Ownable] struct {
	// The MCMS contracts that the contract might own
	McmsContracts commonchangeset.MCMSWithTimelockState
	// The actual contract instance
	Contract T
}

// NewOwnable creates an OwnedContract instance
func NewOwnable[T Ownable](contract T, ab deployment.AddressBook, chain deployment.Chain) (*OwnedContract[T], error) {
	var (
		// We expect one of each contract on the chain.
		timelock  = deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0)
		callProxy = deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0)
		proposer  = deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0)
		canceller = deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0)
		bypasser  = deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0)

		// the same contract can have different roles
		// multichain    = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		proposerMCMS  = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		bypasserMCMS  = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
		cancellerMCMS = deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
	)

	// Look for MCMS contracts that might be owned by the contract
	addresses, err := ab.AddressesForChain(chain.Selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}

	// Convert map keys to a slice
	proposerMCMS.Labels.Add(types.ProposerRole.String())
	bypasserMCMS.Labels.Add(types.BypasserRole.String())
	cancellerMCMS.Labels.Add(types.CancellerRole.String())
	wantTypes := []deployment.TypeAndVersion{timelock, proposer, canceller, bypasser, callProxy,
		proposerMCMS, bypasserMCMS, cancellerMCMS,
	}

	// Ensure we either have the bundle or not.
	_, err = deployment.EnsureDeduped(addresses, wantTypes)
	if err != nil {
		return nil, fmt.Errorf("unable to check MCMS contracts on chain %s error: %w", chain.Name(), err)
	}

	// Get the contract owner
	owner, err := contract.Owner(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract owner: %w", err)
	}
	fmt.Println("Contract owner:", owner.Hex())

	// Filter for potential MCMS contracts
	// TODO: figure out if the contract is owned by MCMS and only load the state if that's the case.
	// Load MCMS state
	mcmsState, err := commonchangeset.MaybeLoadMCMSWithTimelockChainState(chain, addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to load MCMS state: %w", err)
	}

	return &OwnedContract[T]{
		McmsContracts: *mcmsState,
		Contract:      contract,
	}, nil
}

func GetContract[T Ownable](ab deployment.AddressBook, chain deployment.Chain, targetAddr string) (*T, error) {
	var contractType deployment.ContractType

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
	if len(targetAddr) > 0 {
		tv, exists := addresses[targetAddr]
		if !exists {
			return nil, fmt.Errorf("address %s not found in address book", targetAddr)
		}

		if tv.Type != contractType {
			return nil, fmt.Errorf("address %s is not a %s, got %s", targetAddr, contractType, tv.Type)
		}

		// Create and return the contract instance
		var instance T
		var err error
		switch any(*new(T)).(type) {
		case *forwarder.KeystoneForwarder:
			c, e := forwarder.NewKeystoneForwarder(common.HexToAddress(targetAddr), chain.Client)
			instance, err = any(c).(T), e
		case *capabilities_registry.CapabilitiesRegistry:
			c, e := capabilities_registry.NewCapabilitiesRegistry(common.HexToAddress(targetAddr), chain.Client)
			instance, err = any(c).(T), e
		case *ocr3_capability.OCR3Capability:
			c, e := ocr3_capability.NewOCR3Capability(common.HexToAddress(targetAddr), chain.Client)
			instance, err = any(c).(T), e
		case *workflow_registry.WorkflowRegistry:
			c, e := workflow_registry.NewWorkflowRegistry(common.HexToAddress(targetAddr), chain.Client)
			instance, err = any(c).(T), e
		}

		if err != nil {
			return nil, fmt.Errorf("failed to create contract instance: %w", err)
		}

		return &instance, nil
	}

	// Find first contract of the required type
	for addr, tv := range addresses {
		if tv.Type == contractType {
			// Create and return the contract instance
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
			}

			if err != nil {
				return nil, fmt.Errorf("failed to create contract instance: %w", err)
			}

			return &instance, nil
		}
	}

	return nil, fmt.Errorf("no contract of type %s found", contractType)
}
