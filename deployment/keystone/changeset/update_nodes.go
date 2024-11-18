package changeset

import (
	"fmt"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/internal"
)

var _ deployment.ChangeSet[*UpdateNodesRequest] = UpdateNodes

type UpdateNodesRequest = internal.UpdateNodesRequest
type NodeUpdate = internal.NodeUpdate

// UpdateNodes updates the a set of nodes.
// This a complex action in practice that involves registering missing capabilities, adding the nodes, and updating
// the capabilities of the DON
func UpdateNodes(env deployment.Environment, req *UpdateNodesRequest) (deployment.ChangesetOutput, error) {
	_, err := internal.UpdateNodes(env.Logger, req)
	if err != nil {
		return deployment.ChangesetOutput{}, fmt.Errorf("failed to update don: %w", err)
	}
	return deployment.ChangesetOutput{}, nil
}
