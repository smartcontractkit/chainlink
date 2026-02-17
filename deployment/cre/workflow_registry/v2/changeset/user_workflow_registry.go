package changeset

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"

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

// ChangeSet interface implementations for user functions
var _ cldf.ChangeSetV2[UserWorkflowUpsertInput] = UserWorkflowUpsert{}
var _ cldf.ChangeSetV2[UserWorkflowPauseInput] = UserWorkflowPause{}
var _ cldf.ChangeSetV2[UserWorkflowBatchPauseInput] = UserWorkflowBatchPause{}
var _ cldf.ChangeSetV2[UserWorkflowActivateInput] = UserWorkflowActivate{}
var _ cldf.ChangeSetV2[UserWorkflowDeleteInput] = UserWorkflowDelete{}
var _ cldf.ChangeSetV2[UserLinkOwnerInput] = UserLinkOwner{}
var _ cldf.ChangeSetV2[UserUnlinkOwnerInput] = UserUnlinkOwner{}
var _ cldf.ChangeSetV2[UserAllowlistRequestInput] = UserAllowlistRequest{}

func validateWorkflowIDHex(workflowID string) error {
	if workflowID == "" {
		return errors.New("workflow ID cannot be empty")
	}
	decoded, err := hex.DecodeString(workflowID)
	if err != nil || len(decoded) != 32 {
		return errors.New("workflow ID must be a valid 32-byte hex string")
	}
	return nil
}

// UserWorkflowUpsert creates or updates a user workflow
type UserWorkflowUpsert struct{}
type UserWorkflowUpsertInput struct {
	WorkflowID     string `json:"workflowID"`     // Workflow ID
	WorkflowName   string `json:"workflowName"`   // Workflow Name
	WorkflowTag    string `json:"workflowTag"`    // Workflow Tag
	WorkflowStatus uint8  `json:"workflowStatus"` // Workflow Status
	DonFamily      string `json:"donFamily"`      // DON Family
	BinaryURL      string `json:"binaryURL"`      // Binary URL
	ConfigURL      string `json:"configURL"`      // Config URL
	Attributes     string `json:"attributes"`     // Attributes
	KeepAlive      bool   `json:"keepAlive"`      // Keep Alive flag

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserWorkflowUpsert) VerifyPreconditions(e cldf.Environment, config UserWorkflowUpsertInput) error {
	if err := validateWorkflowIDHex(config.WorkflowID); err != nil {
		return err
	}
	if config.WorkflowName == "" {
		return errors.New("workflow name cannot be empty")
	}
	if config.DonFamily == "" {
		return errors.New("DON family cannot be empty")
	}
	if config.BinaryURL == "" {
		return errors.New("binary URL cannot be empty")
	}
	return nil
}

func (u UserWorkflowUpsert) Apply(e cldf.Environment, config UserWorkflowUpsertInput) (cldf.ChangesetOutput, error) {
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
		contracts.UpsertWorkflowUserDescription,
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
		contracts.UserUpsertWorkflowOp, deps, contracts.UserUpsertWorkflowOpInput{
			WorkflowID:     [32]byte(common.Hex2Bytes(config.WorkflowID)),
			WorkflowName:   config.WorkflowName,
			WorkflowTag:    config.WorkflowTag,
			WorkflowStatus: config.WorkflowStatus,
			DonFamily:      config.DonFamily,
			BinaryURL:      config.BinaryURL,
			ConfigURL:      config.ConfigURL,
			Attributes:     common.Hex2Bytes(config.Attributes),
			KeepAlive:      config.KeepAlive,
			ChainSelector:  config.ChainSelector,
			MCMSConfig:     config.MCMSConfig,
			Qualifier:      config.WorkflowRegistryQualifier,
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

// UserWorkflowPause pauses a user workflow
type UserWorkflowPause struct{}
type UserWorkflowPauseInput struct {
	WorkflowID string `json:"workflowID"` // Workflow ID

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserWorkflowPause) VerifyPreconditions(e cldf.Environment, config UserWorkflowPauseInput) error {
	return validateWorkflowIDHex(config.WorkflowID)
}

func (u UserWorkflowPause) Apply(e cldf.Environment, config UserWorkflowPauseInput) (cldf.ChangesetOutput, error) {
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
		contracts.PauseWorkflowUserDescription,
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
		contracts.UserPauseWorkflowOp, deps, contracts.UserPauseWorkflowOpInput{
			WorkflowID:    [32]byte(common.Hex2Bytes(config.WorkflowID)),
			ChainSelector: config.ChainSelector,
			MCMSConfig:    config.MCMSConfig,
			Qualifier:     config.WorkflowRegistryQualifier,
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

// UserWorkflowBatchPause pauses multiple user workflows
type UserWorkflowBatchPause struct{}
type UserWorkflowBatchPauseInput struct {
	WorkflowIDs     string `json:"workflowIDs"` // Workflow IDs
	workflowIDsByte [][32]byte

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserWorkflowBatchPause) VerifyPreconditions(e cldf.Environment, config UserWorkflowBatchPauseInput) error {
	workflowIDs := strings.Split(config.WorkflowIDs, ",")
	idMap := make(map[string]bool)
	for _, id := range workflowIDs {
		if _, exists := idMap[id]; exists {
			return fmt.Errorf("duplicate workflow ID detected: '%s'", id)
		}
		idMap[id] = true
		if err := validateWorkflowIDHex(id); err != nil {
			return fmt.Errorf("invalid workflow ID '%s': %w", id, err)
		}
	}
	return nil
}

func (u UserWorkflowBatchPause) Apply(e cldf.Environment, config UserWorkflowBatchPauseInput) (cldf.ChangesetOutput, error) {
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
		contracts.BatchPauseWorkflowUserDescription,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to create strategy: %w", err)
	}

	workflowIDs := strings.Split(config.WorkflowIDs, ",")
	var parsedIDs [][32]byte
	for _, id := range workflowIDs {
		var idArray [32]byte
		copy(idArray[:], common.Hex2Bytes(id))
		parsedIDs = append(parsedIDs, idArray)
	}
	config.workflowIDsByte = parsedIDs

	// Execute operation
	deps := contracts.WorkflowRegistryOpDeps{
		Env:      &e,
		Strategy: strategy,
		Registry: registry,
	}

	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		contracts.UserBatchPauseWorkflowOp, deps, contracts.UserBatchPauseWorkflowOpInput{
			WorkflowIDs:   config.workflowIDsByte,
			ChainSelector: config.ChainSelector,
			MCMSConfig:    config.MCMSConfig,
			Qualifier:     config.WorkflowRegistryQualifier,
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

// UserWorkflowActivate activates a user workflow
type UserWorkflowActivate struct{}
type UserWorkflowActivateInput struct {
	WorkflowID string `json:"workflowID"` // Workflow ID
	DonFamily  string `json:"donFamily"`  // DON Family

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserWorkflowActivate) VerifyPreconditions(e cldf.Environment, config UserWorkflowActivateInput) error {
	if config.DonFamily == "" {
		return errors.New("DON family cannot be empty")
	}
	return validateWorkflowIDHex(config.WorkflowID)
}

func (u UserWorkflowActivate) Apply(e cldf.Environment, config UserWorkflowActivateInput) (cldf.ChangesetOutput, error) {
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
		contracts.ActivateWorkflowUserDescription,
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
		contracts.UserActivateWorkflowOp, deps, contracts.UserActivateWorkflowOpInput{
			WorkflowID:    [32]byte(common.Hex2Bytes(config.WorkflowID)),
			DonFamily:     config.DonFamily,
			ChainSelector: config.ChainSelector,
			MCMSConfig:    config.MCMSConfig,
			Qualifier:     config.WorkflowRegistryQualifier,
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

// UserWorkflowDelete deletes a user workflow
type UserWorkflowDelete struct{}
type UserWorkflowDeleteInput struct {
	WorkflowID string `json:"workflowID"` // Workflow ID

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserWorkflowDelete) VerifyPreconditions(e cldf.Environment, config UserWorkflowDeleteInput) error {
	return validateWorkflowIDHex(config.WorkflowID)
}

func (u UserWorkflowDelete) Apply(e cldf.Environment, config UserWorkflowDeleteInput) (cldf.ChangesetOutput, error) {
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
		contracts.DeleteWorkflowUserDescription,
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
		contracts.UserDeleteWorkflowOp, deps, contracts.UserDeleteWorkflowOpInput{
			WorkflowID:    [32]byte(common.Hex2Bytes(config.WorkflowID)),
			ChainSelector: config.ChainSelector,
			MCMSConfig:    config.MCMSConfig,
			Qualifier:     config.WorkflowRegistryQualifier,
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

// UserLinkOwner links a user to a workflow owner
type UserLinkOwner struct{}

type UserLinkOwnerInput struct {
	ValidityTimestamp *big.Int `json:"validityTimestamp"`
	Proof             string   `json:"proof"`
	Signature         string   `json:"signature"`

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserLinkOwner) VerifyPreconditions(e cldf.Environment, config UserLinkOwnerInput) error {
	if config.ValidityTimestamp == nil || config.ValidityTimestamp.Cmp(big.NewInt(0)) == 0 {
		return errors.New("validity timestamp cannot be nil or zero")
	}
	if len(config.Proof) == 0 {
		return errors.New("proof cannot be empty")
	}
	if len(config.Signature) == 0 {
		return errors.New("signature cannot be empty")
	}
	return nil
}

func (u UserLinkOwner) Apply(e cldf.Environment, config UserLinkOwnerInput) (cldf.ChangesetOutput, error) {
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
		contracts.LinkOwnerUserDescription,
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
		contracts.UserLinkOwnerOp, deps, contracts.UserLinkOwnerOpInput{
			ValidityTimestamp: config.ValidityTimestamp,
			Proof:             [32]byte(common.Hex2Bytes(config.Proof)),
			Signature:         common.Hex2Bytes(config.Signature),
			ChainSelector:     config.ChainSelector,
			MCMSConfig:        config.MCMSConfig,
			Qualifier:         config.WorkflowRegistryQualifier,
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

// UserUnlinkOwner unlinks a user from a workflow owner
type UserUnlinkOwner struct{}

type UserUnlinkOwnerInput struct {
	Address           common.Address `json:"address"`
	ValidityTimestamp *big.Int       `json:"validityTimestamp"`
	Signature         string         `json:"signature"`

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserUnlinkOwner) VerifyPreconditions(e cldf.Environment, config UserUnlinkOwnerInput) error {
	if config.Address == (common.Address{}) {
		return errors.New("address cannot be zero address")
	}
	if config.ValidityTimestamp == nil || config.ValidityTimestamp.Cmp(big.NewInt(0)) == 0 {
		return errors.New("validity timestamp cannot be nil or zero")
	}
	if len(config.Signature) == 0 {
		return errors.New("signature cannot be empty")
	}
	return nil
}

func (u UserUnlinkOwner) Apply(e cldf.Environment, config UserUnlinkOwnerInput) (cldf.ChangesetOutput, error) {
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
		contracts.UnlinkOwnerUserDescription,
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
		contracts.UserUnlinkOwnerOp, deps, contracts.UserUnlinkOwnerOpInput{
			Address:           config.Address,
			ValidityTimestamp: config.ValidityTimestamp,
			Signature:         common.Hex2Bytes(config.Signature),
			ChainSelector:     config.ChainSelector,
			MCMSConfig:        config.MCMSConfig,
			Qualifier:         config.WorkflowRegistryQualifier,
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

// UserAllowlistRequest allows a user to request allowlist status
type UserAllowlistRequest struct{}

type UserAllowlistRequestInput struct {
	ExpiryTimestamp uint32 `json:"expiryTimestamp"`
	RequestDigest   string `json:"requestDigest"`

	ChainSelector             uint64                   `json:"chainSelector"`             // Chain Selector
	MCMSConfig                *crecontracts.MCMSConfig `json:"mcmsConfig,omitempty"`      // MCMS configuration
	WorkflowRegistryQualifier string                   `json:"workflowRegistryQualifier"` // Qualifier to identify the specific workflow registry
}

func (u UserAllowlistRequest) VerifyPreconditions(e cldf.Environment, config UserAllowlistRequestInput) error {
	if config.ExpiryTimestamp == 0 {
		return errors.New("expiry timestamp cannot be zero")
	}
	if len(config.RequestDigest) == 0 {
		return errors.New("request digest cannot be empty")
	}
	return nil
}

func (u UserAllowlistRequest) Apply(e cldf.Environment, config UserAllowlistRequestInput) (cldf.ChangesetOutput, error) {
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
		contracts.AllowlistRequest,
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
		contracts.UserAllowlistRequestOp, deps, contracts.UserAllowlistRequestOpInput{
			ExpiryTimestamp: config.ExpiryTimestamp,
			RequestDigest:   [32]byte(common.Hex2Bytes(config.RequestDigest)),
			ChainSelector:   config.ChainSelector,
			MCMSConfig:      config.MCMSConfig,
			Qualifier:       config.WorkflowRegistryQualifier,
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
