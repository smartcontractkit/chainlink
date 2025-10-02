package contracts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	mcmslib "github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
)

const (
	PauseWorkflowDescription       = "adminPauseWorkflow on workflow registry v2"
	PauseAllByOwnerDescription     = "adminPauseAllByOwner on workflow registry v2"
	PauseAllByDONDescription       = "adminPauseAllByDON on workflow registry v2"
	PauseBatchWorkflowsDescription = "adminBatchPauseWorkflows on workflow registry v2"
)

// AdminPauseWorkflow Operation
type AdminPauseWorkflowOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	WorkflowID    [32]byte         `json:"workflowID"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type AdminPauseWorkflowOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var AdminPauseWorkflowOp = operations.NewOperation(
	"admin-pause-workflow-op",
	semver.MustParse("1.0.0"),
	"Admin Pause Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminPauseWorkflowOpInput) (AdminPauseWorkflowOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return AdminPauseWorkflowOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return AdminPauseWorkflowOpOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			PauseWorkflowDescription,
		)
		if err != nil {
			return AdminPauseWorkflowOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.AdminPauseWorkflow(opts, input.WorkflowID)
			if err != nil {
				return nil, fmt.Errorf("failed to call AdminPauseWorkflow: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return AdminPauseWorkflowOpOutput{}, fmt.Errorf("failed to execute AdminPauseWorkflow: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminPauseWorkflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin paused workflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		}

		return AdminPauseWorkflowOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// AdminBatchPauseWorkflows Operation
type AdminBatchPauseWorkflowsOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	WorkflowIDs   [][32]byte       `json:"workflowIDs"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type AdminBatchPauseWorkflowsOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var AdminBatchPauseWorkflowsOp = operations.NewOperation(
	"admin-batch-pause-workflows-op",
	semver.MustParse("1.0.0"),
	"Admin Batch Pause Workflows in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminBatchPauseWorkflowsOpInput) (AdminBatchPauseWorkflowsOpOutput, error) {
		if len(input.WorkflowIDs) == 0 {
			return AdminBatchPauseWorkflowsOpOutput{}, errors.New("must provide at least one workflow ID")
		}

		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return AdminBatchPauseWorkflowsOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return AdminBatchPauseWorkflowsOpOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			PauseBatchWorkflowsDescription,
		)
		if err != nil {
			return AdminBatchPauseWorkflowsOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.AdminBatchPauseWorkflows(opts, input.WorkflowIDs)
			if err != nil {
				return nil, fmt.Errorf("failed to call AdminBatchPauseWorkflows: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return AdminBatchPauseWorkflowsOpOutput{}, fmt.Errorf("failed to execute AdminBatchPauseWorkflows: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminBatchPauseWorkflows (%d workflows) on chain %d", len(input.WorkflowIDs), input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin batch paused %d workflows on chain %d", len(input.WorkflowIDs), input.ChainSelector)
		}

		return AdminBatchPauseWorkflowsOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// AdminPauseAllByOwner Operation
type AdminPauseAllByOwnerOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	Owner         common.Address   `json:"owner"`
	Limit         *big.Int         `json:"limit"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type AdminPauseAllByOwnerOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var AdminPauseAllByOwnerOp = operations.NewOperation(
	"admin-pause-all-by-owner-op",
	semver.MustParse("1.0.0"),
	"Admin Pause All By Owner in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminPauseAllByOwnerOpInput) (AdminPauseAllByOwnerOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return AdminPauseAllByOwnerOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return AdminPauseAllByOwnerOpOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			PauseAllByOwnerDescription,
		)
		if err != nil {
			return AdminPauseAllByOwnerOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.AdminPauseAllByOwner(opts, input.Owner, input.Limit)
			if err != nil {
				return nil, fmt.Errorf("failed to call AdminPauseAllByOwner: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return AdminPauseAllByOwnerOpOutput{}, fmt.Errorf("failed to execute AdminPauseAllByOwner: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminPauseAllByOwner (%s) on chain %d", input.Owner.String(), input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin paused all workflows for owner %s on chain %d", input.Owner.String(), input.ChainSelector)
		}

		return AdminPauseAllByOwnerOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)

// AdminPauseAllByDON Operation
type AdminPauseAllByDONOpInput struct {
	ChainSelector uint64           `json:"chainSelector"`
	Qualifier     string           `json:"qualifier"` // Qualifier to identify the specific workflow registry
	DONFamily     string           `json:"donFamily"`
	Limit         *big.Int         `json:"limit"`
	MCMSConfig    *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`
}

type AdminPauseAllByDONOpOutput struct {
	Success   bool                       `json:"success"`
	Proposals []mcmslib.TimelockProposal `json:"proposals,omitempty"`
}

var AdminPauseAllByDONOp = operations.NewOperation(
	"admin-pause-all-by-don-op",
	semver.MustParse("1.0.0"),
	"Admin Pause All By DON in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminPauseAllByDONOpInput) (AdminPauseAllByDONOpOutput, error) {
		chain, ok := deps.Env.BlockChains.EVMChains()[input.ChainSelector]
		if !ok {
			return AdminPauseAllByDONOpOutput{}, fmt.Errorf("chain with selector %d not found", input.ChainSelector)
		}

		registry, err := getWorkflowRegistryV2FromDatastore(deps.Env, input.ChainSelector, input.Qualifier)
		if err != nil {
			return AdminPauseAllByDONOpOutput{}, fmt.Errorf("failed to get workflow registry: %w", err)
		}

		// Create the appropriate strategy
		strategy, err := strategies.CreateStrategy(
			chain,
			*deps.Env,
			input.MCMSConfig,
			deps.MCMSContracts,
			registry.Address(),
			PauseAllByDONDescription,
		)
		if err != nil {
			return AdminPauseAllByDONOpOutput{}, fmt.Errorf("failed to create strategy: %w", err)
		}

		// Execute the transaction using the strategy
		proposals, err := strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			tx, err := registry.AdminPauseAllByDON(opts, input.DONFamily, input.Limit)
			if err != nil {
				return nil, fmt.Errorf("failed to call AdminPauseAllByDON: %w", err)
			}
			return tx, nil
		})
		if err != nil {
			return AdminPauseAllByDONOpOutput{}, fmt.Errorf("failed to execute AdminPauseAllByDON: %w", err)
		}

		if input.MCMSConfig != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminPauseAllByDON (%s) on chain %d", input.DONFamily, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin paused all workflows for DON family %s on chain %d", input.DONFamily, input.ChainSelector)
		}

		return AdminPauseAllByDONOpOutput{
			Success:   true,
			Proposals: proposals,
		}, nil
	},
)
