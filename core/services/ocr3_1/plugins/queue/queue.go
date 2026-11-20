// Package queue contains the consensus queue OCR 3.1 reporting plugin, for coordinating execution of triggered events.
// Each node has a local [triggers.Queue] of [triggers.CoordinatedEvent], which are observed (and broadcast) by the plugin.
// When consensus is reached, events are put in to the KeyValueState priority queue (and removed from the [triggers.Queue]).
// The priority queue is ordered by deadline (with nanosecond precision) and then by trigger event ID.
// Each round, as many events as possible are popped from the priority queue and then (fetched and) executed via [triggers.Coordinator.Deliver].
// Capacity, both overall and per-workflow, limits what may execute. Capacities are observed and medianized.
package queue

import (
	"bytes"
	"container/heap"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"iter"
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
	triggerQueue       triggers.Queue
	expirationDur      func(context.Context) (time.Duration, error)
}

func NewConsensusQueuePluginFactory(
	lggr logger.Logger,
	cfg *PluginConfig,
	limitsFactory limits.Factory,
	triggerCoordinator triggers.Coordinator,
	triggerQueue triggers.Queue,
) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	if cfg == nil {
		return nil, errors.New("cfg cannot be nil")
	}
	if triggerCoordinator == nil {
		return nil, errors.New("triggerCoordinator cannot be nil")
	}
	if triggerQueue == nil {
		return nil, errors.New("triggerQueue cannot be nil")
	}
	return &consensusQueuePluginFactory{
		lggr:               lggr,
		cfg:                cfg,
		triggerCoordinator: triggerCoordinator,
		triggerQueue:       triggerQueue,
		expirationDur: func(ctx context.Context) (time.Duration, error) {
			//TODO also allow this grace period in the trigger coordinator
			return cresettings.Default.ConsensusQueueExpirationInterval.GetOrDefault(ctx, limitsFactory.Settings)
		},
	}, nil
}

func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, blobs ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := newConsensusQueuePlugin(c.lggr, config, c.triggerCoordinator, c.triggerQueue, c.expirationDur, blobs)
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
	triggerQueue                        triggers.Queue
	blobBroadcastFetcher                ocr3_1types.BlobBroadcastFetcher
	expirationDur                       func(context.Context) (time.Duration, error)
}

func newConsensusQueuePlugin(
	lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig, triggerCoordinator triggers.Coordinator,
	triggerQueue triggers.Queue, expirationGraceDur func(context.Context) (time.Duration, error),
	blobs ocr3_1types.BlobBroadcastFetcher,
) *consensusQueuePlugin {
	c := &consensusQueuePlugin{
		lggr:                                logger.Named(lggr, "ConsensusQueuePlugin"),
		onchainCfg:                          onchainCfg,
		blobHandleMarshalledBytesUpperBound: ocr3_1types.BlobHandleMarshalledBytesUpperBound(onchainCfg.N, onchainCfg.F),
		triggerCoordinator:                  triggerCoordinator,
		triggerQueue:                        triggerQueue,
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
	// Always attempt to start any executions from the KV Store in the event Committed was not called. commitState is idempotent.
	queuedEvents, err := c.commitState(ctx, kvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to commit state: %w", err)
	}

	observedEvents, _, err := readObservedEvents(kvReader)
	if err != nil {
		return nil, err
	}

	capacity, err := c.triggerCoordinator.AvailableCapacity(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get available capacity: %w", err)
	}

	expirationGraceDur, err := c.expirationDur(ctx)
	if err != nil {
		c.lggr.Warnw("Failed to get expiration grace duration - using default", "err", err)
	}

	emptyBlobHandle := make([]byte, c.blobHandleMarshalledBytesUpperBound)
	obs := &pb.ConsensusQueueObservation{
		// placeholder for accurate size calculation
		BlobHandles:        slices.Repeat([][]byte{emptyBlobHandle}, maxObservedBlobHandles),
		Expiration:         timestamppb.New(time.Now().Add(-expirationGraceDur)),
		Capacity:           uint32(min(capacity.Global, math.MaxUint32)),
		WorkflowCapacities: make(map[string]uint32),
	}
	for _, e := range observedEvents {
		if wfCap, ok := capacity.PerWorkflow[e.WorkflowID]; ok {
			obs.WorkflowCapacities[e.WorkflowID] = uint32(min(wfCap, math.MaxUint32))
		}
	}
	for _, e := range queuedEvents {
		if wfCap, ok := capacity.PerWorkflow[e.WorkflowID]; ok {
			obs.WorkflowCapacities[e.WorkflowID] = uint32(min(wfCap, math.MaxUint32))
		}
	}

	protoOpts := proto.MarshalOptions{Deterministic: true}

	// Assemble an observation up to maxObservationBytes
	var coordEvents []*pb.CoordinatedEvent
	for coord := range toProtoSeq(c.triggerQueue.Sorted()) {
		observedAt, deadline := coord.ObservedAt, coord.Deadline
		coord.ObservedAt, coord.Deadline = nil, nil  // for consistent hashing
		protoBytes, err2 := protoOpts.Marshal(coord) // TODO opt: avoid remarshalling
		if err2 != nil {
			c.lggr.Errorw("Failed to marshal event. Excluding from observation", "eventID", coord.Event.Event.Id, "err", err2)
			//TODO remove from triggerQueue?
			continue
		}
		timelessHash := sha256.Sum256(protoBytes)

		// skip if kv observation exists already
		if obsEvent, ok := observedEvents[string(timelessHash[:])]; ok {
			if slices.Contains(obsEvent.NodeIds, uint8(c.onchainCfg.OracleID)) {
				continue
			}
		}

		obs.Events = append(obs.Events, &pb.ConsensusQueueEvent{
			WorkflowID:     coord.WorkflowID,
			TriggerEventID: coord.Event.Event.Id,
			ObservedAt:     observedAt,
			Deadline:       deadline,
			TimelessHash:   timelessHash[:],
			Blobs: []*pb.CoordinatedEventBlobHandle{{
				Blob:    0,
				Element: uint32(len(coordEvents)),
			}},
		})
		coordEvents = append(coordEvents, coord)

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

	obs.BlobHandles = nil // remove the placeholders
	if len(coordEvents) > 0 {
		blobHandle, n, err2 := broadcastAsBlob(ctx, blobBroadcaster, coordEvents, seqNr, c.onchainCfg.EstimatedRoundInterval)
		if err2 != nil {
			return nil, err2
		}
		obs.BlobHandles = append(obs.BlobHandles, blobHandle)
		// note: elements already reference blob 0 and natural element index is aligned

		// we verified space for up to maxObservedBlobHandles blobs
		var offset int
		for batch := 1; batch < maxObservedBlobHandles && n < len(coordEvents); batch++ {
			offset += n
			coordEvents = coordEvents[n:]
			blobHandle, n, err2 = broadcastAsBlob(ctx, blobBroadcaster, coordEvents, seqNr, c.onchainCfg.EstimatedRoundInterval)
			if err2 != nil {
				return nil, err2
			}
			obs.BlobHandles = append(obs.BlobHandles, blobHandle)
			for i := offset; i < n; i++ { // update indices
				obs.Events[offset+i].Blobs[0].Blob = uint32(batch)
				obs.Events[offset+i].Blobs[0].Element = uint32(i)
			}
		}
	}

	return protoOpts.Marshal(obs)
}

func toProtoSeq(in iter.Seq[triggers.CoordinatedEvent]) iter.Seq[*pb.CoordinatedEvent] {
	return func(yield func(*pb.CoordinatedEvent) bool) {
		for event := range in {
			if !yield(event.ToProto()) {
				return
			}
		}
	}
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
	if err = eg.Wait(); err != nil {
		return err
	}

	//TODO ensure capacity map is only observations and KV observed & queued events

	// ensure no duplicated trigger event ID or timelessHash
	triggerEventIDs, timelessHashes := make(map[string]struct{}), make(map[string]struct{})
	for _, event := range obs.Events {
		if _, ok := triggerEventIDs[event.TriggerEventID]; ok {
			return fmt.Errorf("duplicate trigger event ID: %s", event.TriggerEventID)
		}
		triggerEventIDs[event.TriggerEventID] = struct{}{}

		timelessHash := string(event.TimelessHash)
		if _, ok := timelessHashes[timelessHash]; ok {
			return fmt.Errorf("duplicate timeless hash: %x", event.TimelessHash)
		}
		timelessHashes[timelessHash] = struct{}{}
	}

	// wait to do this since we depend on every event being valid
	if !slices.IsSortedFunc(obs.Events, (*pb.ConsensusQueueEvent).CompareDeadline) {
		return fmt.Errorf("events must be sorted")
	}

	return nil
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
	done := make(chan struct{})
	var obs observations

	go func() {
		defer close(done)
		obs.aggregate(ctx, c.onchainCfg.F, aos)
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

	// update state
	state.pruneExpiredObservations(obs.medianExpiration)
	var queuedTriggerEventIDs []string // to be removed from the trigger.Queue
	state.addObservedEvents(c.onchainCfg.F, obs.newObservedEvents, obs.observedBlobs, func(event *pb.ConsensusQueueEvent) {
		queuedTriggerEventIDs = append(queuedTriggerEventIDs, event.TriggerEventID)
	})
	// pop as many queuedEvents as possible in to consensusEvents (up to capacity or KV limits)
	for len(state.queuedEvents) > 0 && len(state.consensusEvents) < int(obs.medianCap) {
		event := state.queuedEvents.popLimited(obs.medianCapPerWF, obs.medianExpiration)
		if event == nil {
			break
		}
		state.consensusEvents = append(state.consensusEvents, event)
		// stop when either value or the overall size is too large
		if events, blobs := state.sizeOfConsensusEvents(), state.sizeOfBlobHandles(); events > maxKeyValueValueBytes ||
			blobs > maxKeyValueValueBytes || events+blobs > maxKeyValueModifiedKeysPlusValuesBytes/2 {
			state.consensusEvents = state.consensusEvents[:len(state.consensusEvents)-1]
			break
		}
	}

	// prune observations if that value or the overall size is too large
	for len(state.observedEventEntries) > 0 &&
		(state.sizeOfObservedEvents() > maxKeyValueValueBytes ||
			state.sizeOf() > maxKeyValueModifiedKeysPlusValuesBytes) {
		// sorted by decreasing priority, so prune the last
		toRemove := state.observedEventEntries[len(state.observedEventEntries)-1]
		delete(state.observedEvents, toRemove.Event.TriggerEventID)
		//TODO opt: drop more than one each time?
	}

	// prune (newly) queuedEvents if that value, blob handles, or the overall size is too large
	if len(queuedTriggerEventIDs) > 0 &&
		(state.sizeOfQueuedEvents() > maxKeyValueValueBytes ||
			state.sizeOfBlobHandles() > maxKeyValueValueBytes ||
			state.sizeOf() > maxKeyValueModifiedKeysPlusValuesBytes) {
		// prune lowest priority queuedEvents from newly added queuedTriggerEventIDs
		queuedTriggerEventIDSet := map[string]struct{}{}
		for _, id := range queuedTriggerEventIDs {
			queuedTriggerEventIDSet[id] = struct{}{}
		}
		eligibleEvents := make([]*pb.ConsensusQueueEvent, 0, len(state.queuedEvents))
		for _, e := range state.queuedEvents {
			if _, ok := queuedTriggerEventIDSet[e.TriggerEventID]; ok {
				eligibleEvents = append(eligibleEvents, e)
			}
		}
		backwardsQueuedEvents := &revPriorityQueue{priorityQueue: new(priorityQueue(eligibleEvents))}
		heap.Init(backwardsQueuedEvents)

		// prune until we fit
		for backwardsQueuedEvents.Len() > 0 &&
			(state.sizeOfQueuedEvents() > maxKeyValueValueBytes ||
				state.sizeOfBlobHandles() > maxKeyValueValueBytes ||
				state.sizeOf() > maxKeyValueModifiedKeysPlusValuesBytes) {
			toRemove := heap.Pop(backwardsQueuedEvents).(*pb.ConsensusQueueEvent)
			//TODO opt: faster delete...but need to maintain state.queuedEvents (though not order)...map indices?
			state.queuedEvents = slices.DeleteFunc(state.queuedEvents, func(e *pb.ConsensusQueueEvent) bool { return e.TriggerEventID == toRemove.TriggerEventID })
			//TODO opt: drop more than one each time?
		}
		// rebuild the priority queue
		heap.Init(&state.queuedEvents)
	}

	// note: no need to prune consensus events, since already limited what we collected

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

	err error
}

func (res *observations) aggregate(ctx context.Context, f int, aos []ocrtypes.AttributedObservation) {
	res.medianCapPerWF = make(map[string]uint32)
	res.newObservedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)

	expirations := make([]*timestamppb.Timestamp, 0, len(aos))
	capacities := make([]uint32, 0, len(aos))
	wfCapacities := make(map[string][]uint32)
	for _, ao := range aos {
		var obs pb.ConsensusQueueObservation
		err := proto.Unmarshal(ao.Observation, &obs)
		if err != nil {
			res.err = fmt.Errorf("failed to unmarshal observation from oracle %d: %w", ao.Observer, err)
			return
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

	// If we used median, then honest nodes reporting reduced capacity would not be able to keep up.
	// Instead, we use the F+1th lowest value - the least honest value.
	var medianCount int
	res.medianCap = slices.Sorted(slices.Values(capacities))[f]
	for wf, caps := range wfCapacities {
		if len(caps) > 2*f+1 {
			res.medianCapPerWF[wf] = slices.Sorted(slices.Values(caps))[f]
			if medianCount < 2*f+1 {
				res.medianCapPerWF[wf] = (res.medianCapPerWF[wf] * 9) / 10
			}
		}
	}
	return
}

// addEvent adds the observed event to the observation set, either merging it into an existing event or creating a new one.
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
	_, err := c.commitState(ctx, kvReader)
	return err
}

// commitState executes the consensus events, and remove events from the trigger queue that have been queued or executed.
func (c *consensusQueuePlugin) commitState(ctx context.Context, kvReader ocr3_1types.KeyValueStateReader) ([]*pb.ConsensusQueueEvent, error) {
	consensusEvents, blobHandles, err := readConsensusState(kvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read consensus state: %w", err)
	}

	done := make(chan struct{})
	defer func() { <-done }()
	go c.tryStartExecutions(ctx, done, consensusEvents, blobHandles)

	c.triggerQueue.Remove(ctx, func(yield func(string) bool) {
		for _, e := range consensusEvents {
			if !yield(e.TriggerEventID) {
				break
			}
		}
	})
	queuedEvents, err := readQueuedEvents(kvReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read queued events: %w", err)
	}
	c.triggerQueue.Remove(ctx, func(yield func(string) bool) {
		for _, e := range queuedEvents {
			if !yield(e.TriggerEventID) {
				break
			}
		}
	})
	return queuedEvents, nil
}

// tryStartExecutions calls StartExecution asynchronously on the corresponding engine for each trigger event.
// It is idempotent, thus safe to call from Observation and Committed.
func (c *consensusQueuePlugin) tryStartExecutions(ctx context.Context, done chan struct{}, consensusEvents []*pb.ConsensusQueueEvent, blobHandles [][]byte) {
	defer close(done)
	fetcher := blobCache{fetcher: c.blobBroadcastFetcher, cache: make(map[string]func() ([]byte, error))}
	var wg sync.WaitGroup
	for _, consensusEvent := range consensusEvents {
		wg.Add(1)
		c.eng.GoCtx(ctx, func(ctx context.Context) {
			defer wg.Done()

			if err2 := c.tryStartExecution(ctx, &fetcher, consensusEvent, blobHandles); err2 != nil {
				c.lggr.Errorw("Failed to execute event", "eventID", consensusEvent.TriggerEventID, "err", err2)
				return
			}
		})
	}

	wg.Wait()
	return
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

	return c.triggerCoordinator.StartExecution(ctx, event)
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
