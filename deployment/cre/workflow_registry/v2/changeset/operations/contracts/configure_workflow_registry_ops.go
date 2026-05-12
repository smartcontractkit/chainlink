package contracts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	cldf_evm "github.com/smartcontractkit/chainlink-deployments-framework/chain/evm"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
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
	BatchSetUserDONOverrideDescription = "batchSetUserDonOverride on workflow registry v2"
	SetCapabilitiesRegistryDescription = "setCapabilitiesRegistry on workflow registry v2"
)

// Common dependencies for workflow registry operations
type WorkflowRegistryOpDeps struct {
	Env      *cldf.Environment
	Strategy strategies.TransactionStrategy
	Registry *workflow_registry_v2.WorkflowRegistry
	// Chain is the EVM chain whose datastore produced deps.Registry. Operations that bypass the
	// strategy abstraction (e.g. BatchSetUserDONOverrideOp) require it when MCMSConfig is nil,
	// to sign and confirm transactions on-chain via DeployerKey + Confirm. MCMS-only invocations
	// of those operations may leave it nil; the chain selector then comes from the input.
	// When provided, it must agree with the input's chain selector (the op validates this to
	// avoid silently emitting a proposal targeting the wrong chain). Other operations that
	// route everything through the strategy may leave it nil.
	Chain *cldf_evm.Chain
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

// BatchSetUserDONOverride Operation
//
// Iterates over Overrides and either confirms each setUserDONOverride transaction directly (no MCMS)
// or assembles a single MCMS BatchOperation containing one transaction per override. Using a single
// BatchOperation keeps the resulting proposal compact (1 batch, N transactions) instead of producing
// N separate batch operations as repeated SetUserDONOverride calls would.
type SetUserDONOverrideEntry struct {
	User      common.Address `json:"user"`
	DONFamily string         `json:"donFamily"`
	Limit     uint32         `json:"limit"`
	Enabled   bool           `json:"enabled"`
}

type BatchSetUserDONOverrideOpInput struct {
	// ChainSelector and Qualifier are kept on the input to make the operation invocation uniquely
	// identifiable (the registry itself is passed via deps).
	ChainSelector uint64 `json:"chainSelector"`
	Qualifier     string `json:"qualifier"`

	Overrides  []SetUserDONOverrideEntry `json:"overrides"`
	MCMSConfig *contracts.MCMSConfig     `json:"mcmsConfig,omitempty"`
}

type BatchSetUserDONOverrideOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var BatchSetUserDONOverrideOp = operations.NewOperation(
	"batch-set-user-don-override-op",
	semver.MustParse("1.0.0"),
	"Batch Set User DON Override in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input BatchSetUserDONOverrideOpInput) (BatchSetUserDONOverrideOpOutput, error) {
		if len(input.Overrides) == 0 {
			return BatchSetUserDONOverrideOpOutput{}, errors.New("must provide at least one override")
		}
		// deps.Chain is only strictly required for the non-MCMS path, where we need the deployer
		// key to sign transactions and Confirm() to wait for them on-chain. In the MCMS path the
		// op only builds calldata (via SimTransactOpts) and assembles an MCMS BatchOperation, so
		// an MCMS-only caller can supply deps.Registry + strategy without a chain pointer.
		if input.MCMSConfig == nil && deps.Chain == nil {
			return BatchSetUserDONOverrideOpOutput{}, errors.New("deps.Chain is required when MCMSConfig is nil (needed to sign and confirm transactions on-chain)")
		}

		// Resolve the authoritative chain selector. Prefer deps.Chain (the chain whose lookup
		// produced deps.Registry) when present; that keeps the on-chain lookup and the proposal
		// target trivially in sync. Fall back to input.ChainSelector for MCMS-only callers that
		// deliberately don't pass a chain pointer. When both are present they must agree, else
		// we'd silently emit a proposal targeting the wrong chain.
		chainSelector := input.ChainSelector
		if deps.Chain != nil {
			chainSelector = deps.Chain.ChainSelector()
			if input.ChainSelector != chainSelector {
				return BatchSetUserDONOverrideOpOutput{}, fmt.Errorf(
					"input.ChainSelector (%d) does not match deps.Chain.Selector (%d); refusing to build proposal with ambiguous target chain",
					input.ChainSelector, chainSelector,
				)
			}
		}

		// MCMS path uses simulated tx opts to produce calldata without sending; non-MCMS uses the deployer key.
		var txOpts *bind.TransactOpts
		if input.MCMSConfig != nil {
			txOpts = cldf.SimTransactOpts()
		} else {
			txOpts = deps.Chain.DeployerKey
		}

		var mcmsTxs []mcmstypes.Transaction
		for _, entry := range input.Overrides {
			tx, err := deps.Registry.SetUserDONOverride(txOpts, entry.User, entry.DONFamily, entry.Limit, entry.Enabled)
			if err != nil {
				err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
				return BatchSetUserDONOverrideOpOutput{}, fmt.Errorf("failed to build SetUserDONOverride for %s: %w", entry.User.Hex(), err)
			}

			if input.MCMSConfig == nil {
				if _, cErr := deps.Chain.Confirm(tx); cErr != nil {
					return BatchSetUserDONOverrideOpOutput{}, fmt.Errorf("failed to confirm SetUserDONOverride for %s (tx %s): %w", entry.User.Hex(), tx.Hash().String(), cErr)
				}
				continue
			}

			mtx, err := cldfproposalutils.TransactionForChain(chainSelector, deps.Registry.Address().Hex(), tx.Data(), big.NewInt(0), "", nil)
			if err != nil {
				return BatchSetUserDONOverrideOpOutput{}, fmt.Errorf("failed to build MCMS transaction for %s: %w", entry.User.Hex(), err)
			}
			mcmsTxs = append(mcmsTxs, mtx)
		}

		var mergedOp *mcmstypes.BatchOperation
		if input.MCMSConfig != nil {
			mergedOp = &mcmstypes.BatchOperation{
				ChainSelector: mcmstypes.ChainSelector(chainSelector),
				Transactions:  mcmsTxs,
			}
			deps.Env.Logger.Infof("Created MCMS batch with %d SetUserDONOverride transactions on chain %d", len(mcmsTxs), chainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully applied %d SetUserDONOverride transactions on chain %d", len(input.Overrides), chainSelector)
		}

		return BatchSetUserDONOverrideOpOutput{
			Success:         true,
			MCMSOperation:   mergedOp,
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
