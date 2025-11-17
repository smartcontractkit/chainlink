package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"

	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/contracts"
)

const (
	SetConfigDescription               = "setConfig on workflow registry v2"
	UpdateAllowedSignersDescription    = "updateAllowedSigners on workflow registry v2"
	SetWorkflowOwnerConfigDescription  = "setWorkflowOwner config on workflow registry v2"
	SetDONLimitDescription             = "setDonLimit on workflow registry v2"
	SetUserDONOverrideDescription      = "setUserDonOverride on workflow registry v2"
	SetCapabilitiesRegistryDescription = "setCapabilitiesRegistry on workflow registry v2"
)

// Common dependencies for workflow registry operations
type WorkflowRegistryOpDeps struct {
	Env      *cldf.Environment
	Strategy strategies.TransactionStrategy
	Registry *workflow_registry_v2.WorkflowRegistry
}

// SetConfig Operation
type SetConfigOpInput struct {
	// We are passing the registry via the deps, but we keep chainSelector and qualifier to allow the operation to be
	// unique.
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"`

	NameLen    uint8                 `json:"nameLen"`
	TagLen     uint8                 `json:"tagLen"`
	URLLen     uint8                 `json:"urlLen"`
	AttrLen    uint16                `json:"attrLen"`
	ExpiryLen  uint32                `json:"expiryLen"`
	MCMSConfig *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetConfigOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var SetConfigOp = operations.NewOperation(
	"set-metadata-config-op",
	semver.MustParse("1.0.0"),
	"Set Config in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetConfigOpInput) (SetConfigOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.SetConfig(opts, input.NameLen, input.TagLen, input.URLLen, input.AttrLen, input.ExpiryLen)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return SetConfigOpOutput{}, fmt.Errorf("failed to execute SetConfig: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetConfig on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set metadata config on chain %d", input.ChainSelector)
		}

		return SetConfigOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// UpdateAllowedSigners Operation
type UpdateAllowedSignersOpInput struct {
	// We are passing the registry via the deps, but we keep chainSelector and qualifier to allow the operation to be
	// unique.
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"` // Qualifier to identify the specific workflow registry

	Signers    []common.Address      `json:"signers"`
	Allowed    bool                  `json:"allowed"`
	MCMSConfig *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type UpdateAllowedSignersOpOutput struct {
	Success         bool                      `json:"success"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
	RegistryAddress common.Address            `json:"registryAddress"`
}

var UpdateAllowedSignersOp = operations.NewOperation(
	"update-allowed-signers-op",
	semver.MustParse("1.0.0"),
	"Update Allowed Signers in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UpdateAllowedSignersOpInput) (UpdateAllowedSignersOpOutput, error) {
		if len(input.Signers) == 0 {
			return UpdateAllowedSignersOpOutput{}, errors.New("must provide at least one signer")
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.UpdateAllowedSigners(opts, input.Signers, input.Allowed)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return UpdateAllowedSignersOpOutput{}, fmt.Errorf("failed to execute UpdateAllowedSigners: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UpdateAllowedSigners on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully updated allowed signers on chain %d", input.ChainSelector)
		}

		return UpdateAllowedSignersOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// SetWorkflowOwnerConfig Operation
type SetWorkflowOwnerConfigOpInput struct {
	// We are passing the registry via the deps, but we keep chainSelector and qualifier to allow the operation to be
	// unique.
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"` // Qualifier to identify the specific workflow registry

	Owner      common.Address        `json:"owner"`
	Config     []byte                `json:"config"`
	MCMSConfig *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetWorkflowOwnerConfigOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var SetWorkflowOwnerConfigOp = operations.NewOperation(
	"set-workflow-owner-config-op",
	semver.MustParse("1.0.0"),
	"Set Workflow Owner Config in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetWorkflowOwnerConfigOpInput) (SetWorkflowOwnerConfigOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.SetWorkflowOwnerConfig(opts, input.Owner, input.Config)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return SetWorkflowOwnerConfigOpOutput{}, fmt.Errorf("failed to execute SetWorkflowOwnerConfig: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetWorkflowOwnerConfig on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set workflow owner config for %s on chain %d", input.Owner.String(), input.ChainSelector)
		}

		return SetWorkflowOwnerConfigOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// SetDONLimit Operation
type SetDONLimitOpInput struct {
	// We are passing the registry via the deps, but we keep chainSelector and qualifier to allow the operation to be
	// unique.
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"` // Qualifier to identify the specific workflow registry

	DONFamily        string                `json:"donFamily"`
	DONLimit         uint32                `json:"donlimit"`
	UserDefaultLimit uint32                `json:"userDefaultLimit"`
	MCMSConfig       *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetDONLimitOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var SetDONLimitOp = operations.NewOperation(
	"set-don-limit-op",
	semver.MustParse("1.0.0"),
	"Set DON DONLimit in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetDONLimitOpInput) (SetDONLimitOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.SetDONLimit(opts, input.DONFamily, input.DONLimit, input.UserDefaultLimit)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return SetDONLimitOpOutput{}, fmt.Errorf("failed to execute SetDONLimit: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetDONLimit on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set DON limit for family %s on chain %d", input.DONFamily, input.ChainSelector)
		}

		return SetDONLimitOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// SetUserDONOverride Operation
type SetUserDONOverrideOpInput struct {
	// We are passing the registry via the deps, but we keep chainSelector and qualifier to allow the operation to be
	// unique.
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"` // Qualifier to identify the specific workflow registry

	User       common.Address        `json:"user"`
	DONFamily  string                `json:"donFamily"`
	Limit      uint32                `json:"limit"`
	Enabled    bool                  `json:"enabled"`
	MCMSConfig *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetUserDONOverrideOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var SetUserDONOverrideOp = operations.NewOperation(
	"set-user-don-override-op",
	semver.MustParse("1.0.0"),
	"Set User DON Override in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetUserDONOverrideOpInput) (SetUserDONOverrideOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.SetUserDONOverride(opts, input.User, input.DONFamily, input.Limit, input.Enabled)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return SetUserDONOverrideOpOutput{}, fmt.Errorf("failed to execute SetUserDONOverride: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetUserDONOverride on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set user DON override for %s on chain %d", input.User.String(), input.ChainSelector)
		}

		return SetUserDONOverrideOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// SetCapabilitiesRegistry MCMSOperation
type SetCapabilitiesRegistryOpInput struct {
	// We are passing the registry via the deps, but we keep chainSelector and qualifier to allow the operation to be
	// unique.
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"` // Qualifier to identify the specific workflow registry

	Registry         common.Address        `json:"registry"`
	ChainSelectorDON uint64                `json:"chainSelectorDON"`
	MCMSConfig       *contracts.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetCapabilitiesRegistryOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var SetCapabilitiesRegistryOp = operations.NewOperation(
	"set-capabilities-registry-op",
	semver.MustParse("1.0.0"),
	"Set DON Registry in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetCapabilitiesRegistryOpInput) (SetCapabilitiesRegistryOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.SetCapabilitiesRegistry(opts, input.Registry, input.ChainSelectorDON)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return SetCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to execute SetCapabilitiesRegistry: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetCapabilitiesRegistry on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set DON registry %s on chain %d", input.Registry.String(), input.ChainSelector)
		}

		return SetCapabilitiesRegistryOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// Helper function to get registry instance from datastore

func GetWorkflowRegistryV2FromDatastore(env *cldf.Environment, chainSelector uint64, qualifier string) (*workflow_registry_v2.WorkflowRegistry, error) {
	addresses := env.DataStore.Addresses().Filter(datastore.AddressRefByChainSelector(chainSelector))
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no addresses found for chain selector %d", chainSelector)
	}

	var registryAddr common.Address
	found := false
	for _, addr := range addresses {
		if addr.Type == "WorkflowRegistry" && addr.Qualifier == qualifier {
			registryAddr = common.HexToAddress(addr.Address)
			found = true
			env.Logger.Infof("Found WorkflowRegistry at %s with qualifier %s", addr.Address, qualifier)
			break
		}
	}

	if !found {
		// Debug: log all available addresses for troubleshooting
		env.Logger.Infof("Available addresses for chain %d:", chainSelector)
		for _, addr := range addresses {
			env.Logger.Infof("  Type: %s, Address: %s, Qualifier: %s", string(addr.Type), addr.Address, addr.Qualifier)
		}
		return nil, fmt.Errorf("workflow registry address not found for chain selector %d and qualifier %s", chainSelector, qualifier)
	}

	chain, ok := env.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return nil, fmt.Errorf("chain with selector %d not found", chainSelector)
	}

	registry, err := workflow_registry_v2.NewWorkflowRegistry(registryAddr, chain.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow registry v2 instance: %w", err)
	}

	return registry, nil
}
