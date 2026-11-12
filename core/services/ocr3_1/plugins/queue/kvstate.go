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

// KV store Keys
const (
	observedEventsKey  = "observed_events::index"  // observed events waiting for consensus (<2f+1 observers)
	priorityQueueKey   = "priority_queue::index"   // priority queue of events waiting for capacity
	consensusEventsKey = "consensus_events::index" // events the plugin committed to sending to nodes
)

// readObservedAndQueued reads both observed events and priority queue events from the KV store.
func (s stateReader) readObservedAndQueued() (*observedAndQueued, error) {
	var obsQ observedAndQueued
	v, err := s.kv.Read([]byte(observedEventsKey))
	if err != nil {
		return nil, err
	}
	if v != nil {
		var queued *pb.ConsensusQueueObservedEvents
		if err = proto.Unmarshal(v, queued); err != nil {
			return nil, err
		}
		obsQ.observedEvents = queued.Observations
	} else {
		obsQ.observedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)
	}

	v, err = s.kv.Read([]byte(priorityQueueKey))
	if err != nil {
		return nil, err
	}
	if v != nil {
		var consensus *pb.ConsensusQueueEventList
		if err = proto.Unmarshal(v, consensus); err != nil {
			return nil, err
		}
		obsQ.queuedEvents = new(priorityQueue(consensus.Events))
	}

	return &obsQ, nil
}

func (s stateReadWriter) update(state *observedAndQueued, consensus []*pb.ConsensusQueueEvent) error {
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

func (p *priorityQueue) push(x *pb.ConsensusQueueEvent) { heap.Push(p, x) }

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
			v.Deadlines = slices.Delete(v.Deadlines, i, i+1)
			v.ObservedAts = slices.Delete(v.ObservedAts, i, i+1)
			v.NodeIds = slices.Delete(v.NodeIds, i, i+1)
		}

		return false
	})
}

// TODO on promotion to priority queue, ack the trigger coordinator
func (s *observedAndQueued) addObservedEvents(f int, observedEvents map[string]*pb.ConsensusQueueObservedEvent) {
	for blobHash, event := range observedEvents {
		obs, ok := s.observedEvents[blobHash]
		if !ok && len(event.NodeIds) >= 2*f+1 {
			// new proposal straight to priority queue
			s.queuedEvents.push(&pb.ConsensusQueueEvent{
				WorkflowID:     event.WorkflowID,
				TriggerEventID: event.TriggerEventID,
				ObservedAt:     medianTimestamp(event.ObservedAts),
				Deadline:       medianTimestamp(event.Deadlines),
				Blobs:          event.Blobs,
			})
		} else if !ok {
			// new proposal < f + 1
			s.observedEvents[blobHash] = event
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
			if len(obs.NodeIds) >= 2*f+1 {
				// promote proposal to priority queue
				delete(s.observedEvents, blobHash)
				s.queuedEvents.push(&pb.ConsensusQueueEvent{
					WorkflowID:     event.WorkflowID,
					TriggerEventID: event.TriggerEventID,
					ObservedAt:     medianTimestamp(event.ObservedAts),
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
