package shardownership

import (
	"context"

	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"

	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

type Verdict int

const (
	Allow Verdict = iota
	DenyNotOwner
	DenyOrchestratorError
)

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
