package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink-deployments-framework/deployment"
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
