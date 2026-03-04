package ring

import (
	"context"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

var _ ocr3types.ContractTransmitter[[]byte] = (*Transmitter)(nil)

// Transmitter handles transmission of shard orchestration outcomes
type Transmitter struct {
	lggr          logger.Logger
	ringStore     *Store
	arbiterScaler ringpb.ArbiterScalerClient
	fromAccount   types.Account
}

func NewTransmitter(lggr logger.Logger, ringStore *Store, arbiterScaler ringpb.ArbiterScalerClient, fromAccount types.Account) *Transmitter {
	return &Transmitter{
		lggr:          lggr,
		ringStore:     ringStore,
		arbiterScaler: arbiterScaler,
		fromAccount:   fromAccount,
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

	t.ringStore.SetRoutingState(outcome.State)

	applyRoutes := outcome.State == nil || IsInSteadyState(outcome.State)
	if !applyRoutes {
		if tr := outcome.State.GetTransition(); tr != nil && tr.WantShards == 0 {
			applyRoutes = true // fallback when no healthy shards
		}
	}
	if applyRoutes {
		flat := make(map[string]uint32, len(outcome.Routes))
		for wfID, route := range outcome.Routes {
			flat[wfID] = route.Shard
		}
		t.ringStore.SyncRoutes(flat)
		t.lggr.Debugw("Synced workflow shard mappings", "count", len(flat))
	} else {
		t.lggr.Debugw("Skipping route updates while in transition", "state", outcome.State)
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
