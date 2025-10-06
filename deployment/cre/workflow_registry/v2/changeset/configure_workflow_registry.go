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

// ChangeSet interface implementations for each configuration function
var _ cldf.ChangeSetV2[SetConfigInput] = SetConfig{}
var _ cldf.ChangeSetV2[UpdateAllowedSignersInput] = UpdateAllowedSigners{}
var _ cldf.ChangeSetV2[SetWorkflowOwnerConfigInput] = SetWorkflowOwnerConfig{}
var _ cldf.ChangeSetV2[SetDONLimitInput] = SetDONLimit{}
var _ cldf.ChangeSetV2[SetUserDONOverrideInput] = SetUserDONOverride{}
var _ cldf.ChangeSetV2[SetCapabilitiesRegistryInput] = SetCapabilitiesRegistry{}

// SetConfigInput configures metadata validation settings for workflow registry v2
type SetConfigInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	NameLen                   uint8            `json:"nameLen"`                   // Maximum length for workflow names
	TagLen                    uint8            `json:"tagLen"`                    // Maximum length for workflow tags
	URLLen                    uint8            `json:"urlLen"`                    // Maximum length for URLs
	AttrLen                   uint16           `json:"attrLen"`                   // Maximum length for attributes
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetConfig struct{}

func (l SetConfig) VerifyPreconditions(e cldf.Environment, config SetConfigInput) error {
	return nil
}

func (l SetConfig) Apply(e cldf.Environment, config SetConfigInput) (cldf.ChangesetOutput, error) {
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
		contracts.SetConfigOp, deps, contracts.SetConfigOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			NameLen:       config.NameLen,
			TagLen:        config.TagLen,
			URLLen:        config.URLLen,
			AttrLen:       config.AttrLen,
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

// UpdateAllowedSignersInput updates the list of allowed signers for workflow registry v2
type UpdateAllowedSignersInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Signers                   []common.Address `json:"signers"`                   // List of signer addresses
	Allowed                   bool             `json:"allowed"`                   // Whether to allow or disallow these signers
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	Description               string           `json:"description,omitempty"`     // Description for MCMS proposal
}

type UpdateAllowedSigners struct{}

func (l UpdateAllowedSigners) VerifyPreconditions(e cldf.Environment, config UpdateAllowedSignersInput) error {
	if len(config.Signers) == 0 {
		return errors.New("must provide at least one signer")
	}
	return nil
}

func (l UpdateAllowedSigners) Apply(e cldf.Environment, config UpdateAllowedSignersInput) (cldf.ChangesetOutput, error) {
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
		contracts.UpdateAllowedSignersOp, deps, contracts.UpdateAllowedSignersOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			Signers:       config.Signers,
			Allowed:       config.Allowed,

			MCMSConfig: config.MCMSConfig,
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

// SetWorkflowOwnerConfigInput configures workflow owner-specific settings
type SetWorkflowOwnerConfigInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Owner                     common.Address   `json:"owner"`                     // Workflow owner address
	Config                    []byte           `json:"config"`                    // Owner-specific configuration data
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetWorkflowOwnerConfig struct{}

func (l SetWorkflowOwnerConfig) VerifyPreconditions(e cldf.Environment, config SetWorkflowOwnerConfigInput) error {
	return nil
}

func (l SetWorkflowOwnerConfig) Apply(e cldf.Environment, config SetWorkflowOwnerConfigInput) (cldf.ChangesetOutput, error) {
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
		contracts.SetWorkflowOwnerConfigOp, deps, contracts.SetWorkflowOwnerConfigOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			Owner:         config.Owner,
			Config:        config.Config,

			MCMSConfig: config.MCMSConfig,
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

// SetDONLimitInput configures DON workflow limits
type SetDONLimitInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	DONFamily                 string           `json:"donFamily"`                 // DON family identifier
	DONLimit                  uint32           `json:"donlimit"`                  // Maximum number of workflows per owner
	UserDefaultLimit          uint32           `json:"userDefaultLimit"`          // Whether the limit is enabled
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetDONLimit struct{}

func (l SetDONLimit) VerifyPreconditions(e cldf.Environment, config SetDONLimitInput) error {
	return nil
}

func (l SetDONLimit) Apply(e cldf.Environment, config SetDONLimitInput) (cldf.ChangesetOutput, error) {
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
		contracts.SetDONLimitOp, deps, contracts.SetDONLimitOpInput{
			ChainSelector:    config.ChainSelector,
			Qualifier:        config.WorkflowRegistryQualifier,
			DONFamily:        config.DONFamily,
			DONLimit:         config.DONLimit,
			UserDefaultLimit: config.UserDefaultLimit,
			MCMSConfig:       config.MCMSConfig,
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

// SetUserDONOverrideInput configures user-specific DON overrides
type SetUserDONOverrideInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	User                      common.Address   `json:"user"`                      // User address
	DONFamily                 string           `json:"donFamily"`                 // DON family identifier
	Limit                     uint32           `json:"limit"`                     // User-specific limit
	Enabled                   bool             `json:"enabled"`                   // Whether the override is enabled
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetUserDONOverride struct{}

func (l SetUserDONOverride) VerifyPreconditions(e cldf.Environment, config SetUserDONOverrideInput) error {
	return nil
}

func (l SetUserDONOverride) Apply(e cldf.Environment, config SetUserDONOverrideInput) (cldf.ChangesetOutput, error) {
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
		contracts.SetUserDONOverrideOp, deps, contracts.SetUserDONOverrideOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			User:          config.User,
			DONFamily:     config.DONFamily,
			Limit:         config.Limit,
			Enabled:       config.Enabled,
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

// SetCapabilitiesRegistryInput configures the Capabilities registry address
type SetCapabilitiesRegistryInput struct {
	ChainSelector             uint64           `json:"chainSelector"`
	WorkflowRegistryQualifier string           `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Registry                  common.Address   `json:"registry"`                  // DON registry contract address
	ChainSelectorDON          uint64           `json:"chainSelectorDON"`          // Chain selector where the DON registry exists
	MCMSConfig                *ocr3.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetCapabilitiesRegistry struct{}

func (l SetCapabilitiesRegistry) VerifyPreconditions(e cldf.Environment, config SetCapabilitiesRegistryInput) error {
	return nil
}

func (l SetCapabilitiesRegistry) Apply(e cldf.Environment, config SetCapabilitiesRegistryInput) (cldf.ChangesetOutput, error) {
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
		contracts.SetCapabilitiesRegistryOp, deps, contracts.SetCapabilitiesRegistryOpInput{
			ChainSelector:    config.ChainSelector,
			Qualifier:        config.WorkflowRegistryQualifier,
			Registry:         config.Registry,
			ChainSelectorDON: config.ChainSelectorDON,
			MCMSConfig:       config.MCMSConfig,
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
