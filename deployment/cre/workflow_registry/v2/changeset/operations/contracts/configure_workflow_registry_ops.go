package contracts

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
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
	Env           *cldf.Environment
	MCMSContracts *commonchangeset.MCMSWithTimelockState // Required if MCMSConfig is not nil
}

// SetConfig Operation
type SetConfigOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	NameLen       uint8            `json:"nameLen"`
	TagLen        uint8            `json:"tagLen"`
	URLLen        uint8            `json:"urlLen"`
	AttrLen       uint16           `json:"attrLen"`
	ExpiryLen     uint32           `json:"expiryLen"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetConfigOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var SetConfigOp = operations.NewOperation(
	"set-metadata-config-op",
	semver.MustParse("1.0.0"),
	"Set Config in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetConfigOpInput) (SetConfigOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetConfigOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return SetConfigOpOutput{}, fmt.Errorf("failed to get workflow registry v2: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			SetConfigDescription,
		)
		if err != nil {
			return SetConfigOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.SetConfig(opts, input.NameLen, input.TagLen, input.URLLen, input.AttrLen, input.ExpiryLen)
			if err != nil {
				return nil, fmt.Errorf("failed to call SetConfig: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return SetConfigOpOutput{}, fmt.Errorf("failed to execute SetConfig: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetConfig on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set metadata config on chain %d", input.ChainSelector)
		}

		return SetConfigOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// UpdateAllowedSigners Operation
type UpdateAllowedSignersOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	Signers       []common.Address `json:"signers"`
	Allowed       bool             `json:"allowed"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type UpdateAllowedSignersOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var UpdateAllowedSignersOp = operations.NewOperation(
	"update-allowed-signers-op",
	semver.MustParse("1.0.0"),
	"Update Allowed Signers in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input UpdateAllowedSignersOpInput) (UpdateAllowedSignersOpOutput, error) {
		if len(input.Signers) == 0 {
			return UpdateAllowedSignersOpOutput{}, errors.New("must provide at least one signer")
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return UpdateAllowedSignersOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return UpdateAllowedSignersOpOutput{}, fmt.Errorf("failed to get workflow registry v2: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			UpdateAllowedSignersDescription,
		)
		if err != nil {
			return UpdateAllowedSignersOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.UpdateAllowedSigners(opts, input.Signers, input.Allowed)
			if err != nil {
				return nil, fmt.Errorf("failed to call UpdateAllowedSigners: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return UpdateAllowedSignersOpOutput{}, fmt.Errorf("failed to execute UpdateAllowedSigners: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for UpdateAllowedSigners on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully updated allowed signers on chain %d", input.ChainSelector)
		}

		return UpdateAllowedSignersOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// SetWorkflowOwnerConfig Operation
type SetWorkflowOwnerConfigOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	Owner         common.Address   `json:"owner"`
	Config        []byte           `json:"config"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetWorkflowOwnerConfigOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var SetWorkflowOwnerConfigOp = operations.NewOperation(
	"set-workflow-owner-config-op",
	semver.MustParse("1.0.0"),
	"Set Workflow Owner Config in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetWorkflowOwnerConfigOpInput) (SetWorkflowOwnerConfigOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetWorkflowOwnerConfigOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return SetWorkflowOwnerConfigOpOutput{}, fmt.Errorf("failed to get workflow registry v2: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			SetWorkflowOwnerConfigDescription,
		)
		if err != nil {
			return SetWorkflowOwnerConfigOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.SetWorkflowOwnerConfig(opts, input.Owner, input.Config)
			if err != nil {
				return nil, fmt.Errorf("failed to call SetWorkflowOwnerConfig: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return SetWorkflowOwnerConfigOpOutput{}, fmt.Errorf("failed to execute SetWorkflowOwnerConfig: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetWorkflowOwnerConfig on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set workflow owner config for %s on chain %d", input.Owner.String(), input.ChainSelector)
		}

		return SetWorkflowOwnerConfigOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// SetDONLimit Operation
type SetDONLimitOpInput struct {
	ChainSelector    uint64           `json:"chainSelector"`
	Qualifier        string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	DONFamily        string           `json:"donFamily"`
	DONLimit         uint32           `json:"donlimit"`
	UserDefaultLimit uint32           `json:"userDefaultLimit"`
	MCMSConfig       *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetDONLimitOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var SetDONLimitOp = operations.NewOperation(
	"set-don-limit-op",
	semver.MustParse("1.0.0"),
	"Set DON DONLimit in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetDONLimitOpInput) (SetDONLimitOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetDONLimitOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return SetDONLimitOpOutput{}, fmt.Errorf("failed to get workflow registry v2: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			SetDONLimitDescription,
		)
		if err != nil {
			return SetDONLimitOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.SetDONLimit(opts, input.DONFamily, input.DONLimit, input.UserDefaultLimit)
			if err != nil {
				return nil, fmt.Errorf("failed to call SetDONLimit: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return SetDONLimitOpOutput{}, fmt.Errorf("failed to execute SetDONLimit: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetDONLimit on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set DON limit for family %s on chain %d", input.DONFamily, input.ChainSelector)
		}

		return SetDONLimitOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// SetUserDONOverride Operation
type SetUserDONOverrideOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	User          common.Address   `json:"user"`
	DONFamily     string           `json:"donFamily"`
	Limit         uint32           `json:"limit"`
	Enabled       bool             `json:"enabled"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetUserDONOverrideOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var SetUserDONOverrideOp = operations.NewOperation(
	"set-user-don-override-op",
	semver.MustParse("1.0.0"),
	"Set User DON Override in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetUserDONOverrideOpInput) (SetUserDONOverrideOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetUserDONOverrideOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return SetUserDONOverrideOpOutput{}, fmt.Errorf("failed to get workflow registry v2: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			SetUserDONOverrideDescription,
		)
		if err != nil {
			return SetUserDONOverrideOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.SetUserDONOverride(opts, input.User, input.DONFamily, input.Limit, input.Enabled)
			if err != nil {
				return nil, fmt.Errorf("failed to call SetUserDONOverride: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return SetUserDONOverrideOpOutput{}, fmt.Errorf("failed to execute SetUserDONOverride: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetUserDONOverride on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set user DON override for %s on chain %d", input.User.String(), input.ChainSelector)
		}

		return SetUserDONOverrideOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// SetCapabilitiesRegistry Operation
type SetCapabilitiesRegistryOpInput struct {
	ChainSelector    uint64           `json:"chainSelector"`
	Qualifier        string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	Registry         common.Address   `json:"registry"`
	ChainSelectorDON uint64           `json:"chainSelectorDON"`
	MCMSConfig       *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type SetCapabilitiesRegistryOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var SetCapabilitiesRegistryOp = operations.NewOperation(
	"set-capabilities-registry-op",
	semver.MustParse("1.0.0"),
	"Set DON Registry in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input SetCapabilitiesRegistryOpInput) (SetCapabilitiesRegistryOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return SetCapabilitiesRegistryOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return SetCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to get workflow registry v2: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			SetCapabilitiesRegistryDescription,
		)
		if err != nil {
			return SetCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.SetCapabilitiesRegistry(opts, input.Registry, input.ChainSelectorDON)
			if err != nil {
				return nil, fmt.Errorf("failed to call SetCapabilitiesRegistry: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return SetCapabilitiesRegistryOpOutput{}, fmt.Errorf("failed to execute SetCapabilitiesRegistry: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for SetCapabilitiesRegistry on chain %d", input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully set DON registry %s on chain %d", input.Registry.String(), input.ChainSelector)
		}

		return SetCapabilitiesRegistryOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// Helper function to get registry instance from datastore

func getWorkflowRegistryV2FromDatastore(env *cldf.Environment, chainSelector uint64, qualifier string) (*workflow_registry_v2.WorkflowRegistry, error) {
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
