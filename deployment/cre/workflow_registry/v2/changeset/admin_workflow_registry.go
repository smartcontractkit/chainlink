package changeset

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

// ChangeSet interface implementations for admin functions
var _ cldf.ChangeSetV2[AdminPauseWorkflowInput] = AdminPauseWorkflow{}
var _ cldf.ChangeSetV2[AdminBatchPauseWorkflowsInput] = AdminBatchPauseWorkflows{}
var _ cldf.ChangeSetV2[AdminPauseAllByOwnerInput] = AdminPauseAllByOwner{}
var _ cldf.ChangeSetV2[AdminPauseAllByDONInput] = AdminPauseAllByDON{}

// AdminPauseWorkflowInput pauses a specific workflow
type AdminPauseWorkflowInput struct {
	ChainSelector             uint64                   `json:"chainSelector"`
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	WorkflowID                [32]byte                 `json:"workflowID"`                // Workflow ID to pause
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type AdminPauseWorkflow struct{}

func (l AdminPauseWorkflow) VerifyPreconditions(e cldf.Environment, config AdminPauseWorkflowInput) error {
	return nil
}

func (l AdminPauseWorkflow) Apply(e cldf.Environment, config AdminPauseWorkflowInput) (cldf.ChangesetOutput, error) {
	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	chain, ok := e.BlockChains.EVMChains()[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", config.ChainSelector)
	}

	registry, err := contracts.GetWorkflowRegistryV2FromDatastore(&e, config.ChainSelector, config.WorkflowRegistryQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry address from datastore: %w", err)
	}

	// Create the appropriate strategy
	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		config.MCMSConfig,
		mcmsContracts,
		registry.Address(),
		contracts.PauseWorkflowDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:      &e,
		Strategy: strategy,
		Registry: registry,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminPauseWorkflowOp, deps, contracts.AdminPauseWorkflowOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			WorkflowID:    config.WorkflowID,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if report.Output.MCMSOperation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]types.BatchOperation{*report.Output.MCMSOperation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mcmsErr)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		}, nil
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

// AdminBatchPauseWorkflowsInput pauses multiple workflows in a batch
type AdminBatchPauseWorkflowsInput struct {
	ChainSelector             uint64                   `json:"chainSelector"`
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	WorkflowIDs               [][32]byte               `json:"workflowIDs"`               // List of workflow IDs to pause
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type AdminBatchPauseWorkflows struct{}

func (l AdminBatchPauseWorkflows) VerifyPreconditions(e cldf.Environment, config AdminBatchPauseWorkflowsInput) error {
	if len(config.WorkflowIDs) == 0 {
		return errors.New("must provide at least one workflow ID")
	}
	return nil
}

func (l AdminBatchPauseWorkflows) Apply(e cldf.Environment, config AdminBatchPauseWorkflowsInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	chain, ok := e.BlockChains.EVMChains()[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", config.ChainSelector)
	}

	registry, err := contracts.GetWorkflowRegistryV2FromDatastore(&e, config.ChainSelector, config.WorkflowRegistryQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry address from datastore: %w", err)
	}

	// Create the appropriate strategy
	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		config.MCMSConfig,
		mcmsContracts,
		registry.Address(),
		contracts.PauseBatchWorkflowsDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:      &e,
		Strategy: strategy,
		Registry: registry,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminBatchPauseWorkflowsOp, deps, contracts.AdminBatchPauseWorkflowsOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			WorkflowIDs:   config.WorkflowIDs,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if report.Output.MCMSOperation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]types.BatchOperation{*report.Output.MCMSOperation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mcmsErr)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		}, nil
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

// AdminPauseAllByOwnerInput pauses all workflows for a specific owner
type AdminPauseAllByOwnerInput struct {
	ChainSelector             uint64                   `yaml:"chainSelector"`
	WorkflowRegistryQualifier string                   `yaml:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Owner                     common.Address           `yaml:"owner"`                     // Owner whose workflows should be paused
	MCMSConfig                *crecontracts.MCMSConfig `yaml:"mcmsConfig,omitempty"`      // MCMS configuration
	Limit                     *big.Int                 `yaml:"limit"`
}

type AdminPauseAllByOwner struct{}

func (l AdminPauseAllByOwner) VerifyPreconditions(e cldf.Environment, config AdminPauseAllByOwnerInput) error {
	return nil
}

func (l AdminPauseAllByOwner) Apply(e cldf.Environment, config AdminPauseAllByOwnerInput) (cldf.ChangesetOutput, error) {
	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	chain, ok := e.BlockChains.EVMChains()[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", config.ChainSelector)
	}

	registry, err := contracts.GetWorkflowRegistryV2FromDatastore(&e, config.ChainSelector, config.WorkflowRegistryQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry address from datastore: %w", err)
	}

	// Create the appropriate strategy
	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		config.MCMSConfig,
		mcmsContracts,
		registry.Address(),
		contracts.PauseAllByOwnerDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:      &e,
		Strategy: strategy,
		Registry: registry,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminPauseAllByOwnerOp, deps, contracts.AdminPauseAllByOwnerOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			Owner:         config.Owner,
			Limit:         config.Limit,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if report.Output.MCMSOperation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]types.BatchOperation{*report.Output.MCMSOperation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mcmsErr)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		}, nil
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

// AdminPauseAllByDONInput pauses all workflows for a specific DON family
type AdminPauseAllByDONInput struct {
	ChainSelector             uint64                   `yaml:"chainSelector"`
	WorkflowRegistryQualifier string                   `yaml:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	DONFamily                 string                   `yaml:"donFamily"`                 // DON family whose workflows should be paused
	Limit                     *big.Int                 `yaml:"limit"`
	MCMSConfig                *crecontracts.MCMSConfig `yaml:"mcmsConfig,omitempty"` // MCMS configuration
}

type AdminPauseAllByDON struct{}

func (l AdminPauseAllByDON) VerifyPreconditions(e cldf.Environment, config AdminPauseAllByDONInput) error {
	return nil
}

func (l AdminPauseAllByDON) Apply(e cldf.Environment, config AdminPauseAllByDONInput) (cldf.ChangesetOutput, error) {
	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, *config.MCMSConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	chain, ok := e.BlockChains.EVMChains()[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("chain with selector %d not found", config.ChainSelector)
	}

	registry, err := contracts.GetWorkflowRegistryV2FromDatastore(&e, config.ChainSelector, config.WorkflowRegistryQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get workflow registry address from datastore: %w", err)
	}

	// Create the appropriate strategy
	strategy, err := strategies.CreateStrategy(
		chain,
		e,
		config.MCMSConfig,
		mcmsContracts,
		registry.Address(),
		contracts.PauseAllByDONDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:      &e,
		Strategy: strategy,
		Registry: registry,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminPauseAllByDONOp, deps, contracts.AdminPauseAllByDONOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			DONFamily:     config.DONFamily,
			Limit:         config.Limit,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if report.Output.MCMSOperation != nil {
		proposal, mcmsErr := strategy.BuildProposal([]types.BatchOperation{*report.Output.MCMSOperation})
		if mcmsErr != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", mcmsErr)
		}

		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
			Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		}, nil
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}
