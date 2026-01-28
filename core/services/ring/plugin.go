package ring

import (
	"context"
	"errors"
	"slices"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
)

type Plugin struct {
	store         *Store
	arbiterScaler ringpb.ArbiterScalerClient
	config        ocr3types.ReportingPluginConfig
	lggr          logger.Logger

	batchSize  int
	timeToSync time.Duration
}

var _ ocr3types.ReportingPlugin[[]byte] = (*Plugin)(nil)

type ConsensusConfig struct {
	BatchSize  int
	TimeToSync time.Duration
}

const (
	DefaultBatchSize  = 100
	DefaultTimeToSync = 5 * time.Minute
)

func NewPlugin(store *Store, arbiterScaler ringpb.ArbiterScalerClient, config ocr3types.ReportingPluginConfig, lggr logger.Logger, cfg *ConsensusConfig) (*Plugin, error) {
	if arbiterScaler == nil {
		return nil, errors.New("RingOCR arbiterScaler is required")
	}
	if cfg == nil {
		cfg = &ConsensusConfig{
			BatchSize:  DefaultBatchSize,
			TimeToSync: DefaultTimeToSync,
		}
	}

	if cfg.BatchSize <= 0 {
		lggr.Infow("using default batchSize", "default", DefaultBatchSize)
		cfg.BatchSize = DefaultBatchSize
	}
	if cfg.TimeToSync <= 0 {
		lggr.Infow("using default timeToSync", "default", DefaultTimeToSync)
		cfg.TimeToSync = DefaultTimeToSync
	}

	lggr.Infow("RingPlugin config",
		"batchSize", cfg.BatchSize,
		"timeToSync", cfg.TimeToSync,
	)

	return &Plugin{
		store:         store,
		arbiterScaler: arbiterScaler,
		config:        config,
		lggr:          logger.Named(lggr, "RingPlugin"),
		batchSize:     cfg.BatchSize,
		timeToSync:    cfg.TimeToSync,
	}, nil
}

//coverage:ignore
func (p *Plugin) Query(_ context.Context, _ ocr3types.OutcomeContext) (types.Query, error) {
	return nil, nil
}

func (p *Plugin) Observation(ctx context.Context, _ ocr3types.OutcomeContext, _ types.Query) (types.Observation, error) {
	var wantShards uint32
	var shardStatus map[uint32]*ringpb.ShardStatus

	status, err := p.arbiterScaler.Status(ctx, &emptypb.Empty{})
	if err != nil {
		// NOTE: consider a fallback data source if Arbiter is not available
		p.lggr.Errorw("RingOCR failed to get arbiter scaler status", "error", err)
		return nil, err
	}
	wantShards = status.WantShards
	shardStatus = status.Status

	allWorkflowIDs := make([]string, 0)
	for wfID := range p.store.GetAllRoutingState() {
		allWorkflowIDs = append(allWorkflowIDs, wfID)
	}

	pendingAllocs := p.store.GetPendingAllocations()
	p.lggr.Infow("RingOCR Observation pending allocations", "pendingAllocs", pendingAllocs)

	allWorkflowIDs = append(allWorkflowIDs, pendingAllocs...)
	allWorkflowIDs = uniqueSorted(allWorkflowIDs)
	p.lggr.Infow("RingOCR Observation all workflow IDs unique", "allWorkflowIDs", allWorkflowIDs, "wantShards", wantShards)

	observation := &ringpb.Observation{
		ShardStatus: shardStatus,
		WorkflowIds: allWorkflowIDs,
		Now:         timestamppb.Now(),
		WantShards:  wantShards,
	}

	return proto.MarshalOptions{Deterministic: true}.Marshal(observation)
}

func (p *Plugin) ValidateObservation(_ context.Context, _ ocr3types.OutcomeContext, _ types.Query, ao types.AttributedObservation) error {
	observation := &ringpb.Observation{}
	if err := proto.Unmarshal(ao.Observation, observation); err != nil {
		return err
	}
	if observation.Now == nil {
		return errors.New("observation missing timestamp")
	}
	if observation.WantShards == 0 {
		return errors.New("observation missing WantShards")
	}
	return nil
}

func (p *Plugin) ObservationQuorum(_ context.Context, _ ocr3types.OutcomeContext, _ types.Query, aos []types.AttributedObservation) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, p.config.N, p.config.F, aos), nil
}

func (p *Plugin) collectShardInfo(aos []types.AttributedObservation) (shardHealth map[uint32]int, workflows []string, timestamps []time.Time, wantShardVotes map[commontypes.OracleID]uint32) {
	shardHealth = make(map[uint32]int)
	wantShardVotes = make(map[commontypes.OracleID]uint32)
	for _, ao := range aos {
		observation := &ringpb.Observation{}
		_ = proto.Unmarshal(ao.Observation, observation) // validated in ValidateObservation

		for shardID, status := range observation.ShardStatus {
			if status != nil && status.IsHealthy {
				shardHealth[shardID]++
			}
		}

		workflows = append(workflows, observation.WorkflowIds...)
		timestamps = append(timestamps, observation.Now.AsTime())

		wantShardVotes[ao.Observer] = observation.WantShards
	}
	return shardHealth, workflows, timestamps, wantShardVotes
}

func (p *Plugin) getHealthyShards(shardHealth map[uint32]int) []uint32 {
	var healthyShards []uint32
	for shardID, votes := range shardHealth {
		if votes > p.config.F {
			healthyShards = append(healthyShards, shardID)
			p.store.SetShardHealth(shardID, true)
		}
	}
	slices.Sort(healthyShards)

	return healthyShards
}

func (p *Plugin) Outcome(_ context.Context, outctx ocr3types.OutcomeContext, _ types.Query, aos []types.AttributedObservation) (ocr3types.Outcome, error) {
	currentShardHealth, allWorkflows, nows, wantShardVotes := p.collectShardInfo(aos)
	p.lggr.Infow("RingOCR Outcome collect shard info", "currentShardHealth", currentShardHealth, "wantShardVotes", wantShardVotes)

	// Use the median timestamp to determine the current time
	slices.SortFunc(nows, time.Time.Compare)
	now := nows[len(nows)/2]

	// Use median for wantShards consensus (all validated observations have WantShards > 0)
	votes := make([]uint32, 0, len(wantShardVotes))
	for _, v := range wantShardVotes {
		votes = append(votes, v)
	}
	slices.Sort(votes)
	wantShards := votes[len(votes)/2]

	// Bootstrap from Arbiter's current shard count on 1st round; subsequent rounds build on prior outcome
	prior := &ringpb.Outcome{}
	if outctx.PreviousOutcome == nil {
		prior.Routes = map[string]*ringpb.WorkflowRoute{}
		prior.State = &ringpb.RoutingState{Id: outctx.SeqNr, State: &ringpb.RoutingState_RoutableShards{RoutableShards: wantShards}}
	} else if err := proto.Unmarshal(outctx.PreviousOutcome, prior); err != nil {
		return nil, err
	}

	allWorkflows = uniqueSorted(allWorkflows)

	healthyShards := p.getHealthyShards(currentShardHealth)

	nextState, err := NextState(prior.State, wantShards, now, p.timeToSync)
	if err != nil {
		return nil, err
	}

	// Deterministic hashing ensures all nodes agree on workflow-to-shard assignments
	// without coordination, preventing protocol failures from inconsistent routing
	ring := newShardRing(healthyShards)
	routes := make(map[string]*ringpb.WorkflowRoute)
	for _, wfID := range allWorkflows {
		shard, err := locateShard(ring, wfID)
		if err != nil {
			p.lggr.Warnw("RingOCR failed to locate shard for workflow", "workflowID", wfID, "error", err)
			shard = 0 // fallback to shard 0 when no healthy shards
		}
		routes[wfID] = &ringpb.WorkflowRoute{Shard: shard}
	}

	outcome := &ringpb.Outcome{
		State:  nextState,
		Routes: routes,
	}

	p.lggr.Infow("RingOCR Outcome", "healthyShards", len(healthyShards), "totalObservations", len(aos), "workflowCount", len(routes))

	return proto.MarshalOptions{Deterministic: true}.Marshal(outcome)
}

func (p *Plugin) Reports(_ context.Context, _ uint64, outcome ocr3types.Outcome) ([]ocr3types.ReportPlus[[]byte], error) {
	allOraclesTransmitNow := &ocr3types.TransmissionSchedule{
		Transmitters:       make([]commontypes.OracleID, p.config.N),
		TransmissionDelays: make([]time.Duration, p.config.N),
	}

	for i := 0; i < p.config.N; i++ {
		allOraclesTransmitNow.Transmitters[i] = commontypes.OracleID(i) //nolint:gosec // G115: i bounded by config.N
	}

	info, err := structpb.NewStruct(map[string]any{
		"keyBundleName": "evm",
	})
	if err != nil {
		return nil, err
	}
	infoBytes, err := proto.MarshalOptions{Deterministic: true}.Marshal(info)
	if err != nil {
		return nil, err
	}

	return []ocr3types.ReportPlus[[]byte]{
		{
			ReportWithInfo: ocr3types.ReportWithInfo[[]byte]{
				Report: types.Report(outcome),
				Info:   infoBytes,
			},
		},
	}, nil
}

//coverage:ignore
func (p *Plugin) ShouldAcceptAttestedReport(_ context.Context, _ uint64, _ ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

//coverage:ignore
func (p *Plugin) ShouldTransmitAcceptedReport(_ context.Context, _ uint64, _ ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

//coverage:ignore
func (p *Plugin) Close() error {
	return nil
}
