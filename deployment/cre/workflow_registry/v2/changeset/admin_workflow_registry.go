package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset/state"
	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	"github.com/smartcontractkit/chainlink/deployment/cre/ocr3"
	"github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

// ChangeSet interface implementations for admin functions
var _ cldf.ChangeSetV2[AdminPauseWorkflowInput] = AdminPauseWorkflow{}
var _ cldf.ChangeSetV2[AdminBatchPauseWorkflowsInput] = AdminBatchPauseWorkflows{}
var _ cldf.ChangeSetV2[AdminPauseAllByOwnerInput] = AdminPauseAllByOwner{}
var _ cldf.ChangeSetV2[AdminPauseAllByDONInput] = AdminPauseAllByDON{}

// AdminPauseWorkflowInput pauses a specific workflow
type AdminPauseWorkflowInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	WorkflowID                [32]byte         `json:"workflowID"`                // Workflow ID to pause
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type AdminPauseWorkflow struct{}

// emptyQualifier is used when no specific workflow registry qualifier is needed
const emptyQualifier = ""

func (l AdminPauseWorkflow) VerifyPreconditions(e cldf.Environment, config AdminPauseWorkflowInput) error {
	return nil
}

func (l AdminPauseWorkflow) Apply(e cldf.Environment, config AdminPauseWorkflowInput) (cldf.ChangesetOutput, error) {
	// Get MCMS contracts if needed
	var mcmsContracts *commonchangeset.MCMSWithTimelockState
	if config.MCMSConfig != nil {
		var err error
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:           &e,
		MCMSContracts: mcmsContracts,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminPauseWorkflowOp, deps, contracts.AdminPauseWorkflowOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			WorkflowID:    config.WorkflowID,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: report.Output.Proposals,
	}, nil
}

// AdminBatchPauseWorkflowsInput pauses multiple workflows in a batch
type AdminBatchPauseWorkflowsInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	WorkflowIDs               [][32]byte       `json:"workflowIDs"`               // List of workflow IDs to pause
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
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
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:           &e,
		MCMSContracts: mcmsContracts,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminBatchPauseWorkflowsOp, deps, contracts.AdminBatchPauseWorkflowsOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			WorkflowIDs:   config.WorkflowIDs,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: report.Output.Proposals,
	}, nil
}

// AdminPauseAllByOwnerInput pauses all workflows for a specific owner
type AdminPauseAllByOwnerInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Owner                     common.Address   `json:"owner"`                     // Owner whose workflows should be paused
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
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
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:           &e,
		MCMSContracts: mcmsContracts,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminPauseAllByOwnerOp, deps, contracts.AdminPauseAllByOwnerOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			Owner:         config.Owner,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: report.Output.Proposals,
	}, nil
}

// AdminPauseAllByDONInput pauses all workflows for a specific DON family
type AdminPauseAllByDONInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	DONFamily                 string           `json:"donFamily"`                 // DON family whose workflows should be paused
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
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
		mcmsContracts, err = strategies.GetMCMSContracts(e, config.ChainSelector, emptyQualifier)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
	}

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:           &e,
		MCMSContracts: mcmsContracts,
	}
	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.AdminPauseAllByDONOp, deps, contracts.AdminPauseAllByDONOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			DONFamily:     config.DONFamily,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: report.Output.Proposals,
	}, nil
}
