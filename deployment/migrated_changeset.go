package deployment

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// ChangeSets

type (
	ChangeSet[C any]            = deployment.ChangeSet[C]
	ChangeLogic[C any]          = deployment.ChangeLogic[C]
	PreconditionVerifier[C any] = deployment.PreconditionVerifier[C]
	ChangeSetV2[C any]          = deployment.ChangeSetV2[C]
)

var (
	ErrInvalidConfig      = deployment.ErrInvalidConfig
	ErrInvalidEnvironment = deployment.ErrInvalidEnvironment
)

// Changeset Output

type (
	ChangesetOutput = deployment.ChangesetOutput
	ProposedJob     = deployment.ProposedJob
)

var (
	MergeChangesetOutput = deployment.MergeChangesetOutput
)

// ViewState

type (
	ViewState   = deployment.ViewState
	ViewStateV2 = deployment.ViewStateV2
)

type WrappedChangeSet[C any] struct {
	operation deployment.ChangeSetV2[C]
}

// RunChangeset is used to run a changeset in another changeset
// It executes VerifyPreconditions internally to handle changeset errors.
func RunChangeset[C any](
	operation deployment.ChangeSetV2[C],
	env deployment.Environment,
	config C,
) (deployment.ChangesetOutput, error) {
	cs := WrappedChangeSet[C]{operation: operation}

	err := cs.operation.VerifyPreconditions(env, config)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to run precondition: %w", err)
	}

	return cs.operation.Apply(env, config)
}
