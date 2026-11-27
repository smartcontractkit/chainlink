package queue

import (
	"cmp"
	"container/heap"
	"fmt"
	"maps"
	"math"
	"slices"
	"sort"
	"time"

	"google.golang.org/protobuf/encoding/protowire"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3_1types"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3_1/plugins/queue/queuelimits"
)

// numKVKeys is the number of KV keys written each round.
const numKVKeys = queuelimits.NumKVKeys

// KV store keys
const (
	// observed events waiting for consensus (F+1 <= observers <2F+1). Each event's Blobs are aligned with its NodeIds: the
	// blob of each vote.
	observedEventsKey = queuelimits.ObservedEventsKey // pb.ConsensusQueueObservedEvents { map[string]*pb.ConsensusQueueObservedEvent; [][]byte }
	// priority queue of events waiting for capacity
	queuedEventsKey = queuelimits.QueuedEventsKey // pb.ConsensusQueueEventList { []*pb.ConsensusQueueEvent }
	// events the plugin committed to sending to nodes
	consensusEventsKey = queuelimits.ConsensusEventsKey // pb.ConsensusQueueEventList { []*pb.ConsensusQueueEvent }
	// blob handles referenced from events in the PQ and consensus set
	blobHandlesKey = queuelimits.BlobHandlesKey // pb.ConsensusQueueBlobHandleList{ [][]byte }
)

type stateWriter struct {
	kv    ocr3_1types.KeyValueStateReadWriter
	sizes map[string]int // key -> bytes written
}

func (s *stateWriter) writeValue[M proto.Message](key string, m M) error {
	val, err := proto.MarshalOptions{Deterministic: true}.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal %T: %w", m, err)
	}

	if s.sizes == nil {
		s.sizes = make(map[string]int)
	}
	s.sizes[key] = len(val)
	return s.kv.Write([]byte(key), val)
}

var _ heap.Interface = (*priorityQueue)(nil)

// priorityQueue is a heap of coordinated events, ordered by deadline and then trigger event ID.
type priorityQueue []*pb.ConsensusQueueEvent

func (p *priorityQueue) Len() int {
	if p == nil {
		return 0
	}

	return len(*p)
}

func (p *priorityQueue) Less(i, j int) bool {
	a, b := (*p)[i], (*p)[j]
	return a.CompareDeadline(b) < 0
}

func (p *priorityQueue) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *priorityQueue) Push(x any) {
	*p = append(*p, x.(*pb.ConsensusQueueEvent))
}

func (p *priorityQueue) Pop() any {
	if p.Len() == 0 {
		return nil
	}

	e := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return e
}

func (p *priorityQueue) push(x *pb.ConsensusQueueEvent) { heap.Push(p, x) }

func (s *kvState) queueObservations(timelessHash string, obs *pb.ConsensusQueueObservedEvent) *pb.ConsensusQueueEvent {
	event := &pb.ConsensusQueueEvent{
		WorkflowID:   obs.WorkflowID,
		ExecutionID:  obs.ExecutionID,
		ObservedAt:   medianTimestamp(obs.ObservedAts),
		Deadline:     medianTimestamp(obs.Deadlines),
		TimelessHash: []byte(timelessHash),
		// note: observedBlobs tracks blobs while in memory
	}

	s.queuedEvents.push(event)
	return event
}

// mustPop pops the next event from the queue, or returns nil if the queue is empty.
// The value must be of type *pb.ConsensusQueueEvent. Use push to ensure this.
func (p *priorityQueue) mustPop() *pb.ConsensusQueueEvent {
	if p.Len() == 0 {
		return nil
	}

	return heap.Pop(p).(*pb.ConsensusQueueEvent)
}

// popEligible pops up to n events from the queue, in priority order, each of which also has capacity left for its
// workflow, which it uses up. Expired events are dropped. Events without workflow capacity stay in the queue.
// Note: capPerWF is mutated.
//
// It makes a single pass, and puts back the events it skipped once at the end. A workflow's capacity only goes down
// while popping, so an event that was skipped stays ineligible, and looking at it again for each event popped (as popping
// one at a time did) is wasted work: quadratic in the events waiting on workflows without capacity (1.9s for 10000).
func (p *priorityQueue) popEligible(n int, capPerWF map[string]uint32, expiration time.Time) (popped []*pb.ConsensusQueueEvent, expired int) {
	var skipped []*pb.ConsensusQueueEvent
	for len(popped) < n {
		event := p.mustPop()
		if event == nil {
			break
		}

		if event.Deadline.AsTime().Before(expiration) {
			expired++
			continue // drop expired
		}

		if capPerWF[event.WorkflowID] == 0 {
			skipped = append(skipped, event)
			continue
		}

		capPerWF[event.WorkflowID]--
		popped = append(popped, event)
	}

	for _, event := range skipped {
		p.push(event)
	}

	return popped, expired
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

	maps.Copy(m[eventHash], blobHandles)
}

// voterBlobs is the blob each voter observed each event in: eventHash -> voter -> blob. There is one per voter, so a voter
// that votes again replaces its blob instead of adding to those that the event refers to, which would let it grow
// without limit.
type voterBlobs map[string]map[uint8]blobRef

type blobRef struct {
	handle  string
	element uint32 // the index of the event in the blob
}

func (m voterBlobs) add(eventHash string, voter uint8, blobHandle string, element uint32) {
	if _, ok := m[eventHash]; !ok {
		m[eventHash] = make(map[uint8]blobRef)
	}

	m[eventHash][voter] = blobRef{handle: blobHandle, element: element}
}

// handles returns the distinct blobs of the event's voters, and where the event is in each. Voters may share a blob,
// but if they disagree on where the event is in it, the lowest index is used so that the result is deterministic.
func (m voterBlobs) handles(eventHash string) map[string]uint32 {
	handles := make(map[string]uint32, len(m[eventHash]))
	for _, ref := range m[eventHash] {
		if element, ok := handles[ref.handle]; !ok || ref.element < element {
			handles[ref.handle] = ref.element
		}
	}

	return handles
}

// kvStats counts what a state transition removed from the state, other than by executing it.
type kvStats struct {
	expiredObserved int // observed events whose votes had all expired
	expiredQueued   int // queued events that had expired
	prunedObserved  int // observed events dropped to fit the state limits
	prunedQueued    int // queued events dropped to fit the state limits
}

type kvState struct {
	stats kvStats

	// observedEventsKey
	observedEvents       map[string]*pb.ConsensusQueueObservedEvent // key: timeless hash
	observedEventEntries []*pb.ConsensusQueueObservedEventEntry     // only populated before write
	observedBlobs        voterBlobs
	observedBlobHandles  [][]byte // only populated before write

	// priorityQueueKey
	queuedEvents priorityQueue
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

	queuedEvents, err := readQueuedEvents(kv)
	if err != nil {
		return err
	}

	s.queuedEvents = queuedEvents

	s.queuedBlobs = make(blobHandleMap)
	queuedBlobs, err := readBlobHandles(kv)
	if err != nil {
		return err
	}

	for _, event := range s.queuedEvents {
		for _, blob := range event.GetBlobs() {
			if int(blob.Blob) < len(queuedBlobs) {
				s.queuedBlobs.add(string(event.TimelessHash), string(queuedBlobs[blob.Blob]), blob.Element)
			}
		}
	}

	return nil
}

func readQueuedEvents(kv ocr3_1types.KeyValueStateReader) ([]*pb.ConsensusQueueEvent, error) {
	v, err := kv.Read([]byte(queuedEventsKey))
	if err != nil {
		return nil, fmt.Errorf("failed to read queued events: %w", err)
	}

	if v != nil {
		var queued pb.ConsensusQueueEventList
		if err = proto.Unmarshal(v, &queued); err != nil {
			return nil, fmt.Errorf("failed to unmarshal queued events: %w", err)
		}

		return queued.GetEvents(), nil
	}

	return nil, nil
}

func readObservedEvents(kv ocr3_1types.KeyValueStateReader) (observedEvents map[string]*pb.ConsensusQueueObservedEvent, observedBlobs voterBlobs, err error) {
	v, err := kv.Read([]byte(observedEventsKey))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read observed events: %w", err)
	}

	observedBlobs = make(voterBlobs)
	if v != nil {
		var observed pb.ConsensusQueueObservedEvents
		if err = proto.Unmarshal(v, &observed); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal observed events: %w", err)
		}

		observedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)
		for _, e := range observed.GetObservations() {
			observedEvents[string(e.TimelessHash)] = e.Event
		}

		blobHandles := observed.GetBlobHandles()
		for eventHash, event := range observedEvents {
			// Blobs[i] is the blob of the vote of NodeIds[i]
			for i, voter := range event.GetNodeIds() {
				if i < len(event.Blobs) && int(event.Blobs[i].Blob) < len(blobHandles) {
					observedBlobs.add(eventHash, voter, string(blobHandles[event.Blobs[i].Blob]), event.Blobs[i].Element)
				}
			}
		}
	} else {
		observedEvents = make(map[string]*pb.ConsensusQueueObservedEvent)
	}

	return observedEvents, observedBlobs, nil
}

// sizeOf returns the size of the state when serialized and written to the KV store via writeAll.
func (s *kvState) sizeOf() (size int) {
	size += len(observedEventsKey) + len(queuedEventsKey) + len(consensusEventsKey) + len(blobHandlesKey)
	size += s.sizeOfObservedEvents() + s.sizeOfQueuedEvents() + s.sizeOfConsensusEvents() + s.sizeOfBlobHandles()
	return
}

func (s *kvState) sizeOfObservedEvents() int {
	s.prepareWriteObservedEvents()
	return proto.MarshalOptions{Deterministic: true}.Size(&pb.ConsensusQueueObservedEvents{Observations: s.observedEventEntries, BlobHandles: s.observedBlobHandles})
}

func (s *kvState) sizeOfQueuedEvents() int {
	s.prepareWriteBlobHandles() // events must carry their blob refs to be sized accurately
	return proto.MarshalOptions{Deterministic: true}.Size(&pb.ConsensusQueueEventList{Events: s.queuedEvents})
}

func (s *kvState) sizeOfConsensusEvents() int {
	s.prepareWriteBlobHandles() // events must carry their blob refs to be sized accurately
	return proto.MarshalOptions{Deterministic: true}.Size(&pb.ConsensusQueueEventList{Events: s.consensusEvents})
}

func (s *kvState) sizeOfBlobHandles() int {
	s.prepareWriteBlobHandles()
	return proto.MarshalOptions{Deterministic: true}.Size(&pb.ConsensusQueueBlobHandleList{BlobHandles: s.queuedBlobHandles})
}

// prepareWrite initializes observed blob handles & entries, and updates the event references.
func (s *kvState) prepareWriteObservedEvents() {
	// gather blob handles and update events
	s.observedBlobHandles = nil
	observedBlobIndexes := make(map[string]int)
	s.observedEventEntries = nil

	// sort so that we can drop lowest priority efficiently
	for _, eventHash := range slices.SortedFunc(maps.Keys(s.observedEvents), func(a string, b string) int {
		ea, eb := s.observedEvents[a], s.observedEvents[b]
		// Most votes first, so the events dropped first are those furthest from the F+1 they need. The deadline is not
		// used: until an event has F+1 votes it has no value that is certain to be honest to take it from.
		return cmp.Or(cmp.Compare(len(eb.NodeIds), len(ea.NodeIds)), cmp.Compare(a, b))
	}) {
		event := s.observedEvents[eventHash]
		event.Blobs = nil // replace all, with one per vote
		for _, voter := range event.NodeIds {
			ref, ok := s.observedBlobs[eventHash][voter]
			if !ok { // every vote has a blob, but keep Blobs aligned with NodeIds regardless
				event.Blobs = append(event.Blobs, &pb.CoordinatedEventBlobHandle{Blob: math.MaxUint32})
				continue
			}

			idx, ok := observedBlobIndexes[ref.handle]
			if !ok {
				idx = len(s.observedBlobHandles)
				observedBlobIndexes[ref.handle] = idx
				s.observedBlobHandles = append(s.observedBlobHandles, []byte(ref.handle))
			}

			event.Blobs = append(event.Blobs, &pb.CoordinatedEventBlobHandle{
				Blob:    uint32(idx), //nolint:gosec // G115: an index into at most maxObservedBlobHandles blob handles
				Element: ref.element,
			})
		}

		s.observedEventEntries = append(s.observedEventEntries, &pb.ConsensusQueueObservedEventEntry{
			TimelessHash: []byte(eventHash), Event: event})
	}
}

// prepareWrite initializes queued blob handles and updates the event references.
func (s *kvState) prepareWriteBlobHandles() {
	// gather blob handles and update events
	s.queuedBlobHandles = nil
	queuedBlobIndexes := make(map[string]int)
	addBlobs := func(blobs map[string]uint32) []*pb.CoordinatedEventBlobHandle {
		if len(blobs) == 0 {
			return nil
		}

		blobHandles := make([]*pb.CoordinatedEventBlobHandle, 0, len(blobs))
		for _, blobHandle := range slices.Sorted(maps.Keys(blobs)) {
			ei := blobs[blobHandle]
			idx, ok := queuedBlobIndexes[blobHandle]
			if !ok {
				idx = len(s.queuedBlobHandles)
				queuedBlobIndexes[blobHandle] = idx
				s.queuedBlobHandles = append(s.queuedBlobHandles, []byte(blobHandle))
			}

			blobHandles = append(blobHandles, &pb.CoordinatedEventBlobHandle{
				Blob:    uint32(idx), //nolint:gosec // G115: an index into at most maxObservedBlobHandles blob handles
				Element: ei,
			})
		}

		return blobHandles
	}

	for _, event := range s.queuedEvents {
		event.Blobs = addBlobs(s.queuedBlobs[string(event.TimelessHash)])
	}

	for _, event := range s.consensusEvents {
		event.Blobs = addBlobs(s.queuedBlobs[string(event.TimelessHash)])
	}
}

// trimConsensusEvents moves the lowest priority consensus events (the last ones) back to the queue until the consensus
// events and the blob handles fit: each in a value, and together in half of maxTotalBytes. The blob handles cover every
// queued and consensus event, so they are the same wherever an event sits, and the events can be sized once.
func (s *kvState) trimConsensusEvents(maxTotalBytes int) {
	s.prepareWriteBlobHandles() // events carry their blob refs, and the handles are collected
	blobs := proto.MarshalOptions{Deterministic: true}.Size(&pb.ConsensusQueueBlobHandleList{BlobHandles: s.queuedBlobHandles})

	var events, n int
	for ; n < len(s.consensusEvents); n++ {
		next := events + protowire.SizeTag(1) + protowire.SizeBytes(proto.Size(s.consensusEvents[n])) // as an element of the list
		if next > maxKeyValueValueBytes || blobs > maxKeyValueValueBytes || next+blobs > maxTotalBytes/2 {
			break
		}

		events = next
	}

	for _, event := range s.consensusEvents[n:] {
		s.queuedEvents.push(event)
	}

	s.consensusEvents = s.consensusEvents[:n]
}

// queuedWorkflowCapacities returns a capacity of n for every workflow that has a queued event.
func (s *kvState) queuedWorkflowCapacities(n uint32) map[string]uint32 {
	caps := make(map[string]uint32)
	for _, event := range s.queuedEvents {
		caps[event.WorkflowID] = n
	}
	return caps
}

// popConsensusEvents pops as many queuedEvents as possible in to consensusEvents, up to capacity (overall, and per
// workflow) or the KV limits. Expired events are dropped.
func (s *kvState) popConsensusEvents(capacity int, capPerWF map[string]uint32, expiration time.Time, maxTotalBytes int) {
	popped, expired := s.queuedEvents.popEligible(capacity-len(s.consensusEvents), capPerWF, expiration)
	s.stats.expiredQueued += expired
	s.consensusEvents = append(s.consensusEvents, popped...)
	s.trimConsensusEvents(maxTotalBytes)
}

// pruneObservedToFit drops the lowest priority observed events until that value, and the overall size, fit.
//
// Size only ever shrinks as more events are dropped, so the fewest to drop is found by binary search, rather than by
// dropping one at a time and measuring the whole state again each time (which is quadratic, and took 300ms to drop 500).
func (s *kvState) pruneObservedToFit(maxTotalBytes int) {
	s.prepareWriteObservedEvents() // populates the entries, sorted by decreasing priority, so the last are dropped first
	fits := func() bool { return s.sizeOfObservedEvents() <= maxKeyValueValueBytes && s.sizeOf() <= maxTotalBytes }
	if len(s.observedEventEntries) == 0 || fits() {
		return
	}

	entries, all := s.observedEventEntries, s.observedEvents
	keepFirst := func(n int) { // n is how many of the highest priority events to keep
		kept := make(map[string]*pb.ConsensusQueueObservedEvent, n)
		for _, e := range entries[:n] {
			kept[string(e.TimelessHash)] = all[string(e.TimelessHash)]
		}
		s.observedEvents = kept
	}
	// the most that fit, in terms of those dropped: fewest dropped is the most kept
	dropped := sort.Search(len(entries), func(drop int) bool {
		keepFirst(len(entries) - drop)
		return fits()
	})
	keepFirst(len(entries) - dropped)
	s.stats.prunedObserved += dropped
}

// pruneQueuedToFit drops the lowest priority events of those that were newly queued until the queued events, the blob
// handles, and the overall size fit. As for pruneObservedToFit, the fewest to drop is found by binary search.
func (s *kvState) pruneQueuedToFit(newlyQueued []string, maxTotalBytes int) {
	fits := func() bool {
		return s.sizeOfQueuedEvents() <= maxKeyValueValueBytes && s.sizeOfBlobHandles() <= maxKeyValueValueBytes && s.sizeOf() <= maxTotalBytes
	}
	if len(newlyQueued) == 0 || fits() {
		return
	}

	newlyQueuedSet := map[string]struct{}{}
	for _, hash := range newlyQueued {
		newlyQueuedSet[hash] = struct{}{}
	}
	eligibleEvents := make([]*pb.ConsensusQueueEvent, 0, len(s.queuedEvents))
	for _, e := range s.queuedEvents {
		if _, ok := newlyQueuedSet[string(e.TimelessHash)]; ok {
			eligibleEvents = append(eligibleEvents, e)
		}
	}
	// lowest priority (latest deadline) first, which is the order they are dropped in
	slices.SortFunc(eligibleEvents, func(a, b *pb.ConsensusQueueEvent) int { return b.CompareDeadline(a) })

	queued := s.queuedEvents
	drop := func(n int) {
		dropped := make(map[*pb.ConsensusQueueEvent]struct{}, n)
		for _, e := range eligibleEvents[:n] {
			dropped[e] = struct{}{}
		}
		s.queuedEvents = slices.DeleteFunc(slices.Clone(queued), func(e *pb.ConsensusQueueEvent) bool { _, ok := dropped[e]; return ok })
	}
	n := sort.Search(len(eligibleEvents), func(n int) bool {
		drop(n)
		return fits()
	})
	drop(n)
	s.stats.prunedQueued += n

	// rebuild the priority queue
	heap.Init(&s.queuedEvents)
}

// writeAll writes full state to the KV store, including consensus events.
// prepareWrite must be called to include blob handles.
func (s *kvState) writeTo(w *stateWriter) error {
	s.prepareWriteObservedEvents()
	s.prepareWriteBlobHandles()

	err := w.writeValue(observedEventsKey, &pb.ConsensusQueueObservedEvents{Observations: s.observedEventEntries, BlobHandles: s.observedBlobHandles})
	if err != nil {
		return fmt.Errorf("failed to write observed events: %w", err)
	}

	err = w.writeValue(queuedEventsKey, &pb.ConsensusQueueEventList{Events: s.queuedEvents})
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

// pruneExpiredObservations deletes expired observations, but does not evict on <f+1 since we need to add observations first.
func (s *kvState) pruneExpiredObservations(expiration time.Time) {
	maps.DeleteFunc(s.observedEvents, func(k string, v *pb.ConsensusQueueObservedEvent) bool {
		// a deadline equal to the expiration is not expired, consistent with popLimited
		if !slices.ContainsFunc(v.Deadlines, func(t *timestamppb.Timestamp) bool {
			return !t.AsTime().Before(expiration)
		}) {
			delete(s.observedBlobs, k)
			s.stats.expiredObserved++
			return true // all deadlines expired, so delete the whole thing
		}

		// some expired, so prune them
		for i := 0; i < len(v.Deadlines); {
			if v.Deadlines[i].AsTime().Before(expiration) {
				delete(s.observedBlobs[k], v.NodeIds[i]) // the vote goes, and so does its blob
				v.NodeIds = slices.Delete(v.NodeIds, i, i+1)
				v.ObservedAts = slices.Delete(v.ObservedAts, i, i+1)
				v.Deadlines = slices.Delete(v.Deadlines, i, i+1)
			} else {
				i++
			}
		}

		return false
	})
}

// pruneObservations removes all observed events that have fewer than f+1 observations.
// This is necessary to prevent long expirations from outliving pruned honest ones.
func (s *kvState) pruneObservations(f int) {
	maps.DeleteFunc(s.observedEvents, func(k string, v *pb.ConsensusQueueObservedEvent) bool {
		return len(v.NodeIds) < f+1
	})
}

// addObserved adds the observed events to the set, promoting to the priority queue on 2f+1 and calling onQueued.
// Events with fewer than f+1 observations are ignored.
func (s *kvState) addObservedEvents(f int, toAdd map[string]*pb.ConsensusQueueObservedEvent, observedBlobs voterBlobs, onQueued func(event *pb.ConsensusQueueEvent)) {
	alreadyQueued := make(map[string]struct{}, len(s.queuedEvents))
	for _, e := range s.queuedEvents {
		alreadyQueued[string(e.TimelessHash)] = struct{}{}
	}

	for _, eventHash := range slices.Sorted(maps.Keys(toAdd)) {
		if _, ok := alreadyQueued[eventHash]; ok {
			continue // already waiting for capacity; promoting it again would duplicate it
		}

		eventToAdd := toAdd[eventHash]
		obsEvent, ok := s.observedEvents[eventHash]
		if !ok {
			switch {
			case len(eventToAdd.NodeIds) >= 2*f+1:
				// new proposal straight to priority queue
				s.queuedBlobs.addAll(eventHash, observedBlobs.handles(eventHash))
				event := s.queueObservations(eventHash, eventToAdd)
				onQueued(event)
			case len(eventToAdd.NodeIds) >= f+1:
				// new proposal to the observation set
				s.observedBlobs[eventHash] = maps.Clone(observedBlobs[eventHash])
				s.observedEvents[eventHash] = eventToAdd
			default:
				// new proposal without enough observations < f + 1
			}
		} else {
			// update proposal
			// (quadratic in the number of nodes, which is small)
			for i, node := range eventToAdd.NodeIds {
				if prev := slices.Index(obsEvent.NodeIds, node); prev == -1 {
					obsEvent.NodeIds = append(obsEvent.NodeIds, node)
					obsEvent.ObservedAts = append(obsEvent.ObservedAts, eventToAdd.ObservedAts[i])
					obsEvent.Deadlines = append(obsEvent.Deadlines, eventToAdd.Deadlines[i])
				} else { // overwrite previous observation
					obsEvent.ObservedAts[prev] = eventToAdd.ObservedAts[i]
					obsEvent.Deadlines[prev] = eventToAdd.Deadlines[i]
				}

				// the voter's blob, which replaces its previous one
				if ref, ok := observedBlobs[eventHash][node]; ok {
					s.observedBlobs.add(eventHash, node, ref.handle, ref.element)
				}
			}

			if len(obsEvent.NodeIds) >= 2*f+1 {
				// promote proposal to priority queue
				delete(s.observedEvents, eventHash)
				s.queuedBlobs.addAll(eventHash, s.observedBlobs.handles(eventHash))
				delete(s.observedBlobs, eventHash)
				event := s.queueObservations(eventHash, obsEvent)
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
		var pbs pb.ConsensusQueueEventList
		if err = proto.Unmarshal(v, &pbs); err != nil {
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
		var bs pb.ConsensusQueueBlobHandleList
		if err = proto.Unmarshal(v, &bs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal consensus blob handles: %w", err)
		}

		blobHandles = bs.GetBlobHandles()
	}

	return blobHandles, err
}
