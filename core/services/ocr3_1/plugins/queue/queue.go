package queue

import (
	"container/heap"
	"context"
	"fmt"
	"slices"
	"sync"

	"google.golang.org/protobuf/proto"

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
	lggr logger.Logger
	cfg  *PluginConfig
}

func NewConsensusQueuePluginFactory(lggr logger.Logger, cfg *PluginConfig, limitsFactory limits.Factory) (ocr3_1types.ReportingPluginFactory[[]byte], error) {
	// TODO pass CentralTriggerQueue, CapacityReporter, EngineRegistry
	// TODO construct limiters
	return &consensusQueuePluginFactory{lggr: lggr, cfg: cfg}, nil
}

func (c *consensusQueuePluginFactory) NewReportingPlugin(ctx context.Context, config ocr3types.ReportingPluginConfig, fetcher ocr3_1types.BlobBroadcastFetcher) (ocr3_1types.ReportingPlugin[[]byte], ocr3_1types.ReportingPluginInfo, error) {
	plugin := newConsensusQueuePlugin(c.lggr, config)
	pluginInfo := ocr3_1types.ReportingPluginInfo1{
		Name:   "ConsensusQueueReportingPlugin",
		Limits: c.cfg.Limits,
	}
	return plugin, pluginInfo, plugin.Start(ctx)
}

var _ ocr3_1types.ReportingPlugin[[]byte] = &consensusQueuePlugin{}

type consensusQueuePlugin struct {
	services.Service
	eng        *services.Engine
	lggr       logger.Logger
	onchainCfg ocr3types.ReportingPluginConfig
}

func newConsensusQueuePlugin(lggr logger.Logger, onchainCfg ocr3types.ReportingPluginConfig) *consensusQueuePlugin {
	c := &consensusQueuePlugin{lggr: logger.Named(lggr, "ConsensusQueuePlugin"), onchainCfg: onchainCfg}
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
	ctx context.Context, seqNr uint64, _ ocrtypes.AttributedQuery, kvState ocr3_1types.KeyValueStateReader, _ ocr3_1types.BlobBroadcastFetcher,
) (ocrtypes.Observation, error) {
	// Always attempt to start any executions from the KV Store in the event Committed was not called.
	// tryStartExecutions is idempotent.
	c.tryStartExecutions(kvState)

	// TODO Fetch limits from node

	// Build Observation
	var events []triggers.CoordinatedEvent
	// TODO observe queue via c.CentralTriggerQueue.TakeForObservation()... one at at time?
	obs := &pb.ConsensusQueueObservation{
		Events:             make([]*pb.CoordinatedEvent, 0, len(events)),
		Capacity:           0,   //TODO set
		WorkflowCapacities: nil, //TODO for everything proposed and queued in the KV store?
	}

	const maxObservationSize = 0 //TODO
	protoOpts := proto.MarshalOptions{Deterministic: true}
	for _, e := range events {
		m := e.ToProto()
		obs.Events = append(obs.Events, m)
		id := e.Event.Event.ID
		var addedCap bool
		if _, ok := obs.WorkflowCapacities[id]; !ok {
			addedCap = true
			//TODO set it
		}
		if protoOpts.Size(obs) > maxObservationSize { // too large, end of batch
			// remove this one
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
	//TODO validate
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
	kvState ocr3_1types.KeyValueStateReadWriter, _ ocr3_1types.BlobFetcher,
) (ocr3_1types.ReportsPlusPrecursor, error) {
	var wg sync.WaitGroup

	// TODO In parallel:
	var capacity int
	var wfCapacity map[string]int
	var observedProposals map[string][]uint8           // trigger event ID -> observer node IDs
	var observedEvents map[string]*pb.CoordinatedEvent // trigger event ID -> event
	go func() {
		defer wg.Done()

		// - Aggregate each node’s view of how many trigger events we can pop from the KV store queue (i.e. how many executions we can start)
		// - Aggregate the per workflow free capacity
		// - Achieve quorum on observations -> any item with F+1 observations goes into the queue
		// - Medianize the logical clock value in the observations to order events
	}()

	// Deserialize the KV store into memory
	state := newStateReadWriter(kvState)
	proposedEvents, err := state.readProposedEvents()
	if err != nil {
		return nil, err
	}
	priorityQueue, err := state.readPriorityQueue()
	if err != nil {
		return nil, err
	}
	heap.Init(&priorityQueue) //TODO or can we trust it?

	//TODO prune expired events from proposed and priorityQueue

	wg.Wait()

	// Collect observed proposals and promote to priority queue on f+1
	for id, nodes := range observedProposals {
		if obs, ok := proposedEvents.Observers[id]; !ok && len(nodes) >= c.onchainCfg.F+1 {
			// new proposal straight to priority queue
			heap.Push(&priorityQueue, proposedEvents.Events[id])
			continue
		} else if !ok {
			// new proposal < f + 1
			proposedEvents.Observers[id] = &pb.NodeIDList{Ids: nodes}
			proposedEvents.Events[id] = observedEvents[id]
			continue
		} else if len(nodes) < c.onchainCfg.F+1 {
			// update proposal
			// TODO opt: build a map to avoid n^2? or nbd since n < f+1
			for _, node := range nodes {
				if !slices.Contains(obs.Ids, node) {
					obs.Ids = append(obs.Ids, node)
				}
			}
			if len(obs.Ids) < c.onchainCfg.F+1 {
				continue
			}
		}
		// f+1 - promote proposal to priority queue
		delete(proposedEvents.Observers, id)
		delete(proposedEvents.Events, id)
		heap.Push(&priorityQueue, proposedEvents.Events[id])
	}

	// Select items from queue based on priority and capacity
	var consensusEvents []*pb.CoordinatedEvent
	var workflowLimited []*pb.CoordinatedEvent // omitted due to workflow limit
	for len(consensusEvents) < capacity {
		if len(priorityQueue) == 0 {
			break
		}
		next, ok := heap.Pop(&priorityQueue).(*pb.CoordinatedEvent)
		if !ok {
			return nil, fmt.Errorf("failed to cast *pb.CoordinatedEvent: got type %T", next)
		}
		if wfCapacity[next.WorkflowID] == 0 {
			workflowLimited = append(workflowLimited, next)
			continue
		}
		wfCapacity[next.WorkflowID]--
		consensusEvents = append(consensusEvents, next)
	}
	// restore omitted for next time
	for _, event := range workflowLimited {
		heap.Push(&priorityQueue, event)
	}

	// Updated KV Store
	err = state.writeProposedEvents(proposedEvents)
	if err != nil {
		return nil, err
	}
	err = state.writePriorityQueue(priorityQueue)
	if err != nil {
		return nil, err
	}
	err = state.writeConsensusEvents(&pb.CoordinatedEventList{Events: consensusEvents})
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
		//TODO log?
		return
	}
	errs := make([]error, len(pbs.Events))
	var wg sync.WaitGroup
	for i, pbEvent := range pbs.Events {
		var event triggers.CoordinatedEvent
		if err := event.FromProto(pbEvent); err != nil {
			errs[i] = fmt.Errorf("failed to unmarshal event: %s: %w", pbEvent, err)
			continue
		}
		wg.Add(1)
		c.eng.Go(func(ctx context.Context) {
			defer wg.Done()
			// TODO fetch workflow engine instance
			// TODO execute StartExecution
			errs[i] = nil
		})
	}
	wg.Wait()
	//TODO process errors? or just log as we go?
}

func (c *consensusQueuePlugin) Reports(ctx context.Context, seqNr uint64, reportsPlusPrecursor ocr3_1types.ReportsPlusPrecursor) ([]ocr3types.ReportPlus[[]byte], error) {
	return nil, nil
}

func (c *consensusQueuePlugin) ShouldAcceptAttestedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}

func (c *consensusQueuePlugin) ShouldTransmitAcceptedReport(ctx context.Context, seqNr uint64, reportWithInfo ocr3types.ReportWithInfo[[]byte]) (bool, error) {
	return true, nil
}
