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
	observedEventsKey = "observed_events::index" // pb.ConsensusQueueObservedEvents { map[string]*pb.ConsensusQueueObservedEvent; [][]byte }
	// priority queue of events waiting for capacity
	queuedEventsKey = "queued_events::index" // pb.ConsensusQueueEventList { []*pb.ConsensusQueueEvent }
	// events the plugin committed to sending to nodes
	consensusEventsKey = "consensus_events::index" // pb.ConsensusQueueEventList { []*pb.ConsensusQueueEvent }
	// blob handles referenced from events in the PQ and consensus set
	blobHandlesKey = "blob_handles::index" // pb.ConsensusQueueBlobHandleList{ [][]byte }
)

type stateWriter struct {
	kv ocr3_1types.KeyValueStateReadWriter
}

func (s *stateWriter) writeValue[M proto.Message](key string, m M) error {
	val, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal %T: %w", m, err)
	}
	return s.kv.Write([]byte(key), val)
}

var _ heap.Interface = (*priorityQueue)(nil)

// priorityQueue is a heap of coordinated events, ordered by deadline and then event ID.
type priorityQueue []*pb.ConsensusQueueEvent

func (p *priorityQueue) Len() int { return len(*p) }

func (p *priorityQueue) Less(i, j int) bool {
	a, b := (*p)[i], (*p)[j]
	return cmpConsensusQueueEvents(a, b) < 0
}

func cmpConsensusQueueEvents(a, b *pb.ConsensusQueueEvent) int {
	return cmp.Or(cmp.Compare(a.Deadline.Seconds, b.Deadline.Seconds),
		cmp.Compare(a.Deadline.Nanos, b.Deadline.Nanos),
		cmp.Compare(a.TriggerEventID, b.TriggerEventID))
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

func (p *priorityQueue) pushObservations(timelessHash string, obs *pb.ConsensusQueueObservedEvent) *pb.ConsensusQueueEvent {
	event := &pb.ConsensusQueueEvent{
		WorkflowID:     obs.WorkflowID,
		TriggerEventID: obs.TriggerEventID,
		ObservedAt:     medianTimestamp(obs.ObservedAts),
		Deadline:       medianTimestamp(obs.Deadlines),
		TimelessHash:   []byte(timelessHash),
		Blobs:          obs.Blobs,
	}
	p.push(event)
	return event
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
func (p *priorityQueue) popLimited(capacity uint32, capPerWF map[string]uint32, expiration time.Time) []*pb.ConsensusQueueEvent {
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
	return events
}

type blobHandleMap map[string]map[string]uint32 // eventHash -> blobHandle -> idx

func (m blobHandleMap) add(eventHash string, blobHandle string, element uint32) {
	if _, ok := m[eventHash]; !ok {
		m[eventHash] = make(map[string]uint32)
	}
	m[eventHash][blobHandle] = element
}

func (m blobHandleMap) addAll(eventHash string, blobHandles map[string]uint32) {
	if _, ok := m[eventHash]; !ok {
		m[eventHash] = make(map[string]uint32)
	}
	for blobHandle, elementIndex := range blobHandles {
		m[eventHash][blobHandle] = elementIndex
	}
}

type kvState struct {
	// observedEventsKey
	observedEvents       map[string]*pb.ConsensusQueueObservedEvent // key: timeless hash
	observedEventEntries []*pb.ConsensusQueueObservedEventEntry     // only populated before write
	observedBlobs        blobHandleMap
	observedBlobHandles  [][]byte // only populated before write

	// priorityQueueKey
	queuedEvents *priorityQueue
	// consensusEventsKey
	consensusEvents []*pb.ConsensusQueueEvent
	// blobHandlesKey
	queuedBlobs       blobHandleMap
	queuedBlobHandles [][]byte // only populated before write
}

// loadFrom loads the kvState from the KV store, excluding consensus events, which are not required to transition.
func (s *kvState) loadFrom(kv ocr3_1types.KeyValueStateReadWriter) error {
	var err error
	s.observedEvents, s.observedBlobs, err = readObservedEvents(kv)
	if err != nil {
		return err
	}

	v, err := kv.Read([]byte(queuedEventsKey))
	if err != nil {
		return fmt.Errorf("failed to read queued events: %w", err)
	}
	if v != nil {
		var queued *pb.ConsensusQueueEventList
		if err = proto.Unmarshal(v, queued); err != nil {
			return fmt.Errorf("failed to unmarshal queued events: %w", err)
		}
		s.queuedEvents = new(priorityQueue(queued.GetEvents()))
	}

	s.queuedBlobs = make(blobHandleMap)
	queuedBlobs, err := readBlobHandles(kv)
	if err != nil {
		return err
	}
	for _, event := range []*pb.ConsensusQueueEvent(*s.queuedEvents) {
		for _, blob := range event.GetBlobs() {
			s.queuedBlobs.add(string(event.TimelessHash), string(queuedBlobs[blob.Blob]), blob.Element)
		}
	}

	return nil
}

func readObservedEvents(kv ocr3_1types.KeyValueStateReader) (observedEvents map[string]*pb.ConsensusQueueObservedEvent, observedBlobs blobHandleMap, err error) {
	v, err := kv.Read([]byte(observedEventsKey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read observed events: %w", err)
	}
	observedBlobs = make(blobHandleMap)
	if v != nil {
		var observed *pb.ConsensusQueueObservedEvents
		if err = proto.Unmarshal(v, observed); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal observed events: %w", err)
		}
		observedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)
		for _, e := range observed.GetObservations() {
			observedEvents[string(e.TimelessHash)] = e.Event
		}
		blobHandles := observed.GetBlobHandles()
		for eventHash, event := range observedEvents {
			for _, b := range event.GetBlobs() {
				observedBlobs.add(eventHash, string(blobHandles[b.Blob]), b.Element)
			}
		}
	} else {
		observedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)
	}
	return
}

// sizeOf returns the size of the state when serialized and written to the KV store via writeAll.
func (s *kvState) sizeOf() (size int) {
	s.prepareWrite()

	//TODO check maxKeyValueValueBytes too?

	protoOpts := proto.MarshalOptions{Deterministic: true}
	size += len(observedEventsKey) + len(queuedEventsKey) + len(consensusEventsKey) + len(blobHandlesKey)
	size += protoOpts.Size(&pb.ConsensusQueueObservedEvents{Observations: s.observedEventEntries, BlobHandles: s.observedBlobHandles})
	size += protoOpts.Size(&pb.ConsensusQueueEventList{Events: *s.queuedEvents})
	size += protoOpts.Size(&pb.ConsensusQueueEventList{Events: s.consensusEvents})
	size += protoOpts.Size(&pb.ConsensusQueueBlobHandleList{BlobHandles: s.queuedBlobHandles})

	return
}

// prepareWrite initializes observed and queued blob handles and updates the event references.
func (s *kvState) prepareWrite() {
	// gather blob handles and update events
	s.observedBlobHandles = nil
	observedBlobIndexes := make(map[string]int)
	s.observedEventEntries = nil
	//TODO deadline order so that we can drop efficiently
	for _, eventHash := range slices.Sorted(maps.Keys(s.observedEvents)) {
		event := s.observedEvents[eventHash]
		event.Blobs = nil // replace all
		for blobHandle, ei := range s.observedBlobs[eventHash] {
			idx, ok := observedBlobIndexes[blobHandle]
			if !ok {
				idx = len(s.observedBlobHandles)
				observedBlobIndexes[blobHandle] = idx
				s.observedBlobHandles = append(s.observedBlobHandles, []byte(blobHandle))
			}
			event.Blobs = append(event.Blobs, &pb.CoordinatedEventBlobHandle{
				Blob:    uint32(idx),
				Element: ei,
			})
		}
		s.observedEventEntries = append(s.observedEventEntries, &pb.ConsensusQueueObservedEventEntry{
			TimelessHash: []byte(eventHash), Event: event})
	}

	// gather blob handles and update events
	s.queuedBlobHandles = nil
	queuedBlobIndexes := make(map[string]int)
	addBlobs := func(blobs map[string]uint32) []*pb.CoordinatedEventBlobHandle {
		var blobHandles []*pb.CoordinatedEventBlobHandle
		for blobHandle, ei := range blobs {
			idx, ok := queuedBlobIndexes[blobHandle]
			if !ok {
				idx = len(s.queuedBlobHandles)
				queuedBlobIndexes[blobHandle] = idx
				s.queuedBlobHandles = append(s.queuedBlobHandles, []byte(blobHandle))
			}
			blobHandles = append(blobHandles, &pb.CoordinatedEventBlobHandle{
				Blob:    uint32(idx),
				Element: ei,
			})
		}
		return blobHandles
	}
	for _, event := range []*pb.ConsensusQueueEvent(*s.queuedEvents) {
		event.Blobs = addBlobs(s.queuedBlobs[string(event.TimelessHash)])
	}
	for _, event := range s.consensusEvents {
		event.Blobs = addBlobs(s.queuedBlobs[string(event.TimelessHash)])
	}
}

// writeAll writes full state to the KV store, including consensus events.
// prepareWrite must be called to include blob handles.
func (s *kvState) writeTo(w *stateWriter) error {
	//TODO MarshalOptions{UseCachedSize}?

	err := w.writeValue(observedEventsKey, &pb.ConsensusQueueObservedEvents{Observations: s.observedEventEntries, BlobHandles: s.observedBlobHandles})
	if err != nil {
		return fmt.Errorf("failed to write observed events: %w", err)
	}

	err = w.writeValue(queuedEventsKey, &pb.ConsensusQueueEventList{Events: *s.queuedEvents})
	if err != nil {
		return fmt.Errorf("failed to write queued events: %w", err)
	}

	err = w.writeValue(consensusEventsKey, &pb.ConsensusQueueEventList{Events: s.consensusEvents})
	if err != nil {
		return fmt.Errorf("failed to write consensus events: %w", err)
	}

	err = w.writeValue(blobHandlesKey, &pb.ConsensusQueueBlobHandleList{BlobHandles: s.queuedBlobHandles})
	if err != nil {
		return fmt.Errorf("failed to write blob handles: %w", err)
	}

	return nil
}

func (s *kvState) pruneExpiredObservations(expiration time.Time) {
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

// addObserved adds the observed events to the set, promoting to the priority queue on 2f+1 and calling onQueued.
func (s *kvState) addObservedEvents(f int, toAdd map[string]*pb.ConsensusQueueObservedEvent, observedBlobs blobHandleMap, onQueued func(event *pb.ConsensusQueueEvent)) {
	for eventHash, eventToAdd := range toAdd {
		obsEvent, ok := s.observedEvents[eventHash]
		if !ok && len(eventToAdd.NodeIds) >= 2*f+1 {
			// new proposal straight to priority queue
			s.queuedBlobs.addAll(eventHash, observedBlobs[eventHash])
			event := s.queuedEvents.pushObservations(eventHash, eventToAdd)
			onQueued(event)
		} else if !ok {
			// new proposal < f + 1
			s.observedBlobs.addAll(eventHash, observedBlobs[eventHash])
			s.observedEvents[eventHash] = eventToAdd
		} else {
			// update proposal
			// TODO opt: build a map to avoid n^2?
			for i, node := range eventToAdd.NodeIds {
				if !slices.Contains(obsEvent.NodeIds, node) {
					obsEvent.NodeIds = append(obsEvent.NodeIds, node)
					obsEvent.ObservedAts = append(obsEvent.ObservedAts, eventToAdd.ObservedAts[i])
					obsEvent.Deadlines = append(obsEvent.Deadlines, eventToAdd.Deadlines[i])
					obsEvent.Blobs = append(obsEvent.Blobs, eventToAdd.Blobs[i])
				} //TODO what if a node reports twice with different deadline or observation? do we need to overwrite?
			}

			if len(obsEvent.NodeIds) >= 2*f+1 {
				// promote proposal to priority queue
				delete(s.observedEvents, eventHash)
				s.queuedBlobs[eventHash] = s.observedBlobs[eventHash]
				delete(s.observedBlobs, eventHash)
				event := s.queuedEvents.pushObservations(eventHash, obsEvent)
				onQueued(event)
			}
		}
	}
}

// readConsensusState reads the consensus events and blob handles from the kvReader.
func readConsensusState(kvReader ocr3_1types.KeyValueStateReader) ([]*pb.ConsensusQueueEvent, [][]byte, error) {
	v, err := kvReader.Read([]byte(consensusEventsKey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read consensus events: %w", err)
	}

	var consensusEvents []*pb.ConsensusQueueEvent
	if v != nil {
		var pbs *pb.ConsensusQueueEventList
		if err := proto.Unmarshal(v, pbs); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal consensus events: %w", err)
		}
		consensusEvents = pbs.GetEvents()
	}

	blobHandles, err := readBlobHandles(kvReader)
	if err != nil {
		return nil, nil, err
	}
	return consensusEvents, blobHandles, nil
}

// readBlobHandles reads the blob handles from the kvReader.
func readBlobHandles(kvReader ocr3_1types.KeyValueStateReader) ([][]byte, error) {
	v, err := kvReader.Read([]byte(blobHandlesKey))
	if err != nil {
		return nil, fmt.Errorf("failed to read blob handles: %w", err)
	}

	var blobHandles [][]byte
	if v != nil {
		var bs *pb.ConsensusQueueBlobHandleList
		if err := proto.Unmarshal(v, bs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal consensus blob handles: %w", err)
		}
		blobHandles = bs.GetBlobHandles()
	}
	return blobHandles, err
}
