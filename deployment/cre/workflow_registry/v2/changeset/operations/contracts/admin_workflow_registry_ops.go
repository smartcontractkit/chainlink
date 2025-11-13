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

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	workflow_registry_v2 "github.com/smartcontractkit/chainlink-evm/gethwrappers/workflow/generated/workflow_registry_wrapper_v2"
)

const (
	PauseWorkflowDescription       = "adminPauseWorkflow on workflow registry v2"
	PauseAllByOwnerDescription     = "adminPauseAllByOwner on workflow registry v2"
	PauseAllByDONDescription       = "adminPauseAllByDON on workflow registry v2"
	PauseBatchWorkflowsDescription = "adminBatchPauseWorkflows on workflow registry v2"
)

// AdminPauseWorkflow Operation
type AdminPauseWorkflowOpInput struct {
	ChainSelector uint64   `json:"chainSelector"`
	Qualifier     string   `json:"qualifier"` // Qualifier to identify the specific workflow registry
	WorkflowID    [32]byte `json:"workflowID"`
}

type AdminPauseWorkflowOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var AdminPauseWorkflowOp = operations.NewOperation(
	"admin-pause-workflow-op",
	semver.MustParse("1.0.0"),
	"Admin Pause Workflow in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminPauseWorkflowOpInput) (AdminPauseWorkflowOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.AdminPauseWorkflow(opts, input.WorkflowID)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return AdminPauseWorkflowOpOutput{}, fmt.Errorf("failed to execute AdminPauseWorkflow: %w", err)
		}

		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminPauseWorkflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin paused workflow %x on chain %d", input.WorkflowID, input.ChainSelector)
		}

		return AdminPauseWorkflowOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// AdminBatchPauseWorkflows Operation
type AdminBatchPauseWorkflowsOpInput struct {
	ChainSelector uint64     `json:"chainSelector"`
	Qualifier     string     `json:"qualifier"` // Qualifier to identify the specific workflow registry
	WorkflowIDs   [][32]byte `json:"workflowIDs"`
}

type AdminBatchPauseWorkflowsOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var AdminBatchPauseWorkflowsOp = operations.NewOperation(
	"admin-batch-pause-workflows-op",
	semver.MustParse("1.0.0"),
	"Admin Batch Pause Workflows in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminBatchPauseWorkflowsOpInput) (AdminBatchPauseWorkflowsOpOutput, error) {
		if len(input.WorkflowIDs) == 0 {
			return AdminBatchPauseWorkflowsOpOutput{}, errors.New("must provide at least one workflow ID")
		}

		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.AdminBatchPauseWorkflows(opts, input.WorkflowIDs)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return AdminBatchPauseWorkflowsOpOutput{}, fmt.Errorf("failed to execute AdminBatchPauseWorkflows: %w", err)
		}

		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminBatchPauseWorkflows (%d workflows) on chain %d", len(input.WorkflowIDs), input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin batch paused %d workflows on chain %d", len(input.WorkflowIDs), input.ChainSelector)
		}

		return AdminBatchPauseWorkflowsOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// AdminPauseAllByOwner Operation
type AdminPauseAllByOwnerOpInput struct {
	ChainSelector uint64         `json:"chainSelector"`
	Qualifier     string         `json:"qualifier"` // Qualifier to identify the specific workflow registry
	Owner         common.Address `json:"owner"`
	Limit         *big.Int       `json:"limit"`
}

type AdminPauseAllByOwnerOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var AdminPauseAllByOwnerOp = operations.NewOperation(
	"admin-pause-all-by-owner-op",
	semver.MustParse("1.0.0"),
	"Admin Pause All By Owner in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminPauseAllByOwnerOpInput) (AdminPauseAllByOwnerOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.AdminPauseAllByOwner(opts, input.Owner, input.Limit)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return AdminPauseAllByOwnerOpOutput{}, fmt.Errorf("failed to execute AdminPauseAllByOwner: %w", err)
		}

		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminPauseAllByOwner (%s) on chain %d", input.Owner.String(), input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin paused all workflows for owner %s on chain %d", input.Owner.String(), input.ChainSelector)
		}

		return AdminPauseAllByOwnerOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)

// AdminPauseAllByDON Operation
type AdminPauseAllByDONOpInput struct {
	ChainSelector uint64   `json:"chainSelector"`
	Qualifier     string   `json:"qualifier"` // Qualifier to identify the specific workflow registry
	DONFamily     string   `json:"donFamily"`
	Limit         *big.Int `json:"limit"`
}

type AdminPauseAllByDONOpOutput struct {
	Success         bool                      `json:"success"`
	RegistryAddress common.Address            `json:"registryAddress"`
	MCMSOperation   *mcmstypes.BatchOperation `json:"mcmsOperation"`
}

var AdminPauseAllByDONOp = operations.NewOperation(
	"admin-pause-all-by-don-op",
	semver.MustParse("1.0.0"),
	"Admin Pause All By DON in WorkflowRegistry V2",
	func(b operations.Bundle, deps WorkflowRegistryOpDeps, input AdminPauseAllByDONOpInput) (AdminPauseAllByDONOpOutput, error) {
		// Execute the transaction using the strategy
		operation, _, err := deps.Strategy.Apply(func(opts *bind.TransactOpts) (*types.Transaction, error) {
			return deps.Registry.AdminPauseAllByDON(opts, input.DONFamily, input.Limit)
		})
		if err != nil {
			err = cldf.DecodeErr(workflow_registry_v2.WorkflowRegistryABI, err)
			return AdminPauseAllByDONOpOutput{}, fmt.Errorf("failed to execute AdminPauseAllByDON: %w", err)
		}

		if operation != nil {
			deps.Env.Logger.Infof("Created MCMS proposal for AdminPauseAllByDON (%s) on chain %d", input.DONFamily, input.ChainSelector)
		} else {
			deps.Env.Logger.Infof("Successfully admin paused all workflows for DON family %s on chain %d", input.DONFamily, input.ChainSelector)
		}

		return AdminPauseAllByDONOpOutput{
			Success:         true,
			MCMSOperation:   operation,
			RegistryAddress: deps.Registry.Address(),
		}, nil
	},
)
