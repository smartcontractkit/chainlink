package queue

import (
	"cmp"
	"container/heap"
	"fmt"
	"slices"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// KV store Keys
const (
	proposedEventsKey  = "proposed_events::index"  // proposed events waiting for consensus (<f+1 observers)
	priorityQueueKey   = "priority_queue::index"   // priority queue of events waiting for capacity
	consensusEventsKey = "consensus_events::index" // events the plugin committed to sending to nodes
)

func (s stateReader) readProposedEvents() (*pb.ConsensusQueueProposals, error) {
	return s.readValue[*pb.ConsensusQueueProposals](proposedEventsKey)
}

func (s stateReadWriter) writeProposedEvents(events *pb.ConsensusQueueProposals) error {
	return s.writeValue(proposedEventsKey, events)
}

func (s stateReader) readPriorityQueue() (priorityQueue, error) {
	list, err := s.readValue[*pb.CoordinatedEventList](priorityQueueKey)
	if err != nil {
		return nil, err
	}
	return list.Events, nil
}

func (s stateReadWriter) writePriorityQueue(queue priorityQueue) error {
	return s.writeValue(priorityQueueKey, &pb.CoordinatedEventList{Events: queue})
}

func (s stateReader) readConsensusEvents() (*pb.CoordinatedEventList, error) {
	return s.readValue[*pb.CoordinatedEventList](consensusEventsKey)
}

func (s stateReadWriter) writeConsensusEvents(events *pb.CoordinatedEventList) error {
	return s.writeValue(consensusEventsKey, events)
}

type stateReader struct {
	kv ocr3_1types.KeyValueStateReader
}

// readValue unmarshals the given proto message if one exists for the given key, otherwise it returns a zero value.
func (s stateReader) readValue[M proto.Message](key string) (m M, err error) {
	val, err := s.kv.Read([]byte(key))
	if err != nil || len(val) == 0 {
		return
	}
	err = proto.Unmarshal(val, m)
	return
}

type stateReadWriter struct {
	stateReader
	kv ocr3_1types.KeyValueStateReadWriter
}

func newStateReadWriter(kv ocr3_1types.KeyValueStateReadWriter) stateReadWriter {
	return stateReadWriter{stateReader: stateReader{kv: kv}, kv: kv}
}

func (s stateReadWriter) writeValue[M proto.Message](key string, m M) error {
	val, err := proto.Marshal(m)
	if err != nil {
		return err
	}
	return s.kv.Write([]byte(key), val)
}

var _ heap.Interface = (*priorityQueue)(nil)

// priorityQueue is a heap of coordinated events.
type priorityQueue []*pb.CoordinatedEvent

func (p *priorityQueue) Len() int { return len(*p) }

func (p *priorityQueue) Less(i, j int) bool {
	a, b := (*p)[i], (*p)[j]
	return cmp.Or(cmp.Compare(a.Deadline.Seconds, b.Deadline.Seconds),
		cmp.Compare(a.Deadline.Nanos, b.Deadline.Nanos)) < 0
}

func (p *priorityQueue) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *priorityQueue) Push(x any) {
	*p = append(*p, x.(*pb.CoordinatedEvent))
}

func (p *priorityQueue) Pop() any {
	return (*p)[len(*p)-1]
}

// PopCap pops up to capacity events from the queue, each of which must have wf capacity. Mutates capPerWF.
func (p *priorityQueue) PopCap(capacity uint32, capPerWF map[string]uint32) ([]*pb.CoordinatedEvent, error) {
	var events []*pb.CoordinatedEvent
	var wfLimited []*pb.CoordinatedEvent // omitted due to workflow limit
	for len(events) < int(capacity) {
		if len(*p) == 0 {
			break
		}
		next, ok := heap.Pop(p).(*pb.CoordinatedEvent)
		if !ok {
			return nil, fmt.Errorf("failed to cast *pb.CoordinatedEvent: got type %T", next)
		}
		if capPerWF[next.WorkflowID] == 0 {
			wfLimited = append(wfLimited, next)
			continue
		}
		capPerWF[next.WorkflowID]--
		events = append(events, next)
	}
	// restore omitted for next time
	for _, event := range wfLimited {
		heap.Push(p, event)
	}
	return events, nil
}

type proposedAndQueued struct {
	proposedEvents *pb.ConsensusQueueProposals
	queuedEvents   *priorityQueue
}

func (s proposedAndQueued) processObserved(f int, observedProposals map[string][]uint8, observedEvents map[string]*pb.CoordinatedEvent) {
	for id, nodes := range observedProposals {
		if observers, ok := s.proposedEvents.Observers[id]; !ok && len(nodes) >= f+1 {
			// new proposal straight to priority queue
			heap.Push(s.queuedEvents, s.proposedEvents.Events[id])
			continue
		} else if !ok {
			// new proposal < f + 1
			s.proposedEvents.Observers[id] = &pb.NodeIDList{Ids: nodes}
			s.proposedEvents.Events[id] = observedEvents[id]
			continue
		} else if len(nodes) < f+1 {
			// update proposal
			// TODO opt: build a map to avoid n^2? or nbd since n < f+1
			for _, node := range nodes {
				if !slices.Contains(observers.Ids, node) {
					observers.Ids = append(observers.Ids, node)
				}
			}
			if len(observers.Ids) < f+1 {
				continue
			}
		}
		// f+1 - promote proposal to priority queue
		delete(s.proposedEvents.Observers, id)
		delete(s.proposedEvents.Events, id)
		heap.Push(s.queuedEvents, s.proposedEvents.Events[id])
	}
}
