package queue

import (
	"container/heap"
	"context"
	"fmt"
	"math"
	"slices"
	"sync"

	"golang.org/x/exp/constraints"
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
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"

	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

var _ ocr3_1types.ReportingPluginFactory[[]byte] = &consensusQueuePluginFactory{}

type consensusQueuePluginFactory struct {
	lggr               logger.Logger
	cfg                *PluginConfig
	triggerCoordinator triggers.Coordinator
}

func NewConsensusQueuePluginFactory(lggr logger.Logger, cfg *PluginConfig, limitsFactory limits.Factory, triggerCoordinator triggers.Coordinator) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	// TODO inject CapacityReporter
	// TODO construct limiters
	return &consensusQueuePluginFactory{lggr: lggr, cfg: cfg, triggerCoordinator: triggerCoordinator}, nil
}

func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, blobs ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := newConsensusQueuePlugin(c.lggr, config, c.triggerCoordinator, blobs)
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
}

func newConsensusQueuePlugin(lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig, triggerCoordinator triggers.Coordinator, blobs ocr3_1types.BlobBroadcastFetcher) *consensusQueuePlugin {
	c := &consensusQueuePlugin{
		lggr:               logger.Named(lggr, "ConsensusQueuePlugin"),
		onchainCfg:         onchainCfg,
		triggerCoordinator: triggerCoordinator,
		blobs:              blobs,
	}
	c.Service, c.eng = services.Config{
		Name:  "ConsensusQueueReportingPlugin",
		Close: c.Close,
	}.NewServiceEngine(lggr)

	return c
}

func (c *consensusQueuePlugin) Query(_ context.Context, _ uint64, _ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher) (ocrtypes.Query, error) {
	return nil, nil
}

func (c *consensusQueuePlugin) Observation(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, kvState ocr3_1types.KeyValueStateReader, blobs ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Observation, error) {
	// Always attempt to start any executions from the KV Store in the event Committed was not called.
	// tryStartExecutions is idempotent.
	c.tryStartExecutions(kvState)

	capacity, err := c.triggerCoordinator.AvailableCapacity(ctx)
	if err != nil {
		//TODO or just log and make a best effort to observe events regardless?
		return nil, fmt.Errorf("failed to get available capacity: %w", err)
	}

	// Build Observation
	var events []triggers.CoordinatedEvent
	// TODO observe queue via c.CentralTriggerQueue.TakeForObservation()... one at at time?
	obs := &pb.ConsensusQueueObservation{
		Events:             make([]*pb.ConsensusQueueEvent, 0, len(events)),
		Capacity:           uint32(min(capacity.Global, math.MaxUint32)),
		WorkflowCapacities: make(map[string]uint32),
	}
	for wf, c := range capacity.PerWorkflow {
		//TODO only for events proposed and queued in the KV store?
		obs.WorkflowCapacities[wf] = uint32(min(c, math.MaxUint32))
	}

	const maxObservationSize = 0 //TODO
	protoOpts := proto.MarshalOptions{Deterministic: true}
	for _, e := range events {
		m := e.ToProto()
		// Try adding to batch

		//TODO broadcast now synchronously? or earlier so it's already available now?
		//TODO pb.CoordinatedEvent? or could omit wfID, eventID, deadline
		// handle, err := blobs.BroadcastBlob(ctx, []byte{}, ocr3_1types.BlobExpirationHintSequenceNumber{SeqNr: 0})
		obs.Events = append(obs.Events, &pb.ConsensusQueueEvent{
			WorkflowID:     m.WorkflowID,
			TriggerEventID: m.Event.Event.Id,
			Deadline:       m.Deadline,
			Blobs: []*pb.CoordinatedEventBlobHandle{{
				Handle: nil, //TODO
				Index:  0,   //TODO
			}},
		})
		id := e.Event.Event.ID
		var addedCap bool
		if _, ok := obs.WorkflowCapacities[e.WorkflowID]; !ok {
			addedCap = true
			obs.WorkflowCapacities[e.WorkflowID] = uint32(min(capacity.PerWorkflow[e.WorkflowID], math.MaxUint32))
		}
		// Check size
		if protoOpts.Size(obs) > maxObservationSize {
			// too large, remove this one
			obs.Events = obs.Events[:len(obs.Events)-1]
			if addedCap {
				delete(obs.WorkflowCapacities, id)
			}
			break
		}
	}
	return protoOpts.Marshal(obs)
}

func (c *consensusQueuePlugin) ValidateObservation(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, ao ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher,
) error {
	var obs pb.ConsensusQueueObservation
	err := proto.Unmarshal(ao.Observation, &obs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal observation: %w", err)
	}
	//TODO ensure each event has single blob+deadline
	//TODO ensure blobs exist?
	return nil
}

func (c *consensusQueuePlugin) ObservationQuorum(
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation,
	_ ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobFetcher,
) (quorumReached bool, err error) {
	return quorumhelper.ObservationCountReachesObservationQuorum(quorumhelper.QuorumTwoFPlusOne, c.onchainCfg.N, c.onchainCfg.F, aos), nil
}

func (c *consensusQueuePlugin) StateTransition(
	ctx context.Context, seqNr uint64, aq ocrtypes.AttributedQuery, aos []ocrtypes.AttributedObservation,
	kvState ocr3_1types.KeyValueStateReadWriter, fetcher ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	var wg sync.WaitGroup

	//TODO is median only logical option for capacity? what if honest nodes have less, don't we ideally want execute min(honest)?
	// consider some configurable x%
	var medianCap uint32
	var medianCapPerWF map[string]uint32 // wf ID -> capacity
	//TODO key on eventID+workflowID?
	var newObservedEvents map[string]*pb.ConsensusQueueObservedEvent // trigger event ID -> event
	go func() {
		defer wg.Done()

		capacities := make([]uint32, 0, len(aos))
		wfCapacities := make(map[string][]uint32)
		for _, ao := range aos {
			var obs pb.ConsensusQueueObservation
			err := proto.Unmarshal(ao.Observation, &obs)
			if err != nil {
				c.lggr.Errorw("Failed to unmarshal observation", "err", err)
				continue
			}
			capacities = append(capacities, obs.Capacity)
			for wf, c := range obs.WorkflowCapacities {
				wfCapacities[wf] = append(wfCapacities[wf], c)
			}
			for _, event := range obs.Events {
				//TODO key=WFID+eventID?
				mergedEvent, ok := newObservedEvents[event.TriggerEventID]
				if !ok {
					mergedEvent = &pb.ConsensusQueueObservedEvent{
						Deadlines: []*timestamppb.Timestamp{event.Deadline},
						NodeIds:   []uint8{uint8(ao.Observer)},
						Blobs:     []*pb.CoordinatedEventBlobHandle{event.Blobs[0]},
					}
				} else {
					mergedEvent.Deadlines = append(mergedEvent.Deadlines, event.Deadline)
					mergedEvent.Blobs = append(mergedEvent.Blobs, event.Blobs...)
				}

				newObservedEvents[event.TriggerEventID] = mergedEvent
			}
		}
		medianCap = median(capacities)
		for wf, c := range wfCapacities {
			medianCapPerWF[wf] = median(c)
		}
	}()

	// Deserialize the KV store into memory
	//TODO package in a constructor for observedAndQueued?
	state := newStateReadWriter(kvState)
	observedEvents, err := state.readObservedEvents()
	if err != nil {
		return nil, err
	}
	queuedEvents, err := state.readPriorityQueue()
	if err != nil {
		return nil, err
	}
	heap.Init(&queuedEvents) //TODO or can we trust it?

	//TODO prune expired events from proposed and priorityQueue

	wg.Wait() // TODO also wait if we return early for error

	// Collect observed proposals and promote to priority queue on f+1
	observedAndQueued{
		observedEvents: observedEvents,
		queuedEvents:   &queuedEvents,
	}.addObservedEvents(c.onchainCfg.F, newObservedEvents)

	// Select items from queue based on priority and capacity
	consensusEvents, err := queuedEvents.PopCap(medianCap, medianCapPerWF)
	if err != nil {
		return nil, err
	}

	// Updated KV Store
	//TODO might we hit a KV size limit in the middle of these? then what? can't just drop proposals
	// alternative: proposals only valid when there is definitely enough space
	err = state.writeObservedEvents(observedEvents)
	if err != nil {
		return nil, err
	}
	err = state.writePriorityQueue(queuedEvents)
	if err != nil {
		return nil, err
	}
	err = state.writeConsensusEvents(&pb.ConsensusQueueEventList{Events: consensusEvents})
	if err != nil {
		return nil, err
	}

	// no report to return, executions start when committed
	return ocr3_1types.ReportsPlusPrecursor{}, nil
}

func (c *consensusQueuePlugin) Committed(ctx context.Context, seqNr uint64, keyValueStateReader ocr3_1types.KeyValueStateReader) error {
	c.tryStartExecutions(keyValueStateReader)
	return nil
}

// tryStartExecutions calls StartExecution asynchronously on the corresponding engine for each trigger event.
// It is idempotent, thus safe to call from Observation and Committed.
func (c *consensusQueuePlugin) tryStartExecutions(reader ocr3_1types.KeyValueStateReader) {
	state := stateReader{reader}
	pbs, err := state.readConsensusEvents()
	if err != nil {
		c.lggr.Errorw("Failed to read consensus events", "err", err)
		return
	}

	var wg sync.WaitGroup
	for _, pbEvent := range pbs.Events {
		var event triggers.CoordinatedEvent
		wg.Add(1)
		c.eng.Go(func(ctx context.Context) {
			defer wg.Done()

			// fetch blobs with observed events
			errGrp := new(errgroup.Group)
			observedEvents := make([]*pb.CoordinatedEvent, len(pbEvent.Blobs))
			for i, blob := range pbEvent.Blobs {
				errGrp.Go(func() error {
					var handle ocr3_1types.BlobHandle
					if err := handle.UnmarshalBinary(blob.Handle); err != nil {
						return fmt.Errorf("failed to unmarshal blob handle: %w", err)
					}
					b, err := c.blobs.FetchBlob(ctx, handle)
					if err != nil {
						return fmt.Errorf("failed to fetch blob: %w", err)
					}
					//TODO allocate element first?
					if err := proto.Unmarshal(b, observedEvents[i]); err != nil {
						return fmt.Errorf("failed to unmarshal CoordinatedEvent from blob: %w", err)
					}
					return nil
				})
			}
			if err := errGrp.Wait(); err != nil {
				c.lggr.Errorw("Failed to fetch blob for event; unable to execute", "eventID", pbEvent.TriggerEventID, "err", err)
				return
			}

			//TODO we don't always want a continuous median, e.g. IDs
			// maybe mode if it exists, otherwise fallback to each node using its own observed value?
			if err := event.FromProto(&pb.CoordinatedEvent{
				WorkflowID:   pbEvent.WorkflowID,
				TriggerCapID: "",  //TODO from blob median?
				TriggerIndex: 0,   //TODO from blob median?
				ObservedAt:   nil, //TODO from blob median?
				Deadline:     pbEvent.Deadline,
				Event:        nil, //TODO from blob ....piecewise median?
			}); err != nil {
				c.lggr.Errorw("Failed to unmarshal CoordinatedEvent", "pb", pbEvent.String(), "err", err)
			}

			// TODO execute StartExecution
			// c.triggerCoordinator.tryStartExecution(ctx, event)
		})
	}
	wg.Wait() //TODO do we need to wait?
}

func (c *consensusQueuePlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, nil
}

func (c *consensusQueuePlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}

func (c *consensusQueuePlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return false, nil
}

func median[N constraints.Integer, S []N](s S) N {
	slices.Sort(s)
	return (s[len(s)/2-1] + s[len(s)/2]) / 2
}
