package changeset

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	evmstate "github.com/smartcontractkit/cld-changesets/legacy/pkg/family/evm"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink/deployment/cre/common/strategies"
	crecontracts "github.com/smartcontractkit/chainlink/deployment/cre/contracts"
	"github.com/smartcontractkit/chainlink/deployment/cre/workflow_registry/v2/changeset/operations/contracts"
)

// ChangeSet interface implementations for each configuration function
var _ cldf.ChangeSetV2[SetConfigInput] = SetConfig{}
var _ cldf.ChangeSetV2[UpdateAllowedSignersInput] = UpdateAllowedSigners{}
var _ cldf.ChangeSetV2[SetWorkflowOwnerConfigInput] = SetWorkflowOwnerConfig{}
var _ cldf.ChangeSetV2[SetDONLimitInput] = SetDONLimit{}
var _ cldf.ChangeSetV2[SetUserDONOverrideInput] = SetUserDONOverride{}
var _ cldf.ChangeSetV2[BatchSetUserDONOverrideInput] = BatchSetUserDONOverride{}
var _ cldf.ChangeSetV2[SetCapabilitiesRegistryInput] = SetCapabilitiesRegistry{}

// validateRegistryTarget checks the two fields every registry changeset uses to find its registry.
// Before, a missing chain selector or qualifier was only noticed (or not) by the datastore lookup.
func validateRegistryTarget(chainSelector uint64, qualifier string) error {
	if chainSelector == 0 {
		return errors.New("chainSelector must be set")
	}
	if qualifier == "" {
		return errors.New("workflowRegistryQualifier must be set: it picks which WorkflowRegistry to change")
	}
	return nil
}

// prepareWorkflowRegistryDeps loads MCMS state, EVM chain, registry contract, and the transaction
// strategy needed to execute a workflow-registry operation. The returned deps is ready to pass to
// operations.ExecuteOperation and the strategy is reachable via deps.Strategy for the subsequent
// BuildProposal call in finalizeWorkflowRegistryOutput.
//
// Extracted to remove the ~30-line setup boilerplate that was previously duplicated across each
// Apply function in this file. New code should prefer this helper; existing Apply functions can
// be migrated incrementally.
func prepareWorkflowRegistryDeps(
	e cldf.Environment,
	chainSelector uint64,
	qualifier string,
	mcmsConfig *crecontracts.MCMSConfig,
	description string,
) (contracts.WorkflowRegistryOpDeps, error) {
	var mcmsContracts *evmstate.MCMSWithTimelockState
	if mcmsConfig != nil {
		loaded, err := strategies.GetMCMSContracts(e, chainSelector, *mcmsConfig)
		if err != nil {
			return contracts.WorkflowRegistryOpDeps{}, fmt.Errorf("failed to get MCMS contracts: %w", err)
		}
		mcmsContracts = loaded
	}

	chain, ok := e.BlockChains.EVMChains()[chainSelector]
	if !ok {
		return contracts.WorkflowRegistryOpDeps{}, fmt.Errorf("chain with selector %d not found", chainSelector)
	}

	registry, err := contracts.GetWorkflowRegistryV2FromDatastore(&e, chainSelector, qualifier)
	if err != nil {
		return contracts.WorkflowRegistryOpDeps{}, fmt.Errorf("failed to get workflow registry address from datastore: %w", err)
	}

	strategy, err := strategies.CreateStrategy(chain, e, mcmsConfig, mcmsContracts, registry.Address(), description)
	if err != nil {
		return contracts.WorkflowRegistryOpDeps{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	return contracts.WorkflowRegistryOpDeps{
		Env:      &e,
		Strategy: strategy,
		Registry: registry,
		Chain:    &chain,
	}, nil
}

// finalizeWorkflowRegistryOutput wraps a completed operation report into a ChangesetOutput,
// building an MCMS timelock proposal from mcmsOp when non-nil.
func finalizeWorkflowRegistryOutput(
	strategy strategies.TransactionStrategy,
	mcmsOp *types.BatchOperation,
	report operations.Report[any, any],
) (cldf.ChangesetOutput, error) {
	if mcmsOp == nil {
		return cldf.ChangesetOutput{
			Reports: []operations.Report[any, any]{report},
		}, nil
	}
	proposal, err := strategy.BuildProposal([]types.BatchOperation{*mcmsOp})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal: %w", err)
	}
	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{*proposal},
		Reports:               []operations.Report[any, any]{report},
	}, nil
}

// SetConfigInput configures metadata validation settings for workflow registry v2
type SetConfigInput struct {
	ChainSelector             uint64 `json:"chainSelector"`
	WorkflowRegistryQualifier string `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	// The five lengths are pointers and must all be set. On-chain 0 means "unlimited", so a length left out
	// of the input used to silently remove that cap. Set 0 explicitly if you really want unlimited.
	NameLen    *uint8                   `json:"nameLen"`              // Maximum length for workflow names (0 = unlimited)
	TagLen     *uint8                   `json:"tagLen"`               // Maximum length for workflow tags (0 = unlimited)
	URLLen     *uint8                   `json:"urlLen"`               // Maximum length for URLs (0 = unlimited)
	AttrLen    *uint16                  `json:"attrLen"`              // Maximum length for attributes (0 = unlimited)
	ExpiryLen  *uint32                  `json:"expiryLen"`            // Maximum expiry duration for allowlisted secret requests (0 = unlimited)
	MCMSConfig *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"` // MCMS configuration
}

type SetConfig struct{}

func (l SetConfig) VerifyPreconditions(e cldf.Environment, config SetConfigInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	const unset = "%s must be set explicitly: omitting it would have set it to 0, which means unlimited on-chain"
	switch {
	case config.NameLen == nil:
		return fmt.Errorf(unset, "nameLen")
	case config.TagLen == nil:
		return fmt.Errorf(unset, "tagLen")
	case config.URLLen == nil:
		return fmt.Errorf(unset, "urlLen")
	case config.AttrLen == nil:
		return fmt.Errorf(unset, "attrLen")
	case config.ExpiryLen == nil:
		return fmt.Errorf(unset, "expiryLen")
	}
	return nil
}

func (l SetConfig) Apply(e cldf.Environment, config SetConfigInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Get MCMS contracts if needed
	var mcmsContracts *evmstate.MCMSWithTimelockState
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
		contracts.SetConfigDescription,
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
		contracts.SetConfigOp, deps, contracts.SetConfigOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			NameLen:       *config.NameLen,
			TagLen:        *config.TagLen,
			URLLen:        *config.URLLen,
			AttrLen:       *config.AttrLen,
			ExpiryLen:     *config.ExpiryLen,
			MCMSConfig:    config.MCMSConfig,
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

// UpdateAllowedSignersInput updates the list of allowed signers for workflow registry v2
type UpdateAllowedSignersInput struct {
	ChainSelector             uint64                   `json:"chainSelector"`
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Signers                   []common.Address         `json:"signers"`                   // List of signer addresses
	Allowed                   *bool                    `json:"allowed"`                   // Whether to allow or disallow these signers; must be set (a missing value used to mean false, which REMOVES the signers)
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	Description               string                   `json:"description,omitempty"`     // Description for MCMS proposal
}

type UpdateAllowedSigners struct{}

func (l UpdateAllowedSigners) VerifyPreconditions(e cldf.Environment, config UpdateAllowedSignersInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	if len(config.Signers) == 0 {
		return errors.New("must provide at least one signer")
	}
	for i, signer := range config.Signers {
		if signer == (common.Address{}) {
			return fmt.Errorf("signers[%d] is the zero address, which the contract rejects", i)
		}
	}
	if config.Allowed == nil {
		return errors.New("allowed must be set explicitly: omitting it would have meant false, which removes the listed signers")
	}
	return nil
}

func (l UpdateAllowedSigners) Apply(e cldf.Environment, config UpdateAllowedSignersInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Get MCMS contracts if needed
	var mcmsContracts *evmstate.MCMSWithTimelockState
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
		contracts.UpdateAllowedSignersDescription,
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
		contracts.UpdateAllowedSignersOp, deps, contracts.UpdateAllowedSignersOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			Signers:       config.Signers,
			Allowed:       *config.Allowed,

			MCMSConfig: config.MCMSConfig,
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

// SetWorkflowOwnerConfigInput configures workflow owner-specific settings
type SetWorkflowOwnerConfigInput struct {
	ChainSelector             uint64                   `json:"chainSelector"`
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Owner                     common.Address           `json:"owner"`                     // Workflow owner address
	Config                    []byte                   `json:"config"`                    // Owner-specific configuration data
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetWorkflowOwnerConfig struct{}

func (l SetWorkflowOwnerConfig) VerifyPreconditions(e cldf.Environment, config SetWorkflowOwnerConfigInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	if config.Owner == (common.Address{}) {
		return errors.New("owner must be set: the zero address would write a config for nobody")
	}
	if len(config.Config) == 0 {
		return errors.New("config must not be empty: omitting it would have overwritten the owner's stored config with empty bytes")
	}
	return nil
}

func (l SetWorkflowOwnerConfig) Apply(e cldf.Environment, config SetWorkflowOwnerConfigInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Get MCMS contracts if needed
	var mcmsContracts *evmstate.MCMSWithTimelockState
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
		contracts.SetWorkflowOwnerConfigDescription,
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

// SetDONLimitInput configures DON workflow limits
type SetDONLimitInput struct {
	ChainSelector             uint64 `json:"chainSelector"`
	WorkflowRegistryQualifier string `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	DONFamily                 string `json:"donFamily"`                 // DON family identifier
	// Both limits are pointers and must be set. On-chain a DON limit of 0 DISABLES the DON family and a user
	// default of 0 means users get no active workflows by default, so a value left out of the input used to
	// silently do that. Set 0 explicitly if you really want it.
	DONLimit         *uint32                  `json:"donlimit"`             // Maximum number of active workflows on the DON family (0 = DON disabled)
	UserDefaultLimit *uint32                  `json:"userDefaultLimit"`     // Default per-user limit of active workflows (0 = none unless overridden)
	MCMSConfig       *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"` // MCMS configuration
}

type SetDONLimit struct{}

func (l SetDONLimit) VerifyPreconditions(e cldf.Environment, config SetDONLimitInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	if config.DONFamily == "" {
		return errors.New("donFamily must be set")
	}
	if config.DONLimit == nil {
		return errors.New("donlimit must be set explicitly: omitting it would have set it to 0, which disables the DON family")
	}
	if config.UserDefaultLimit == nil {
		return errors.New("userDefaultLimit must be set explicitly: omitting it would have set it to 0, so users get no active workflows on this DON family by default")
	}
	if *config.UserDefaultLimit > *config.DONLimit {
		return fmt.Errorf("userDefaultLimit (%d) is above donlimit (%d); the contract reverts on this", *config.UserDefaultLimit, *config.DONLimit)
	}
	return nil
}

func (l SetDONLimit) Apply(e cldf.Environment, config SetDONLimitInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Get MCMS contracts if needed
	var mcmsContracts *evmstate.MCMSWithTimelockState
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
		contracts.SetDONLimitDescription,
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
		contracts.SetDONLimitOp, deps, contracts.SetDONLimitOpInput{
			ChainSelector:    config.ChainSelector,
			Qualifier:        config.WorkflowRegistryQualifier,
			DONFamily:        config.DONFamily,
			DONLimit:         *config.DONLimit,
			UserDefaultLimit: *config.UserDefaultLimit,
			MCMSConfig:       config.MCMSConfig,
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

// SetUserDONOverrideInput configures user-specific DON overrides
type SetUserDONOverrideInput struct {
	ChainSelector             uint64                   `json:"chainSelector"`
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	User                      common.Address           `json:"user"`                      // User address
	DONFamily                 string                   `json:"donFamily"`                 // DON family identifier
	Limit                     *uint32                  `json:"limit"`                     // User-specific limit; must be set when enabled is true
	Enabled                   *bool                    `json:"enabled"`                   // Whether the override is enabled; must be set (a missing value used to mean false, which DELETES the override)
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetUserDONOverride struct{}

func (l SetUserDONOverride) VerifyPreconditions(e cldf.Environment, config SetUserDONOverrideInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	if config.User == (common.Address{}) {
		return errors.New("user must be set: the zero address is not a real user")
	}
	if config.DONFamily == "" {
		return errors.New("donFamily must be set")
	}
	_, _, err := contracts.SetUserDONOverrideEntry{Limit: config.Limit, Enabled: config.Enabled}.Resolve()
	return err
}

func (l SetUserDONOverride) Apply(e cldf.Environment, config SetUserDONOverrideInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}
	limit, enabled, err := contracts.SetUserDONOverrideEntry{Limit: config.Limit, Enabled: config.Enabled}.Resolve()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deps, err := prepareWorkflowRegistryDeps(e, config.ChainSelector, config.WorkflowRegistryQualifier, config.MCMSConfig, contracts.SetUserDONOverrideDescription)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.SetUserDONOverrideOp, deps, contracts.SetUserDONOverrideOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			User:          config.User,
			DONFamily:     config.DONFamily,
			Limit:         limit,
			Enabled:       enabled,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return finalizeWorkflowRegistryOutput(deps.Strategy, report.Output.MCMSOperation, report.ToGenericReport())
}

// BatchSetUserDONOverrideInput configures multiple user DON overrides in a single MCMS BatchOperation.
//
// This is functionally equivalent to applying SetUserDONOverride once per override entry, but produces
// a single MCMS BatchOperation containing N transactions instead of N separate batch operations.
// Use this when you have many overrides to apply to the same WorkflowRegistry to keep proposals compact
// and avoid blocking other queued proposals.
type BatchSetUserDONOverrideInput struct {
	ChainSelector             uint64                              `json:"chainSelector"`
	WorkflowRegistryQualifier string                              `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Overrides                 []contracts.SetUserDONOverrideEntry `json:"overrides"`                 // Per-user overrides to apply
	MCMSConfig                *crecontracts.MCMSConfig            `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type BatchSetUserDONOverride struct{}

func (l BatchSetUserDONOverride) VerifyPreconditions(e cldf.Environment, config BatchSetUserDONOverrideInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	if len(config.Overrides) == 0 {
		return errors.New("must provide at least one override")
	}
	seen := make(map[string]struct{}, len(config.Overrides))
	for i, entry := range config.Overrides {
		if entry.User == (common.Address{}) {
			return fmt.Errorf("override[%d]: user address must not be zero", i)
		}
		if entry.DONFamily == "" {
			return fmt.Errorf("override[%d] (user %s): donFamily must not be empty", i, entry.User.Hex())
		}
		if _, _, err := entry.Resolve(); err != nil {
			return fmt.Errorf("override[%d] (user %s): %w", i, entry.User.Hex(), err)
		}
		key := entry.User.Hex() + "|" + entry.DONFamily
		if _, dup := seen[key]; dup {
			return fmt.Errorf("override[%d]: duplicate (user %s, donFamily %s)", i, entry.User.Hex(), entry.DONFamily)
		}
		seen[key] = struct{}{}
	}
	return nil
}

func (l BatchSetUserDONOverride) Apply(e cldf.Environment, config BatchSetUserDONOverrideInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// The strategy is only used at the BuildProposal step inside finalizeWorkflowRegistryOutput;
	// BatchSetUserDONOverrideOp itself builds calldata directly via deps.Chain + deps.Registry so
	// all overrides land in a single MCMS BatchOperation.
	deps, err := prepareWorkflowRegistryDeps(e, config.ChainSelector, config.WorkflowRegistryQualifier, config.MCMSConfig, contracts.BatchSetUserDONOverrideDescription)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.BatchSetUserDONOverrideOp, deps, contracts.BatchSetUserDONOverrideOpInput{
			ChainSelector: config.ChainSelector,
			Qualifier:     config.WorkflowRegistryQualifier,
			Overrides:     config.Overrides,
			MCMSConfig:    config.MCMSConfig,
		},
	)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return finalizeWorkflowRegistryOutput(deps.Strategy, report.Output.MCMSOperation, report.ToGenericReport())
}

// SetCapabilitiesRegistryInput configures the Capabilities registry address
type SetCapabilitiesRegistryInput struct {
	ChainSelector             uint64                   `json:"chainSelector"`
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
	Registry                  common.Address           `json:"registry"`                  // DON registry contract address
	ChainSelectorDON          uint64                   `json:"chainSelectorDON"`          // Chain selector where the DON registry exists
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
}

type SetCapabilitiesRegistry struct{}

func (l SetCapabilitiesRegistry) VerifyPreconditions(e cldf.Environment, config SetCapabilitiesRegistryInput) error {
	if err := validateRegistryTarget(config.ChainSelector, config.WorkflowRegistryQualifier); err != nil {
		return err
	}
	if config.Registry == (common.Address{}) {
		return errors.New("registry must be set: omitting it would have overwritten the capabilities registry address with the zero address")
	}
	if config.ChainSelectorDON == 0 {
		return errors.New("chainSelectorDON must be set: omitting it would have overwritten the capabilities registry chain selector with 0")
	}
	return nil
}

func (l SetCapabilitiesRegistry) Apply(e cldf.Environment, config SetCapabilitiesRegistryInput) (cldf.ChangesetOutput, error) {
	if err := l.VerifyPreconditions(e, config); err != nil {
		return cldf.ChangesetOutput{}, err
	}

	// Get MCMS contracts if needed
	var mcmsContracts *evmstate.MCMSWithTimelockState
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
		contracts.SetCapabilitiesRegistryDescription,
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
