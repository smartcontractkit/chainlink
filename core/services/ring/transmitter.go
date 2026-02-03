package ring

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*Transmitter)(nil)

// Transmitter handles transmission of shard orchestration outcomes
type Transmitter struct {
	lggr                   logger.Logger
	ringStore              *Store
	shardOrchestratorStore *shardorchestrator.Store
	arbiterScaler          ringpb.ArbiterScalerClient
	fromAccount            types.Account
}

func NewTransmitter(lggr logger.Logger, ringStore *Store, shardOrchestratorStore *shardorchestrator.Store, arbiterScaler ringpb.ArbiterScalerClient, fromAccount types.Account) *Transmitter {
	return &Transmitter{
		lggr:                   lggr,
		ringStore:              ringStore,
		shardOrchestratorStore: shardOrchestratorStore,
		arbiterScaler:          arbiterScaler,
		fromAccount:            fromAccount,
	}
}

func (t *Transmitter) Transmit(ctx context.Context, _ types.ConfigDigest, _ uint64, r ocr3types.ReportWithInfo[[]byte], _ []types.AttributedOnchainSignature) error {
	outcome := &ringpb.Outcome{}
	if err := proto.Unmarshal(r.Report, outcome); err != nil {
		t.lggr.Error("failed to unmarshal report")
		return err
	}

	if err := t.notifyArbiter(ctx, outcome.State); err != nil {
		t.lggr.Errorw("failed to notify arbiter", "err", err)
		return err
	}

	// Update Ring Store
	t.ringStore.SetRoutingState(outcome.State)

	// Determine if system is in transition state
	systemInTransition := false
	if outcome.State != nil {
		if _, ok := outcome.State.State.(*ringpb.RoutingState_Transition); ok {
			systemInTransition = true
		}
	}

	// Update ShardOrchestrator store if available
	if t.shardOrchestratorStore != nil {
		mappings := make([]*shardorchestrator.WorkflowMappingState, 0, len(outcome.Routes))
		for workflowID, route := range outcome.Routes {
			// Get the current shard assignment for this workflow to detect changes
			var oldShardID uint32
			var transitionState shardorchestrator.TransitionState

			existingMapping, err := t.shardOrchestratorStore.GetWorkflowMapping(ctx, workflowID)
			switch {
			case err != nil:
				// New workflow - no previous assignment
				oldShardID = 0
				transitionState = shardorchestrator.StateSteady
			case existingMapping.NewShardID != route.Shard:
				// Workflow is moving to a different shard
				oldShardID = existingMapping.NewShardID
				transitionState = shardorchestrator.StateTransitioning
			default:
				// Same shard - but might be in system transition
				oldShardID = existingMapping.NewShardID
				if systemInTransition {
					transitionState = shardorchestrator.StateTransitioning
				} else {
					transitionState = shardorchestrator.StateSteady
				}
			}

			mappings = append(mappings, &shardorchestrator.WorkflowMappingState{
				WorkflowID:      workflowID,
				OldShardID:      oldShardID,
				NewShardID:      route.Shard,
				TransitionState: transitionState,
			})
		}

		if err := t.shardOrchestratorStore.BatchUpdateWorkflowMappings(ctx, mappings); err != nil {
			t.lggr.Errorw("failed to update ShardOrchestrator store", "err", err, "workflowCount", len(mappings))
			// Don't fail the entire transmission if ShardOrchestrator update fails
		} else {
			t.lggr.Debugw("Updated ShardOrchestrator store", "workflowCount", len(mappings))
		}
	}

	// Update Ring Store workflow mappings
	for workflowID, route := range outcome.Routes {
		t.ringStore.SetShardForWorkflow(workflowID, route.Shard)
		t.lggr.Debugw("Updated workflow shard mapping", "workflowID", workflowID, "shard", route.Shard)
	}

	return nil
}

func (t *Transmitter) notifyArbiter(ctx context.Context, state *ringpb.RoutingState) error {
	if state == nil {
		return nil
	}

	var nShards uint32
	switch s := state.State.(type) {
	case *ringpb.RoutingState_RoutableShards:
		nShards = s.RoutableShards
		t.lggr.Infow("Transmitting shard routing", "routableShards", nShards)
	case *ringpb.RoutingState_Transition:
		nShards = s.Transition.WantShards
		t.lggr.Infow("Transmitting shard routing (in transition)", "wantShards", nShards)
	}

	if t.arbiterScaler != nil && nShards > 0 {
		if _, err := t.arbiterScaler.ConsensusWantShards(ctx, &ringpb.ConsensusWantShardsRequest{NShards: nShards}); err != nil {
			return err
		}
	}

	return nil
}

func (t *Transmitter) FromAccount(ctx context.Context) (types.Account, error) {
	return t.fromAccount, nil
}
