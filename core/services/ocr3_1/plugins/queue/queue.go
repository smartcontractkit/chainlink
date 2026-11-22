// Package queue contains the consensus queue OCR 3.1 reporting plugin, for coordinating execution of triggered events.
// Each node has a local [triggers.Queue] of [triggers.CoordinatedEvent], which are observed (and broadcast) by the plugin.
// When consensus is reached, events are put in to the KeyValueState priority queue (and removed from the [triggers.Queue]).
// The priority queue is ordered by deadline (with nanosecond precision) and then by trigger event ID.
// Each round, as many events as possible are popped from the priority queue and then (fetched and) executed via [triggers.Coordinator.Deliver].
// Capacity, both overall and per-workflow, limits what may execute. Capacities are observed and medianized.
package queue

import (
	"bytes"
	"cmp"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"maps"
	"math"
	"slices"
	"sync"
	"time"

	"github.com/smartcontractkit/libocr/commontypes"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3types"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/libocr/quorumhelper"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
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
		expirationDur:      triggers.ConsensusQueueExpirationGrace(limitsFactory),
	}, nil
}

func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, blobs ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := newConsensusQueuePlugin(c.lggr, config, c.cfg.Limits, c.triggerCoordinator, c.triggerQueue, c.expirationDur, blobs)
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
	limits                              ocr3_1types.ReportingPluginLimits
	blobHandleMarshalledBytesUpperBound int
	triggerCoordinator                  triggers.Coordinator
	triggerQueue                        triggers.Queue
	blobBroadcastFetcher                ocr3_1types.BlobBroadcastFetcher
	expirationDur                       func(context.Context) (time.Duration, error)
}

func newConsensusQueuePlugin(
	lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig, limits ocr3_1types.ReportingPluginLimits, triggerCoordinator triggers.Coordinator,
	triggerQueue triggers.Queue, expirationGraceDur func(context.Context) (time.Duration, error),
	blobs ocr3_1types.BlobBroadcastFetcher,
) *consensusQueuePlugin {
	c := &consensusQueuePlugin{
		lggr:                                logger.Named(lggr, "ConsensusQueuePlugin"),
		onchainCfg:                          onchainCfg,
		limits:                              limits,
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

	obs := &pb.ConsensusQueueObservation{
		// placeholder for accurate size calculation
		BlobHandles:        slices.Repeat([][]byte{make([]byte, c.blobHandleMarshalledBytesUpperBound)}, maxObservedBlobHandles),
		Expiration:         timestamppb.New(time.Now().Add(-expirationGraceDur)),
		Capacity:           clampUint32(capacity.Global),
		WorkflowCapacities: stateWorkflowCapacities(capacity, observedEvents, queuedEvents),
	}
	kvWFCapacities := maps.Clone(obs.WorkflowCapacities) // in case we need to rebuild later

	events, err := c.assembleEvents(ctx, obs, observedEvents, capacity.PerWorkflow)
	if err != nil {
		return nil, err
	}
	obs.BlobHandles = nil // remove placeholders

	fit, err := c.broadcastObservationBlobs(ctx, obs, blobBroadcaster, events, seqNr)
	if err != nil {
		return nil, err
	}
	if !fit { // events did not all fit, so rebuild the wf capacities
		obs.WorkflowCapacities = kvWFCapacities
		for _, event := range obs.Events {
			if _, ok := obs.WorkflowCapacities[event.WorkflowID]; !ok {
				obs.WorkflowCapacities[event.WorkflowID] = clampUint32(capacity.PerWorkflow[event.WorkflowID])
			}
		}
	}

	return proto.MarshalOptions{Deterministic: true}.Marshal(obs)
}

func clampUint32(v int) uint32 { return uint32(max(0, min(v, math.MaxUint32))) }

// stateWorkflowCapacities returns the capacity of every workflow that has an observed or queued event in the KV state.
func stateWorkflowCapacities(
	capacity triggers.Capacity, observedEvents map[string]*pb.ConsensusQueueObservedEvent, queuedEvents []*pb.ConsensusQueueEvent,
) map[string]uint32 {
	capacities := make(map[string]uint32)
	add := func(workflowID string) {
		if wfCap, ok := capacity.PerWorkflow[workflowID]; ok {
			capacities[workflowID] = clampUint32(wfCap)
		}
	}
	for _, e := range observedEvents {
		add(e.WorkflowID)
	}
	for _, e := range queuedEvents {
		add(e.WorkflowID)
	}
	return capacities
}

// assembleEvents adds the local queue's events to obs, in order, up to the observation size limit. It returns the
// events (without timestamps, as they will appear in blobs) to be broadcast.
func (c *consensusQueuePlugin) assembleEvents(
	ctx context.Context, obs *pb.ConsensusQueueObservation, observedEvents map[string]*pb.ConsensusQueueObservedEvent,
	wfCapacities map[string]int,
) ([]*pb.CoordinatedEvent, error) {
	var events []*pb.CoordinatedEvent
	for timelessHash, event := range c.triggerQueue.Sorted(ctx, maxObservedEvents, c.maxEventBytes()) {
		// skip if kv observation exists already
		if obsEvent, ok := observedEvents[string(timelessHash)]; ok {
			if slices.Contains(obsEvent.NodeIds, uint8(c.onchainCfg.OracleID)) {
				continue
			}
		}

		events = append(events, event)
		newWF, err := observeEvent(obs, timelessHash, event, wfCapacities)
		if err != nil {
			return nil, err
		}

		if observationSize(obs) > c.limits.MaxObservationBytes {
			// too large, remove this one
			events = events[:len(events)-1]
			obs.Events = obs.Events[:len(obs.Events)-1]
			if newWF {
				delete(obs.WorkflowCapacities, event.WorkflowID)
			}
			break
		}
	}
	return events, nil
}

// observationSize returns the encoded size of an observation that grows as events are appended, without sizing every
// event again each time (which would make assembling an observation quadratic).
//
// It uses proto.MarshalOptions.UseCachedSize, which reuses the size a previous Size call stored in each message. That
// is only correct for a message that has not changed since it was sized, so it is not applied to obs itself, which does
// change. Instead it sizes a shallow copy, a new message whose own size was never cached, and which shares obs's fields
// (including any added to the message in the future). The events are unchanged since they were added, so their sizes
// are valid. An event not sized yet has no cached size, and is computed as usual.
//
// Callers must not modify events, or any message inside them, between calls. (Marshaling obs, which does not use
// cached sizes, recomputes and refreshes them, so modifying events afterwards is fine.)
func observationSize(obs *pb.ConsensusQueueObservation) int {
	m := obs.ProtoReflect()
	probe := m.New()
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		probe.Set(fd, v)
		return true
	})
	return proto.MarshalOptions{Deterministic: true, UseCachedSize: true}.Size(probe.Interface())
}

// observeEvent adds the event to the observation and updates the workflow capacities if required.
func observeEvent(
	obs *pb.ConsensusQueueObservation, timelessHash []byte, event *pb.CoordinatedEvent, wfCapacities map[string]int,
) (newWF bool, err error) {
	observedAt, deadline := event.ObservedAt, event.Deadline
	event.ObservedAt, event.Deadline = nil, nil // for consistent hashing

	executionID, err := workflows.GenerateExecutionIDWithTriggerIndex(event.WorkflowID, event.Event.Event.Id, int(event.TriggerIndex))
	if err != nil {
		return false, fmt.Errorf("failed to generate execution ID: %w", err)
	}

	idx := len(obs.Events)
	obs.Events = append(obs.Events, &pb.ConsensusQueueEvent{
		WorkflowID:   event.WorkflowID,
		ExecutionID:  executionID,
		ObservedAt:   observedAt,
		Deadline:     deadline,
		TimelessHash: timelessHash,
		Blobs: []*pb.CoordinatedEventBlobHandle{{
			Blob:    0,
			Element: uint32(idx),
		}},
	})

	if _, ok := obs.WorkflowCapacities[event.WorkflowID]; !ok {
		newWF = true
		obs.WorkflowCapacities[event.WorkflowID] = uint32(min(wfCapacities[event.WorkflowID], math.MaxUint32))
	}

	return
}

// broadcastObservationBlobs broadcasts up to maxObservedBlobHandles blobs containing events, and updates the observation.
// It may not be able to broadcast all events, in which case the observed events will be truncated and false returned.
// The events are divided between blobs first, so that the blobs can be broadcast concurrently.
func (c *consensusQueuePlugin) broadcastObservationBlobs(
	ctx context.Context, obs *pb.ConsensusQueueObservation, blobBroadcaster ocr3_1types.BlobBroadcaster,
	events []*pb.CoordinatedEvent, seqNr uint64,
) (bool, error) {
	if len(events) == 0 {
		return true, nil
	}
	if len(events) != len(obs.Events) {
		return false, fmt.Errorf("number of events and observedEvents must match")
	}

	// ends[i] is the index after the last event of blob i. We verified space for maxObservedBlobHandles handles.
	var ends []int
	for start := 0; start < len(events) && len(ends) < maxObservedBlobHandles; {
		n := c.blobPrefix(events[start:])
		if n == 0 {
			break
		}
		start += n
		ends = append(ends, start)
	}
	if len(ends) == 0 {
		c.lggr.Warnw("No events fit in a blob - omitting events from observation", "events", len(events), "blobLimit", c.limits.MaxBlobPayloadBytes)
		obs.Events = obs.Events[:0]
		return false, nil
	}

	handles := make([][]byte, len(ends))
	errs := make([]error, len(ends))
	// At most maxObservedBlobHandles (10, the vault plugin's limit on concurrent broadcasts) are broadcast at once.
	var eg errgroup.Group
	for i, start := 0, 0; i < len(ends); i++ {
		batch := events[start:ends[i]]
		start = ends[i]
		eg.Go(func() error {
			handles[i], errs[i] = c.broadcastBlob(ctx, blobBroadcaster, batch, seqNr)
			return nil
		})
	}
	_ = eg.Wait()
	if err := ctx.Err(); err != nil {
		return false, err // cancelled, so there is no point in observing at all
	}

	// Keep the blobs before the first failure. A failed broadcast must not fail Observation (see libocr
	// BlobBroadcaster), so observe without the events that did not make it.
	keep := len(ends)
	for i, err := range errs {
		if err != nil {
			c.lggr.Warnw("Failed to broadcast events blob - omitting it and the remaining events from observation", "blob", i, "err", err)
			keep = i
			break
		}
	}
	if keep == 0 {
		obs.Events = obs.Events[:0]
		return false, nil
	}

	obs.BlobHandles = append(obs.BlobHandles, handles[:keep]...)
	for i, start := 0, 0; i < keep; i++ { // point every event at its blob and element
		for j, e := range obs.Events[start:ends[i]] {
			e.Blobs[0].Blob, e.Blobs[0].Element = uint32(i), uint32(j)
		}
		start = ends[i]
	}
	limited := ends[keep-1] < len(obs.Events)
	if limited {
		obs.Events = obs.Events[:ends[keep-1]] // drop the events that did not fit in the blobs, or whose blob failed
		// note: caller will fix wfCapacity
	}
	return !limited, nil
}

func (c *consensusQueuePlugin) ValidateObservation(
	ctx context.Context, _ uint64, _ ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation,
	kvReader ocr3_1types.KeyValueStateReader, blobFetcher ocr3_1types.BlobFetcher,
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
	if err = obs.Expiration.CheckValid(); err != nil {
		return fmt.Errorf("invalid expiration: %w", err)
	}
	for i, handle := range obs.BlobHandles {
		if len(handle) > c.blobHandleMarshalledBytesUpperBound {
			return fmt.Errorf("blob handle %d is %d bytes, more than the maximum %d", i, len(handle), c.blobHandleMarshalledBytesUpperBound)
		}
	}

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

	// ensure every event has an entry in the capacity map
	for _, event := range obs.Events {
		if _, ok := obs.WorkflowCapacities[event.WorkflowID]; !ok {
			return fmt.Errorf("missing capacity for workflow ID %s", event.WorkflowID)
		}
	}

	// ensure the capacity map only covers workflows an honest node reports: those with events in this observation,
	// or with observed or queued events in the KV state
	allowed := make(map[string]struct{})
	for _, event := range obs.Events {
		allowed[event.WorkflowID] = struct{}{}
	}
	observedEvents, _, err := readObservedEvents(kvReader)
	if err != nil {
		return fmt.Errorf("failed to read observed events: %w", err)
	}
	for _, event := range observedEvents {
		allowed[event.WorkflowID] = struct{}{}
	}
	queuedEvents, err := readQueuedEvents(kvReader)
	if err != nil {
		return fmt.Errorf("failed to read queued events: %w", err)
	}
	for _, event := range queuedEvents {
		allowed[event.WorkflowID] = struct{}{}
	}
	for wf := range obs.WorkflowCapacities {
		if _, ok := allowed[wf]; !ok {
			return fmt.Errorf("unexpected capacity for workflow ID %s", wf)
		}
	}

	// ensure no duplicate execution ID or timelessHash
	execIDs, timelessHashes := make(map[string]struct{}), make(map[string]struct{})
	for _, event := range obs.Events {
		if _, ok := execIDs[event.ExecutionID]; ok {
			return fmt.Errorf("duplicate workflow execution ID: %s", event.ExecutionID)
		}
		execIDs[event.ExecutionID] = struct{}{}

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
	} // TODO ensure not too long
	if event.ExecutionID == "" {
		return errors.New("event execution id is missing") // (the exact value is checked against the blob below)
	}
	if event.ObservedAt == nil {
		return errors.New("event observation ts is missing")
	} else if err := event.ObservedAt.CheckValid(); err != nil {
		return fmt.Errorf("invalid event observation ts: %w", err)
	}
	if event.Deadline == nil {
		return errors.New("event deadline is missing")
	} else if err := event.Deadline.CheckValid(); err != nil {
		return fmt.Errorf("invalid event deadline: %w", err)
	}
	if len(event.TimelessHash) == 0 {
		return fmt.Errorf("event blob hash is missing")
	} else if len(event.TimelessHash) != sha256.Size {
		return fmt.Errorf("event blob hash must be %d bytes but got %d", sha256.Size, len(event.TimelessHash))
	}

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
	} else if coordEvent == nil {
		return fmt.Errorf("failed to fetch event from blob: event is nil")
	}

	if coordEvent.WorkflowID != event.WorkflowID {
		return fmt.Errorf("blob workflow ID %s does not match event workflow ID %s", coordEvent.WorkflowID, event.WorkflowID)
	}
	executionID, err := workflows.GenerateExecutionIDWithTriggerIndex(event.WorkflowID, coordEvent.GetEvent().GetEvent().GetId(), int(coordEvent.TriggerIndex))
	if err != nil {
		return fmt.Errorf("failed to generate execution ID: %w", err)
	}
	if executionID != event.ExecutionID {
		return fmt.Errorf("blob derived execution ID %s does not match event execution ID %s", executionID, event.ExecutionID)
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
	// Aggregating the observations (CPU) and loading the state (KV) are independent, so they overlap: this takes
	// about 9ms off a 43ms full round at N=31. Wait returns only when both are done, so neither outlives the call.
	var (
		obs   observations
		state kvState
		eg    errgroup.Group
	)
	eg.Go(func() error {
		obs.aggregate(ctx, c.onchainCfg.F, aos) // reports cancellation as its error, so a partial result is never used
		if obs.err != nil {
			return fmt.Errorf("failed to aggregate observations: %w", obs.err)
		}
		return nil
	})
	eg.Go(func() error {
		if err := state.loadFrom(kvReaderWriter); err != nil {
			return fmt.Errorf("failed to read state: %w", err)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// update state
	state.pruneExpiredObservations(obs.medianExpiration)
	var queuedTimelessHashes []string // eligible for pruning later if necessary
	state.addObservedEvents(c.onchainCfg.F, obs.newObservedEvents, obs.observedBlobs, func(event *pb.ConsensusQueueEvent) {
		queuedTimelessHashes = append(queuedTimelessHashes, string(event.TimelessHash))
	})
	state.popConsensusEvents(int(obs.medianCap), obs.medianCapPerWF, obs.medianExpiration, c.limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	state.pruneObservedToFit(c.limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	state.pruneQueuedToFit(queuedTimelessHashes, c.limits.MaxKeyValueModifiedKeysPlusValuesBytes)
	// note: no need to prune consensus events, since already limited what we collected

	if err := state.writeTo(&stateWriter{kv: kvReaderWriter}); err != nil {
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
	res.observedBlobs = make(blobHandleMap)

	// normalize the order, so that the outcome never depends on the order observations arrive in
	if len(aos) <= f { // quorum guarantees more, but never index out of range
		res.err = fmt.Errorf("need more than F=%d observations but got %d", f, len(aos))
		return
	}
	aos = slices.SortedFunc(slices.Values(aos), func(a, b ocrtypes.AttributedObservation) int { return cmp.Compare(a.Observer, b.Observer) })

	expirations := make([]*timestamppb.Timestamp, 0, len(aos))
	capacities := make([]uint32, 0, len(aos))
	wfCapacities := make(map[string][]uint32)
	for _, ao := range aos {
		if err := ctx.Err(); err != nil { // a cancelled aggregation is partial, so it must never look successful
			res.err = err
			return
		}
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
			if err2 := ctx.Err(); err2 != nil {
				res.err = err2
				return
			}
			res.addEvent(ao.Observer, event, obs.BlobHandles)
		}
	}
	res.medianExpiration = medianTimestamp(expirations).AsTime()

	// If we used median, then honest nodes reporting reduced capacity would not be able to keep up.
	// Instead, we use the F+1th lowest value - the least honest value.
	res.medianCap = slices.Sorted(slices.Values(capacities))[f]
	for wf, caps := range wfCapacities {
		// A node that did not report a workflow counts as zero capacity. Otherwise up to F nodes (faulty, or just
		// without the workflow) could omit it and leave too few reports to ever run it.
		caps = append(caps, make([]uint32, len(aos)-len(caps))...)
		res.medianCapPerWF[wf] = slices.Sorted(slices.Values(caps))[f]
	}
	return
}

// addEvent adds the observed event to the observation set, either merging it into an existing event or creating a new one.
func (res *observations) addEvent(obs commontypes.OracleID, event *pb.ConsensusQueueEvent, blobHandles [][]byte) {
	res.observedBlobs.add(string(event.TimelessHash), string(blobHandles[event.Blobs[0].Blob]), event.Blobs[0].Element)
	aggregated, ok := res.newObservedEvents[string(event.TimelessHash)]
	if !ok {
		res.newObservedEvents[string(event.TimelessHash)] = &pb.ConsensusQueueObservedEvent{
			WorkflowID:  event.WorkflowID,
			ExecutionID: event.ExecutionID,
			NodeIds:     []uint8{uint8(obs)},
			ObservedAts: []*timestamppb.Timestamp{event.ObservedAt},
			Deadlines:   []*timestamppb.Timestamp{event.Deadline},
			Blobs:       nil, // res.observedBlobs tracks these while in memory
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
// Starting an execution is best effort: once an event is in consensus, it leaves the local queue whether or not starting
// it succeeds, and a failure is only logged. (Committed and the next Observation each try, as consensus events stay in
// KV until the next StateTransition, but the queue never holds an event to retry it.)
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
			if !yield(e.ExecutionID) {
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
			if !yield(e.ExecutionID) {
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
				c.lggr.Errorw("Failed to execute event", "executionID", consensusEvent.ExecutionID, "err", err2)
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
		if int(blob.Blob) >= len(blobHandles) {
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

	//TODO ensure this only *starts* the execution, and does not block while it runs
	// TestCommitted_DoesNotWaitForExecutions
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
