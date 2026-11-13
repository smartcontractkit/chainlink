package queue

import (
	"cmp"
	"container/heap"
	"fmt"
	"maps"
	"slices"
	"time"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// KV store keys
const (
	// observed events waiting for consensus (<2f+1 observers)
	observedEventsKey = "observed_events::index" // pb.ConsensusQueueObservedEvents { map[string]*pb.ConsensusQueueObservedEvent }
	// priority queue of events waiting for capacity
	priorityQueueKey = "priority_queue::index" // pb.ConsensusQueueEventList { []*pb.ConsensusQueueEvent }
	// events the plugin committed to sending to nodes
	consensusEventsKey = "consensus_events::index" // pb.ConsensusQueueEventList { []*pb.ConsensusQueueEvent }
)

type stateReader struct {
	kv ocr3_1types.KeyValueStateReader
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
		return fmt.Errorf("failed to marshal %T: %w", m, err)
	}
	return s.kv.Write([]byte(key), val)
}

// readObservedAndQueued reads both observed events and priority queue events from the KV store.
func (s stateReader) readObservedAndQueued() (*observedAndQueued, error) {
	var obsQ observedAndQueued
	v, err := s.kv.Read([]byte(observedEventsKey))
	if err != nil {
		return nil, fmt.Errorf("failed to read observed events: %w", err)
	}
	if v != nil {
		var observed *pb.ConsensusQueueObservedEvents
		if err = proto.Unmarshal(v, observed); err != nil {
			return nil, err
		}
		obsQ.observedEvents = observed.GetObservations()
	} else {
		obsQ.observedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)
	}

	v, err = s.kv.Read([]byte(priorityQueueKey))
	if err != nil {
		return nil, fmt.Errorf("failed to read queued events: %w", err)
	}
	if v != nil {
		var queued *pb.ConsensusQueueEventList
		if err = proto.Unmarshal(v, queued); err != nil {
			return nil, err
		}
		obsQ.queuedEvents = new(priorityQueue(queued.GetEvents()))
	}

	return &obsQ, nil
}

func (s stateReadWriter) writeAll(state *observedAndQueued, consensus []*pb.ConsensusQueueEvent) error {
	err := s.writeValue(observedEventsKey, &pb.ConsensusQueueObservedEvents{Observations: state.observedEvents})
	if err != nil {
		return fmt.Errorf("failed to write observed events: %w", err)
	}
	err = s.writeValue(priorityQueueKey, &pb.ConsensusQueueEventList{Events: *state.queuedEvents})
	if err != nil {
		return fmt.Errorf("failed to write queued events: %w", err)
	}
	err = s.writeValue(consensusEventsKey, &pb.ConsensusQueueEventList{Events: consensus})
	if err != nil {
		return fmt.Errorf("failed to write consensus events: %w", err)
	}
	return nil
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

func (p *priorityQueue) push(x *pb.ConsensusQueueEvent) { heap.Push(p, x) }

func (p *priorityQueue) pushObservations(blobHash string, event *pb.ConsensusQueueObservedEvent) {
	p.push(&pb.ConsensusQueueEvent{
		WorkflowID:     event.WorkflowID,
		TriggerEventID: event.TriggerEventID,
		ObservedAt:     medianTimestamp(event.ObservedAts),
		Deadline:       medianTimestamp(event.Deadlines),
		BlobHash:       []byte(blobHash),
		Blobs:          event.Blobs,
	})
}

// mustPop pops the next event from the queue, or returns nil if the queue is empty.
// The value must be of type *pb.ConsensusQueueEvent. Use push to ensure this.
func (p *priorityQueue) mustPop() *pb.ConsensusQueueEvent {
	if len(*p) == 0 {
		return nil
	}
	return heap.Pop(p).(*pb.ConsensusQueueEvent)
}

// popLimited pops up to capacity events from the queue, each of which must also have sufficient wf capacity.
// Note: capPerWF is mutated.
func (p *priorityQueue) popLimited(capacity uint32, capPerWF map[string]uint32, expiration time.Time) ([]*pb.ConsensusQueueEvent, error) {
	var events []*pb.ConsensusQueueEvent
	var wfLimited []*pb.ConsensusQueueEvent // omitted due to workflow limit
	for len(events) < int(capacity) {
		next := p.mustPop()
		if next == nil {
			break
		}
		if next.Deadline.AsTime().Before(expiration) {
			// drop expired
			continue
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
		p.push(event)
	}
	return events, nil
}

type observedAndQueued struct {
	observedEvents map[string]*pb.ConsensusQueueObservedEvent // key: blobHash
	queuedEvents   *priorityQueue
}

func (s *observedAndQueued) pruneExpiredObservations(expiration time.Time) {
	unexpiredFn := func(t *timestamppb.Timestamp) bool { return t.AsTime().After(expiration) }
	maps.DeleteFunc(s.observedEvents, func(k string, v *pb.ConsensusQueueObservedEvent) bool {
		if !slices.ContainsFunc(v.Deadlines, unexpiredFn) {
			return false // no unexpired deadlines, so delete the whole thing
		}

		// some expired, so prune them
		var toDelete []int
		for i, d := range v.Deadlines {
			if !unexpiredFn(d) {
				toDelete = append(toDelete, i)
			}
		}
		for _, i := range toDelete {
			v.NodeIds = slices.Delete(v.NodeIds, i, i+1)
			v.ObservedAts = slices.Delete(v.ObservedAts, i, i+1)
			v.Deadlines = slices.Delete(v.Deadlines, i, i+1)
			v.Blobs = slices.Delete(v.Blobs, i, i+1)
		}

		return false
	})
}

// TODO on promotion to priority queue, ack the trigger coordinator
func (s *observedAndQueued) addObservedEvents(f int, toAdd map[string]*pb.ConsensusQueueObservedEvent) {
	for blobHash, eventToAdd := range toAdd {
		obsEvent, ok := s.observedEvents[blobHash]
		if !ok && len(eventToAdd.NodeIds) >= 2*f+1 {
			// new proposal straight to priority queue
			s.queuedEvents.pushObservations(blobHash, eventToAdd)
		} else if !ok {
			// new proposal < f + 1
			s.observedEvents[blobHash] = eventToAdd
		} else {
			// update proposal
			// TODO opt: build a map to avoid n^2?
			for i, node := range eventToAdd.NodeIds {
				if !slices.Contains(obsEvent.NodeIds, node) {
					obsEvent.NodeIds = append(obsEvent.NodeIds, node)
					obsEvent.ObservedAts = append(obsEvent.ObservedAts, eventToAdd.ObservedAts[i])
					obsEvent.Deadlines = append(obsEvent.Deadlines, eventToAdd.Deadlines[i])
					obsEvent.Blobs = append(obsEvent.Blobs, eventToAdd.Blobs[i])
				} //TODO what if a node reports twice with different deadline or nodeID? do we need to overwrite?
			}
			if len(obsEvent.NodeIds) >= 2*f+1 {
				// promote proposal to priority queue
				delete(s.observedEvents, blobHash)
				s.queuedEvents.pushObservations(blobHash, obsEvent)
			}
		}
	}
}
