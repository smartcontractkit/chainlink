package queue

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"slices"
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
	eng                                 *services.Engine
	lggr                                logger.Logger
	onchainCfg                          ocr3types.ReportingPluginConfig
	blobHandleMarshalledBytesUpperBound int
	triggerCoordinator                  triggers.Coordinator
	blobBroadcastFetcher                ocr3_1types.BlobBroadcastFetcher
	expirationDur                       func(context.Context) (time.Duration, error)
}

func newConsensusQueuePlugin(
	lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig, triggerCoordinator triggers.Coordinator,
	expirationGraceDur func(context.Context) (time.Duration, error), blobs ocr3_1types.BlobBroadcastFetcher,
) *consensusQueuePlugin {
	c := &consensusQueuePlugin{
		lggr:                                logger.Named(lggr, "ConsensusQueuePlugin"),
		onchainCfg:                          onchainCfg,
		blobHandleMarshalledBytesUpperBound: ocr3_1types.BlobHandleMarshalledBytesUpperBound(onchainCfg.F, onchainCfg.N),
		triggerCoordinator:                  triggerCoordinator,
		blobBroadcastFetcher:                blobs,
		expirationDur:                       expirationGraceDur,
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
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, kvReader ocr3_1types.KeyValueStateReader,
	blobBroadcaster ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Observation, error) {
	// Always attempt to start any executions from the KV Store in the event Committed was not called.
	// tryStartExecutions is idempotent.
	c.tryStartExecutions(kvReader)

	capacity, err := c.triggerCoordinator.AvailableCapacity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available capacity: %w", err)
	}

	var events []triggers.CoordinatedEvent
	// TODO observe queue - but can't "take" - need to confirm during StateTransition
	// ignore if observed in kv state already

	expirationGraceDur, err := c.expirationDur(ctx)
	if err != nil {
		c.lggr.Warnw("Failed to get expiration grace duration - using default", "err", err)
	}

	emptyBlobHandle := make([]byte, c.blobHandleMarshalledBytesUpperBound)
	obs := &pb.ConsensusQueueObservation{
		Events: make([]*pb.ConsensusQueueEvent, 0, len(events)),
		// placeholder for accurate size calculation
		BlobHandles:        slices.Repeat([][]byte{emptyBlobHandle}, maxObservedBlobHandles),
		Expiration:         timestamppb.New(time.Now().Add(-expirationGraceDur)),
		Capacity:           uint32(min(capacity.Global, math.MaxUint32)),
		WorkflowCapacities: make(map[string]uint32),
	}
	for wf, wfCap := range capacity.PerWorkflow {
		//TODO only for events proposed and queued in the KV store?
		obs.WorkflowCapacities[wf] = uint32(min(wfCap, math.MaxUint32))
	}

	coordEvents := make([]*pb.CoordinatedEvent, len(events))
	for i, e := range events {
		coordEvents[i] = e.ToProto()
	}
	protoOpts := proto.MarshalOptions{Deterministic: true}
	for i, coord := range coordEvents { //TODO iter.Seq from trigger coordinator
		// Try adding to batch

		observedAt, deadline := coord.ObservedAt, coord.Deadline
		coord.ObservedAt, coord.Deadline = nil, nil  // for consistent hashing
		protoBytes, err2 := protoOpts.Marshal(coord) // TODO opt: avoid remarshalling
		if err2 != nil {
			c.lggr.Errorw("failed to marshal event", "eventID", coord.Event.Event.Id, "err", err2)
			continue
		}
		timelessHash := sha256.Sum256(protoBytes)

		obs.Events = append(obs.Events, &pb.ConsensusQueueEvent{
			WorkflowID:     coord.WorkflowID,
			TriggerEventID: coord.Event.Event.Id,
			ObservedAt:     observedAt,
			Deadline:       deadline,
			TimelessHash:   timelessHash[:],
			Blobs: []*pb.CoordinatedEventBlobHandle{{
				Blob:    0,
				Element: uint32(i),
			}},
		})

		var addedCap bool
		if _, ok := obs.WorkflowCapacities[coord.WorkflowID]; !ok {
			addedCap = true
			obs.WorkflowCapacities[coord.WorkflowID] = uint32(min(capacity.PerWorkflow[coord.WorkflowID], math.MaxUint32))
		}

		// Check size
		if protoOpts.Size(obs) > maxObservationBytes {
			// too large, remove this one
			obs.Events = obs.Events[:len(obs.Events)-1]
			if addedCap {
				delete(obs.WorkflowCapacities, coord.WorkflowID)
			}
			break
		}
	}

	roundDur := c.onchainCfg.EstimatedRoundInterval         //TODO OK?
	expSeqNr := seqNr + uint64(expirationGraceDur/roundDur) //TODO pad extra time and/or rounds?

	obs.BlobHandles = nil // remove the placeholders
	if len(coordEvents) > 0 {
		blobHandle, n, err2 := broadcastAsBlob(ctx, blobBroadcaster, coordEvents, expSeqNr)
		if err2 != nil {
			return nil, err2
		}
		obs.BlobHandles = append(obs.BlobHandles, blobHandle)
		// note: elements already reference blob 0 and natural element index is aligned

		var offset int
		for batch := 1; n < len(coordEvents); batch++ {
			offset += n
			coordEvents = coordEvents[n:]
			blobHandle, n, err2 = broadcastAsBlob(ctx, blobBroadcaster, coordEvents, expSeqNr)
			if err2 != nil {
				return nil, err2
			}
			obs.BlobHandles = append(obs.BlobHandles, blobHandle)
			for i := n; i < n; i++ { // update indices
				obs.Events[offset+i].Blobs[0].Blob = uint32(batch)
				obs.Events[offset+i].Blobs[0].Element = uint32(i)
			}
		}
	}

	return protoOpts.Marshal(obs)
}

func (c *consensusQueuePlugin) ValidateObservation(
	ctx context.Context, _ uint64, _ ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher,
) error {
	var obs pb.ConsensusQueueObservation
	err := proto.Unmarshal(ao.Observation, &obs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal observation: %w", err)
	}
	if len(obs.BlobHandles) > maxObservedBlobHandles {
		return fmt.Errorf("blob handles are limited to %d but got %d", maxObservedBlobHandles, len(obs.BlobHandles))
	}
	if obs.Expiration == nil {
		return errors.New("expiration is missing")
	}
	//TODO what do we expect from WF capacities? check KV state and insist 100% match? or subset - none extra?

	if len(obs.Events) > maxObservedEvents {
		return fmt.Errorf("events are limited to %d but got %d", maxObservedEvents, len(obs.Events))
	}
	if !slices.IsSortedFunc(obs.Events, cmpConsensusQueueEvents) {
		return fmt.Errorf("events must be sorted")
	}
	var eg errgroup.Group
	fetcher := blobCache{fetcher: blobFetcher, cache: make(map[string]func() ([]byte, error))}
	for i, event := range obs.Events {
		eg.Go(func() error {
			if err := validateObservedEvent(ctx, &fetcher, event, obs.BlobHandles); err != nil {
				return fmt.Errorf("event %d failed validation: %w", i, err)
			}
			return nil
		})
	}
	return eg.Wait()
}

func validateObservedEvent(ctx context.Context, fetcher *blobCache, event *pb.ConsensusQueueEvent, blobHandles [][]byte) error {
	if event.WorkflowID == "" {
		return errors.New("event workflow id is missing")
	} ///TODO check not too long - exact length?
	if event.TriggerEventID == "" {
		return errors.New("event trigger event id is missing")
	} ///TODO check not too long - exact length?
	if event.ObservedAt == nil {
		return errors.New("event observation ts is missing")
	}
	if event.Deadline == nil {
		return errors.New("event deadline is missing")
	}
	if event.TimelessHash == nil {
		return fmt.Errorf("event blob hash is missing")
	} ///TODO check not too long - exact length?

	if len(event.Blobs) == 0 {
		return fmt.Errorf("event blob is missing")
	}
	if len(event.Blobs) > 1 {
		return fmt.Errorf("event has more than one blob")
	}

	blobIdx := int(event.Blobs[0].Blob)
	if blobIdx >= len(blobHandles) {
		return fmt.Errorf("event blob index %d is out of range", blobIdx)
	}
	blobHandle := blobHandles[blobIdx]

	coordEvent, pbBytes, err := fetchCoordinatedEventFromBlob(ctx, fetcher, blobHandle, int(event.Blobs[0].Element))
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

	//TODO opt: reduce allocations with a pool
	if shaSum := sha256.Sum256(pbBytes); !bytes.Equal(shaSum[:], event.TimelessHash) {
		return fmt.Errorf("blob hash %x does not match observation hash %x", shaSum, event.TimelessHash)
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
	kvReaderWriter ocr3_1types.KeyValueStateReadWriter, _ ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	var done chan struct{}
	var obs observations

	go func() {
		defer close(done)
		obs.aggregate(ctx, aos)
	}()

	var state kvState
	err := state.loadFrom(kvReaderWriter)
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
	state.addObservedEvents(c.onchainCfg.F, obs.newObservedEvents, obs.observedBlobs, func(queued *pb.ConsensusQueueEvent) {
		//TODO ack the trigger coordinator - no longer need to observe this event
	})
	state.consensusEvents = state.queuedEvents.popLimited(obs.medianCap, obs.medianCapPerWF, obs.medianExpiration)

	// prune observations until we fit
	for state.sizeOf() > maxKeyValueModifiedKeysPlusValuesBytes {
		//TODO drop observations (must be deterministic, fair, efficient)
		// prefer farthest future deadlines and fewest observations
	}

	err = state.writeTo(&stateWriter{kv: kvReaderWriter})
	if err != nil {
		return nil, err
	}

	// no report to return, executions start when committed
	return nil, nil
}

type observations struct {
	medianExpiration  time.Time
	medianCap         uint32
	medianCapPerWF    map[string]uint32                          // wf ID -> capacity
	newObservedEvents map[string]*pb.ConsensusQueueObservedEvent // key: timeless hash
	observedBlobs     blobHandleMap
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
			res.addEvent(ao.Observer, event, obs.BlobHandles)
		}
	}
	res.medianExpiration = medianTimestamp(expirations).AsTime()
	res.medianCap = median(capacities)
	for wf, caps := range wfCapacities {
		res.medianCapPerWF[wf] = median(caps)
	}
	return
}

func (res *observations) addEvent(obs commontypes.OracleID, event *pb.ConsensusQueueEvent, blobHandles [][]byte) {
	res.observedBlobs.add(string(event.TimelessHash), string(blobHandles[event.Blobs[0].Blob]), event.Blobs[0].Element)
	aggregated, ok := res.newObservedEvents[string(event.TimelessHash)]
	if !ok {
		res.newObservedEvents[string(event.TimelessHash)] = &pb.ConsensusQueueObservedEvent{
			TriggerEventID: event.TriggerEventID,
			WorkflowID:     event.WorkflowID,
			NodeIds:        []uint8{uint8(obs)},
			ObservedAts:    []*timestamppb.Timestamp{event.ObservedAt},
			Deadlines:      []*timestamppb.Timestamp{event.Deadline},
			Blobs:          nil, // res.observedBlobs tracks these while in memory
		}
	} else {
		// note: trigger event and workflow ID alignment is ensured by ValidateObservation
		aggregated.NodeIds = append(aggregated.NodeIds, uint8(obs))
		aggregated.ObservedAts = append(aggregated.ObservedAts, event.ObservedAt)
		aggregated.Deadlines = append(aggregated.Deadlines, event.Deadline)
		// note: res.observedBlobs tracks blobs while in memory
	}
}

func (c *consensusQueuePlugin) Committed(ctx context.Context, _ uint64, kvReader ocr3_1types.KeyValueStateReader) error {
	c.tryStartExecutions(kvReader)
	return nil
}

// tryStartExecutions calls StartExecution asynchronously on the corresponding engine for each trigger event.
// It is idempotent, thus safe to call from Observation and Committed.
func (c *consensusQueuePlugin) tryStartExecutions(kvReader ocr3_1types.KeyValueStateReader) {
	consensusEvents, blobHandles, err := readConsensusState(kvReader)
	if err != nil {
		c.lggr.Errorw("Failed to read consensus state. Unable to execute events", "err", err)
		return
	}

	fetcher := blobCache{fetcher: c.blobBroadcastFetcher, cache: make(map[string]func() ([]byte, error))}
	var wg sync.WaitGroup
	for _, consensusEvent := range consensusEvents {
		wg.Add(1)
		c.eng.Go(func(ctx context.Context) {
			defer wg.Done()

			if err2 := c.tryStartExecution(ctx, &fetcher, consensusEvent, blobHandles); err2 != nil {
				c.lggr.Errorw("Failed to execute event", "eventID", consensusEvent.TriggerEventID, "err", err2)
				return
			}
		})
	}
	wg.Wait()
}

func (c *consensusQueuePlugin) tryStartExecution(ctx context.Context, fetcher *blobCache, consensusEvent *pb.ConsensusQueueEvent, blobHandles [][]byte) error {
	// any blob will do, first that exists
	var pbEvent *pb.CoordinatedEvent
	var errs []error
	for _, blob := range consensusEvent.Blobs {
		var err error
		if int(blob.Blob) > len(blobHandles) {
			c.lggr.Errorw("Blob handle index out of bounds", "index", int(blob.Blob), "len", len(blobHandles))
			continue
		}
		blobHandle := blobHandles[int(blob.Blob)]
		pbEvent, _, err = fetchCoordinatedEventFromBlob(ctx, fetcher, blobHandle, int(blob.Element))
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to fetch blob: %w", err))
			c.lggr.Errorw("Failed to fetch blob", "handle", blobHandle, "index", blob.Element, "err", err)
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
