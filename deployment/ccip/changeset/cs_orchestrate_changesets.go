package changeset

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
)

// OrchestrateChangesets orchestrates the validation and application of multiple changesets.
var OrchestrateChangesets = cldf.CreateChangeSet(
	orchestrateChangesetsLogic,
	orchestrateChangesetsPrecondition,
)

// WithConfig is a struct that holds a changeset and its associated configuration.
// Changesets are applied in the provided order.
type WithConfig struct {
	Config    any
	ChangeSet cldf.ChangeSetV2[any]
}

// CreateGenericChangeSetWithConfig creates a ChangeSetWithConfig instance.
// It converts a strictly typed changeset with a specific configuration type C into a generic ChangeSetWithConfig.
// This allows for any changeset to be used with OrchestrateChangesets.
func CreateGenericChangeSetWithConfig[C any](changeSet cldf.ChangeSetV2[C], cfg C) WithConfig {
	applyFunc := func(e cldf.Environment, c any) (cldf.ChangesetOutput, error) {
		// Type assert the config to the expected type C
		configC, ok := c.(C)
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("config type assertion failed: expected %T, got %T", configC, c)
		}
		return changeSet.Apply(e, configC)
	}
	verifyFunc := func(e cldf.Environment, c any) error {
		// Type assert the config to the expected type C
		configC, ok := c.(C)
		if !ok {
			return fmt.Errorf("config type assertion failed: expected %T, got %T", configC, c)
		}
		return changeSet.VerifyPreconditions(e, configC)
	}
	return WithConfig{
		ChangeSet: cldf.CreateChangeSet(applyFunc, verifyFunc),
		Config:    cfg,
	}
}

// OrchestrateChangesetsConfig is the configuration struct for OrchestrateChangesets.
type OrchestrateChangesetsConfig struct {
	Description string
	MCMS        *proposalutils.TimelockConfig
	ChangeSets  []WithConfig
}

func orchestrateChangesetsLogic(e cldf.Environment, c OrchestrateChangesetsConfig) (cldf.ChangesetOutput, error) {
	state, err := stateview.LoadOnchainState(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	// Apply each changeset
	// NOTE: If a changeset fails to apply, we will return the output with reports only.
	finalOutput := cldf.ChangesetOutput{}
	for i, cs := range c.ChangeSets {
		output, err := cs.ChangeSet.Apply(e, cs.Config)
		if err != nil {
			finalOutput.Reports = append(finalOutput.Reports, output.Reports...)
			return cldf.ChangesetOutput{Reports: finalOutput.Reports}, fmt.Errorf("failed to apply changeset at index %d: %w", i, err)
		}
		err = mergeChangesetOutput(e, &finalOutput, output)
		if err != nil {
			finalOutput.Reports = append(finalOutput.Reports, output.Reports...)
			return cldf.ChangesetOutput{Reports: finalOutput.Reports}, fmt.Errorf("failed to merge output of changeset at index %d: %w", i, err)
		}
	}

	// Aggregate all Timelock proposals into 1 proposal
	proposal, err := proposalutils.AggregateProposalsV2(
		e,
		proposalutils.MCMSStates{
			MCMSEVMState:    state.EVMMCMSStateByChain(),
			MCMSSolanaState: state.SolanaMCMSStateByChain(e),
		},
		finalOutput.MCMSTimelockProposals,
		c.Description,
		c.MCMS,
	)
	if err != nil {
		return finalOutput, fmt.Errorf("failed to aggregate proposals: %w", err)
	}

	// If no proposal was created, we return the final output without a proposal
	if proposal == nil {
		return finalOutput, nil
	}

	// Reset proposals to only include the aggregated proposal
	finalOutput.MCMSTimelockProposals = []mcms.TimelockProposal{*proposal}
	return finalOutput, nil
}

func orchestrateChangesetsPrecondition(e cldf.Environment, c OrchestrateChangesetsConfig) error {
	if c.Description == "" {
		return errors.New("description must not be empty")
	}
	if c.MCMS == nil {
		return errors.New("mcms must not be nil")
	}
	for i, cs := range c.ChangeSets {
		if err := cs.ChangeSet.VerifyPreconditions(e, cs.Config); err != nil {
			return fmt.Errorf("precondition failed for changeset at index %d: %w", i, err)
		}
	}

	return nil
}

func mergeChangesetOutput(e cldf.Environment, dest *cldf.ChangesetOutput, src cldf.ChangesetOutput) error {
	err := cldf.MergeChangesetOutput(e, dest, src)
	if err != nil {
		return fmt.Errorf("failed to merge changeset output: %w", err)
	}

	// The following merges are not included in cldf.MergeChangesetOutput
	// TODO @ccip-tooling: Open PR in chainlink-deployments-framework to include these merges
	// 1. Merge DataStores
	if dest.DataStore == nil {
		dest.DataStore = src.DataStore
	} else if src.DataStore != nil {
		err := dest.DataStore.Merge(src.DataStore.Seal())
		if err != nil {
			return fmt.Errorf("failed to merge data store: %w", err)
		}
	}
	// 2. Merge Reports
	if dest.Reports == nil {
		dest.Reports = src.Reports
	} else if src.Reports != nil {
		dest.Reports = append(dest.Reports, src.Reports...)
	}

	return nil
}
