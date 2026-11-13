package queue

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

var _ ocr3_1types.ReportingPluginFactory[[]byte] = &consensusQueuePluginFactory{}

type consensusQueuePluginFactory struct {
	lggr               logger.Logger
	cfg                *PluginConfig
	triggerCoordinator triggers.Coordinator
	expirationDur      func(context.Context) (time.Duration, error)
}

func NewConsensusQueuePluginFactory(lggr logger.Logger, cfg *PluginConfig, limitsFactory limits.Factory, triggerCoordinator triggers.Coordinator) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	return &consensusQueuePluginFactory{
		lggr:               lggr,
		cfg:                cfg,
		triggerCoordinator: triggerCoordinator,
		expirationDur: func(ctx context.Context) (time.Duration, error) {
			//TODO also allow this grace period in the trigger coordinator
			return cresettings.Default.ConsensusQueueExpirationInterval.GetOrDefault(ctx, limitsFactory.Settings)
		},
	}, nil
}

func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, blobs ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := newConsensusQueuePlugin(c.lggr, config, c.triggerCoordinator, c.expirationDur, blobs)
	pluginInfo := ocr3_1types.ReportingPluginInfo1{
		Name:   "ConsensusQueueReportingPlugin",
		Limits: c.cfg.Limits,
	}
	return plugin, pluginInfo, plugin.Start(ctx)
}

var _ ocr3_1types.ReportingPlugin[[]byte] = &consensusQueuePlugin{}

type consensusQueuePlugin struct {
	services.Service
	eng                *services.Engine
	lggr               logger.Logger
	onchainCfg         ocr3types.ReportingPluginConfig
	triggerCoordinator triggers.Coordinator
	blobs              ocr3_1types.BlobBroadcastFetcher
	expirationDur      func(context.Context) (time.Duration, error)
}

func newConsensusQueuePlugin(
	lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig, triggerCoordinator triggers.Coordinator,
	expirationGraceDur func(context.Context) (time.Duration, error), blobs ocr3_1types.BlobBroadcastFetcher,
) *consensusQueuePlugin {
	c := &consensusQueuePlugin{
		lggr:               logger.Named(lggr, "ConsensusQueuePlugin"),
		onchainCfg:         onchainCfg,
		triggerCoordinator: triggerCoordinator,
		blobs:              blobs,
		expirationDur:      expirationGraceDur,
	}
	c.Service, c.eng = services.Config{
		Name: "ConsensusQueueReportingPlugin",
	}.NewServiceEngine(lggr)
	return c
}

func (c *consensusQueuePlugin) Query(_ context.Context, _ uint64, _ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return nil, nil
}

func (c *consensusQueuePlugin) Observation(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, kvState ocr3_1types.KeyValueStateReader,
	blobs ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Observation, error) {
	// Always attempt to start any executions from the KV Store in the event Committed was not called.
	// tryStartExecutions is idempotent.
	c.tryStartExecutions(kvState)

	capacity, err := c.triggerCoordinator.AvailableCapacity(ctx)
	if err != nil {
		//TODO or just log and make a best effort to observe events regardless?
		return nil, fmt.Errorf("failed to get available capacity: %w", err)
	}

	var events []triggers.CoordinatedEvent
	// TODO observe queue - but can't "take" - need to confirm during StateTransition

	expirationGraceDur, err := c.expirationDur(ctx)
	if err != nil {
		c.lggr.Warnw("Failed to get expiration grace duration - using default", "err", err)
	}
	expiration := time.Now().Add(expirationGraceDur)

	obs := &pb.ConsensusQueueObservation{
		Events:             make([]*pb.ConsensusQueueEvent, 0, len(events)),
		Expiration:         timestamppb.New(expiration),
		Capacity:           uint32(min(capacity.Global, math.MaxUint32)),
		WorkflowCapacities: make(map[string]uint32),
	}
	for wf, wfCap := range capacity.PerWorkflow {
		//TODO only for events proposed and queued in the KV store?
		obs.WorkflowCapacities[wf] = uint32(min(wfCap, math.MaxUint32))
	}

	const roundDur = 10 * time.Second                       //TODO dynamic calibration?
	expSeqNr := seqNr + uint64(expirationGraceDur/roundDur) //TODO pad extra time and/or rounds

	pbEvents := make([]*pb.CoordinatedEvent, len(events))
	for i, e := range events {
		pbEvents[i] = e.ToProto()
	}
	blobHandle, err := broadcastAsBlob(ctx, blobs, pbEvents, expSeqNr)
	if err != nil {
		return nil, err
	}

	const maxObservationSize = 0 //TODO from where?
	protoOpts := proto.MarshalOptions{Deterministic: true}
	for i, m := range pbEvents {
		// Try adding to batch

		observedAt, deadline := m.ObservedAt, m.Deadline
		m.ObservedAt, m.Deadline = nil, nil      // for consistent hashing
		protoBytes, err2 := protoOpts.Marshal(m) //TODO opt: don't marshal again, get hash from broadcast helper
		if err2 != nil {
			c.lggr.Errorw("failed to marshal event", "eventID", m.Event.Event.Id, "err", err2)
			continue
		}
		blobHash := sha256.Sum256(protoBytes)

		obs.Events = append(obs.Events, &pb.ConsensusQueueEvent{
			WorkflowID:     m.WorkflowID,
			TriggerEventID: m.Event.Event.Id,
			ObservedAt:     observedAt,
			Deadline:       deadline,
			BlobHash:       blobHash[:],
			Blobs: []*pb.CoordinatedEventBlobHandle{{
				Handle: blobHandle,
				Index:  uint32(i),
			}},
		})

		var addedCap bool
		if _, ok := obs.WorkflowCapacities[m.WorkflowID]; !ok {
			addedCap = true
			obs.WorkflowCapacities[m.WorkflowID] = uint32(min(capacity.PerWorkflow[m.WorkflowID], math.MaxUint32))
		}
		// Check size
		if protoOpts.Size(obs) > maxObservationSize {
			// too large, remove this one
			obs.Events = obs.Events[:len(obs.Events)-1]
			if addedCap {
				delete(obs.WorkflowCapacities, m.WorkflowID)
			}
			//TODO opt: remember the blobHandle to avoid rebroadcasting events
			break
		}
	}
	return protoOpts.Marshal(obs)
}

func (c *consensusQueuePlugin) ValidateObservation(
	ctx context.Context, _ uint64, _ ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, blobs ocr3_1types.BlobFetcher,
) error {
	var obs pb.ConsensusQueueObservation
	err := proto.Unmarshal(ao.Observation, &obs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal observation: %w", err)
	}
	var eg errgroup.Group
	for i, event := range obs.Events {
		eg.Go(func() error {
			if err := validateObservedEvent(ctx, blobs, event); err != nil {
				return fmt.Errorf("event %d failed validation: %w", i, err)
			}
			return nil
		})
	}
	return eg.Wait()
}

func validateObservedEvent(ctx context.Context, blobs ocr3_1types.BlobFetcher, event *pb.ConsensusQueueEvent) error {
	if event.WorkflowID == "" {
		return errors.New("event workflow id is missing")
	}
	if event.TriggerEventID == "" {
		return errors.New("event trigger event id is missing")
	}
	if event.ObservedAt == nil {
		return errors.New("event observation ts is missing")
	}
	if event.Deadline == nil {
		return errors.New("event deadline is missing")
	}
	if event.BlobHash == nil {
		return fmt.Errorf("event blob hash is missing")
	}

	if len(event.Blobs) == 0 {
		return fmt.Errorf("event blob is missing")
	}
	if len(event.Blobs) > 1 {
		return fmt.Errorf("event has more than one blob")
	}

	coordEvent, pbBytes, err := fetchFromBlob(ctx, blobs, event.Blobs[0].Handle, int(event.Blobs[0].Index))
	if err != nil {
		return fmt.Errorf("failed to fetch event from blob: %w", err)
	}

	if coordEvent.WorkflowID != event.WorkflowID {
		return fmt.Errorf("blob workflow ID %s does not match event workflow ID %s", coordEvent.WorkflowID, event.WorkflowID)
	}
	if coordEvent.GetEvent().GetEvent().Id != event.TriggerEventID {
		return fmt.Errorf("blob event ID %s does not match event trigger event ID %s", coordEvent.GetEvent().GetEvent().Id, event.TriggerEventID)
	}

	if coordEvent.ObservedAt != nil {
		return fmt.Errorf("blob observation timestamp is not nil: %v", coordEvent.ObservedAt)
	}
	if coordEvent.Deadline != nil {
		return fmt.Errorf("blob deadline timestamp is not nil: %v", coordEvent.Deadline)
	}

	//TODO reduce allocations with a pool?
	if shaSum := sha256.Sum256(pbBytes); !bytes.Equal(shaSum[:], event.BlobHash) {
		return fmt.Errorf("blob hash %x does not match observation hash %x", shaSum, event.BlobHash)
	}

	return nil
}

func (c *consensusQueuePlugin) ObservationQuorum(
	_ context.Context, _ uint64, _ ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher,
) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, c.onchainCfg.N, c.onchainCfg.F, aos), nil
}

func (c *consensusQueuePlugin) StateTransition(
	ctx context.Context, _ uint64, _ ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation,
	kvState ocr3_1types.KeyValueStateReadWriter, _ ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	var done chan struct{}
	var obs observations

	go func() {
		defer close(done)
		obs.aggregate(ctx, aos)
	}()

	stateRW := newStateReadWriter(kvState)
	state, err := stateRW.readObservedAndQueued()
	if err != nil {
		return nil, fmt.Errorf("failed to read state: %w", err)
	}

	// wait for aggregation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		if err = ctx.Err(); err != nil {
			return nil, err
		}
	}
	if obs.err != nil {
		return nil, fmt.Errorf("failed to aggregate observations: %w", obs.err)
	}
	if obs.warn != nil {
		c.lggr.Warnw("Failed to aggregate some observations", "err", obs.warn)
	}

	// update state
	state.pruneExpiredObservations(obs.medianExpiration)
	state.addObservedEvents(c.onchainCfg.F, obs.newObservedEvents)
	consensusEvents, err := state.queuedEvents.popLimited(obs.medianCap, obs.medianCapPerWF, obs.medianExpiration)
	if err != nil {
		return nil, err
	}
	//TODO might we hit a KV size limit in the middle of these? then what? can't just drop proposals
	// alternative: proposals only valid when there is definitely enough space
	err = stateRW.writeAll(state, consensusEvents)
	if err != nil {
		return nil, err
	}

	// no report to return, executions start when committed
	return ocr3_1types.ReportsPlusPrecursor{}, nil
}

type observations struct {
	medianExpiration  time.Time
	medianCap         uint32
	medianCapPerWF    map[string]uint32                          // wf ID -> capacity
	newObservedEvents map[string]*pb.ConsensusQueueObservedEvent // key: blobHash
	warn              error
	err               error
}

// TODO consider some configurable x% of capacity since some honest nodes have less
func (res *observations) aggregate(ctx context.Context, aos []ocrtypes.AttributedObservation) {
	res.medianCapPerWF = make(map[string]uint32)
	res.newObservedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)

	expirations := make([]*timestamppb.Timestamp, 0, len(aos))
	capacities := make([]uint32, 0, len(aos))
	wfCapacities := make(map[string][]uint32)
	for _, ao := range aos {
		var obs pb.ConsensusQueueObservation
		err := proto.Unmarshal(ao.Observation, &obs)
		if err != nil {
			res.warn = errors.Join(res.warn, fmt.Errorf("failed to unmarshal observation from oracle %d: %w", ao.Observation, err))
			continue
		}
		expirations = append(expirations, obs.Expiration)
		capacities = append(capacities, obs.Capacity)
		for wf, c := range obs.WorkflowCapacities {
			wfCapacities[wf] = append(wfCapacities[wf], c)
		}
		for _, event := range obs.Events {
			if ctx.Err() != nil {
				return
			}
			res.addEvent(ao.Observer, event)
		}
	}
	res.medianExpiration = medianTimestamp(expirations).AsTime()
	res.medianCap = median(capacities)
	for wf, caps := range wfCapacities {
		res.medianCapPerWF[wf] = median(caps)
	}
	return
}

func (res *observations) addEvent(obs commontypes.OracleID, event *pb.ConsensusQueueEvent) {
	aggregated, ok := res.newObservedEvents[string(event.BlobHash)]
	if !ok {
		aggregated = &pb.ConsensusQueueObservedEvent{
			TriggerEventID: event.TriggerEventID,
			WorkflowID:     event.WorkflowID,
			NodeIds:        []uint8{uint8(obs)},
			ObservedAts:    []*timestamppb.Timestamp{event.ObservedAt},
			Deadlines:      []*timestamppb.Timestamp{event.Deadline},
			Blobs:          []*pb.CoordinatedEventBlobHandle{event.Blobs[0]},
		}
	} else {
		// note: trigger event and workflow ID alignment is ensured by ValidateObservation
		aggregated.NodeIds = append(aggregated.NodeIds, uint8(obs))
		aggregated.ObservedAts = append(aggregated.ObservedAts, event.ObservedAt)
		aggregated.Deadlines = append(aggregated.Deadlines, event.Deadline)
		aggregated.Blobs = append(aggregated.Blobs, event.Blobs...)
	}

	res.newObservedEvents[string(event.BlobHash)] = aggregated
}

func (c *consensusQueuePlugin) Committed(ctx context.Context, _ uint64, kvState ocr3_1types.KeyValueStateReader) error {
	c.tryStartExecutions(kvState)
	return nil
}

// tryStartExecutions calls StartExecution asynchronously on the corresponding engine for each trigger event.
// It is idempotent, thus safe to call from Observation and Committed.
func (c *consensusQueuePlugin) tryStartExecutions(kvState ocr3_1types.KeyValueStateReader) {
	v, err := kvState.Read([]byte(consensusEventsKey))
	if err != nil {
		c.lggr.Errorw("Failed to read consensus events", "err", err)
		return
	}
	var consensusEvents []*pb.ConsensusQueueEvent
	if v != nil {
		var pbs *pb.ConsensusQueueEventList
		if err := proto.Unmarshal(v, pbs); err != nil {
			c.lggr.Errorw("Failed to unmarshal consensus events", "err", err)
			return
		}
	}

	var wg sync.WaitGroup
	for _, consensusEvent := range consensusEvents {
		wg.Add(1)
		c.eng.Go(func(ctx context.Context) {
			defer wg.Done()

			if err2 := c.tryStartExecution(ctx, consensusEvent); err2 != nil {
				c.lggr.Errorw("Failed to execute event", "eventID", consensusEvent.TriggerEventID, "err", err2)
				return
			}
		})
	}
	wg.Wait()
}

func (c *consensusQueuePlugin) tryStartExecution(ctx context.Context, consensusEvent *pb.ConsensusQueueEvent) error {
	// any blob will do, first that exists
	var pbEvent *pb.CoordinatedEvent
	var errs []error
	for _, blob := range consensusEvent.Blobs {
		var err error
		pbEvent, _, err = fetchFromBlob(ctx, c.blobs, blob.Handle, int(blob.Index))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to fetch blob: %w", err))
			c.lggr.Warnw("Failed to fetch blob", "handle", blob.Handle, "index", blob.Index, "err", err)
			continue // try next blob
		}
		break
	}
	if pbEvent == nil {
		return fmt.Errorf("failed to fetch blob from all %d handles: %w", len(consensusEvent.Blobs), errors.Join(errs...))
	}
	// fill in medianized timestamps
	pbEvent.ObservedAt = consensusEvent.ObservedAt
	pbEvent.Deadline = consensusEvent.Deadline

	var event triggers.CoordinatedEvent
	if err := event.FromProto(pbEvent); err != nil {
		return fmt.Errorf("failed to convert CoordinatedEvent from proto: %w", err)
	}

	// TODO execute StartExecution
	// c.triggerCoordinator.tryStartExecution(ctx, event)
	return nil
}

func (c *consensusQueuePlugin) Reports(ctx context.Context, _ uint64, _ ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, nil
}

func (c *consensusQueuePlugin) ShouldAcceptAttestedReport(ctx context.Context, _ uint64, _ ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}

func (c *consensusQueuePlugin) ShouldTransmitAcceptedReport(ctx context.Context, _ uint64, _ ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}
