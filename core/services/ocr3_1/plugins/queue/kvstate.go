package queue

import (
	"cmp"
	"container/heap"
	"fmt"
	"slices"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// KV store Keys
const (
	observedEventsKey  = "observed_events::index"  // observed events waiting for consensus (<f+1 observers)
	priorityQueueKey   = "priority_queue::index"   // priority queue of events waiting for capacity
	consensusEventsKey = "consensus_events::index" // events the plugin committed to sending to nodes
)

func (s stateReader) readObservedEvents() (*pb.ConsensusQueueObservedEvents, error) {
	return s.readValue[*pb.ConsensusQueueObservedEvents](observedEventsKey)
}

func (s stateReadWriter) writeObservedEvents(events *pb.ConsensusQueueObservedEvents) error {
	return s.writeValue(observedEventsKey, events)
}

func (s stateReader) readPriorityQueue() (priorityQueue, error) {
	list, err := s.readValue[*pb.ConsensusQueueEventList](priorityQueueKey)
	if err != nil {
		return nil, err
	}
	return list.Events, nil
}

func (s stateReadWriter) writePriorityQueue(queue priorityQueue) error {
	return s.writeValue(priorityQueueKey, &pb.ConsensusQueueEventList{Events: queue})
}

func (s stateReader) readConsensusEvents() (*pb.ConsensusQueueEventList, error) {
	return s.readValue[*pb.ConsensusQueueEventList](consensusEventsKey)
}

func (s stateReadWriter) writeConsensusEvents(events *pb.ConsensusQueueEventList) error {
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
	val, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
	if err != nil {
		return err
	}
	return s.kv.Write([]byte(key), val)
}

var _ heap.Interface = (*priorityQueue)(nil)

// priorityQueue is a heap of coordinated events, ordered by deadline and then event ID.
type priorityQueue []*pb.ConsensusQueueEvent

func (p *priorityQueue) Len() int { return len(*p) }

func (p *priorityQueue) Less(i, j int) bool {
	a, b := (*p)[i], (*p)[j]

	return cmp.Or(cmp.Compare(a.Deadline.Seconds, b.Deadline.Seconds),
		cmp.Compare(a.Deadline.Nanos, b.Deadline.Nanos),
		cmp.Compare(a.TriggerEventID, b.TriggerEventID)) < 0
}

func (p *priorityQueue) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *priorityQueue) Push(x any) {
	*p = append(*p, x.(*pb.ConsensusQueueEvent))
}

func (p *priorityQueue) Pop() any {
	return (*p)[len(*p)-1]
}

// PopCap pops up to capacity events from the queue, each of which must have wf capacity. Mutates capPerWF.
func (p *priorityQueue) PopCap(capacity uint32, capPerWF map[string]uint32) ([]*pb.ConsensusQueueEvent, error) {
	var events []*pb.ConsensusQueueEvent
	var wfLimited []*pb.ConsensusQueueEvent // omitted due to workflow limit
	for len(events) < int(capacity) {
		if len(*p) == 0 {
			break
		}
		next, ok := heap.Pop(p).(*pb.ConsensusQueueEvent)
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

type observedAndQueued struct {
	observedEvents *pb.ConsensusQueueObservedEvents
	queuedEvents   *priorityQueue
}

func (s observedAndQueued) addObservedEvents(f int, observedEvents map[string]*pb.ConsensusQueueObservedEvent) {
	for id, event := range observedEvents {
		obs, ok := s.observedEvents.Observations[id]
		if !ok && len(event.NodeIds) >= f+1 {
			// new proposal straight to priority queue
			heap.Push(s.queuedEvents, &pb.ConsensusQueueEvent{
				WorkflowID:     "", //TODO parse?
				TriggerEventID: "", //TODO parse?
				Deadline:       medianTimestamp(event.Deadlines),
				Blobs:          event.Blobs,
			})
		} else if !ok {
			// new proposal < f + 1
			s.observedEvents.Observations[id] = event
		} else {
			// update proposal
			// TODO opt: build a map to avoid n^2?
			for i, node := range event.NodeIds {
				if !slices.Contains(obs.NodeIds, node) {
					obs.NodeIds = append(obs.NodeIds, node)
					obs.Deadlines = append(obs.Deadlines, event.Deadlines[i])
					obs.Blobs = append(obs.Blobs, event.Blobs[i])
				} //TODO what if a node reports twice with different deadline or nodeID? do we need to overwrite?
			}
			if len(obs.NodeIds) >= f+1 {
				// promote proposal to priority queue
				delete(s.observedEvents.Observations, id)
				heap.Push(s.queuedEvents, &pb.ConsensusQueueEvent{
					WorkflowID:     "", //TODO
					TriggerEventID: "", //TODO
					Deadline:       medianTimestamp(event.Deadlines),
					Blobs:          event.Blobs,
				})
			}
		}
	}
}

func medianTimestamp(s []*timestamppb.Timestamp) *timestamppb.Timestamp {
	var times []time.Time
	for _, t := range s {
		times = append(times, time.Unix(t.Seconds, int64(t.Nanos)))
	}
	slices.SortFunc(times, time.Time.Compare)

	a, b := times[len(times)/2], times[len(times)/2-1]
	diff := a.Sub(b)
	return timestamppb.New(a.Add(diff / 2))
}
