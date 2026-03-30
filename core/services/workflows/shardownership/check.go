// Package shardownership validates this node's shard against the shard-0 ring store via the orchestrator.
package shardownership

import (
	"context"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"

	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

// Verdict is the outcome of a committed-ownership check against GetWorkflowShardMapping.
type Verdict int

const (
	// Allow means this shard matches the orchestrator mapping for the workflow.
	Allow Verdict = iota
	// DenyNotOwner means mapping exists but assigns the workflow to another shard, or the workflow is missing from mappings.
	DenyNotOwner
	// DenyOrchestratorError means the orchestrator RPC failed (caller should fail closed).
	DenyOrchestratorError
)

// CheckCommittedOwner returns whether this node may execute the workflow per shard-0's committed routes.
func CheckCommittedOwner(ctx context.Context, client shardorchestrator.ClientInterface, workflowID string, myShardID uint32) (v Verdict, resp *ringpb.GetWorkflowShardMappingResponse, err error) {
	resp, err = client.GetWorkflowShardMapping(ctx, []string{workflowID})
	if err != nil {
		return DenyOrchestratorError, nil, err
	}
	shard, ok := resp.Mappings[workflowID]
	if !ok {
		return DenyNotOwner, resp, nil
	}
	if shard != myShardID {
		return DenyNotOwner, resp, nil
	}
	return Allow, resp, nil
}
