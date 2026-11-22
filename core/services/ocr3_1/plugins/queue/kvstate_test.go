package queue

import (
	"container/heap"
	"fmt"
	"maps"
	"math/rand"
	"slices"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
)

// memKV is an in-memory ocr3_1types.KeyValueStateReadWriter.
type memKV map[string][]byte

func (m memKV) Read(k []byte) ([]byte, error) { return m[string(k)], nil }
func (m memKV) Write(k, v []byte) error {
	if v == nil {
		v = []byte{}
	}
	m[string(k)] = v
	return nil
}
func (m memKV) Delete(k []byte) error { delete(m, string(k)); return nil }

func qEvent(wf, execID, hash string, deadlineSec int64) *pb.ConsensusQueueEvent {
	return &pb.ConsensusQueueEvent{
		WorkflowID: wf, ExecutionID: execID, TimelessHash: []byte(hash),
		ObservedAt: ts(deadlineSec - 1), Deadline: ts(deadlineSec),
	}
}

// obsEvent is an observed event seen by the given nodes, each with deadline base+nodeID.
func obsEvent(wf, execID string, base int64, nodes ...byte) *pb.ConsensusQueueObservedEvent {
	e := &pb.ConsensusQueueObservedEvent{WorkflowID: wf, ExecutionID: execID, NodeIds: nodes}
	for _, n := range nodes {
		e.ObservedAts = append(e.ObservedAts, ts(base+int64(n)-1))
		e.Deadlines = append(e.Deadlines, ts(base+int64(n)))
	}
	return e
}

func newKVState() *kvState {
	return &kvState{
		observedEvents: map[string]*pb.ConsensusQueueObservedEvent{},
		observedBlobs:  voterBlobs{},
		queuedBlobs:    blobHandleMap{},
	}
}

func TestPriorityQueue_Order(t *testing.T) {
	t.Parallel()

	var pq priorityQueue
	pq.push(qEvent("wf", "c", "h3", 30))
	pq.push(qEvent("wf", "a", "h2", 10)) // same deadline as next, ordered by hash
	pq.push(qEvent("wf", "b", "h1", 10))
	pq.push(qEvent("wf", "d", "h4", 20))

	var got []string
	for pq.Len() > 0 {
		got = append(got, string(pq.mustPop().TimelessHash))
	}
	require.Equal(t, []string{"h1", "h2", "h4", "h3"}, got)
	require.Nil(t, pq.mustPop(), "empty queue pops nil")
}

func TestPopLimited(t *testing.T) {
	t.Parallel()

	now := int64(100)
	expiration := time.Unix(now, 0)

	t.Run("drops expired", func(t *testing.T) {
		t.Parallel()
		var pq priorityQueue
		pq.push(qEvent("wf", "old", "h1", now-10))
		pq.push(qEvent("wf", "new", "h2", now+10))
		got := pq.popLimited(map[string]uint32{"wf": 5}, expiration)
		require.NotNil(t, got)
		assert.Equal(t, "new", got.ExecutionID)
		assert.Equal(t, 0, pq.Len(), "expired event is dropped, not restored")
	})

	t.Run("workflow capacity limits and restores", func(t *testing.T) {
		t.Parallel()
		var pq priorityQueue
		pq.push(qEvent("full", "a", "h1", now+1)) // first by deadline, but no capacity
		pq.push(qEvent("ok", "b", "h2", now+2))
		caps := map[string]uint32{"ok": 1}
		got := pq.popLimited(caps, expiration)
		require.NotNil(t, got)
		assert.Equal(t, "b", got.ExecutionID)
		assert.Equal(t, uint32(0), caps["ok"], "capacity is decremented")
		require.Equal(t, 1, pq.Len(), "limited event is restored")
		assert.Equal(t, "a", pq.mustPop().ExecutionID)
	})

	t.Run("absent workflow has no capacity", func(t *testing.T) {
		t.Parallel()
		var pq priorityQueue
		pq.push(qEvent("unknown", "a", "h1", now+1))
		require.Nil(t, pq.popLimited(map[string]uint32{}, expiration))
		assert.Equal(t, 1, pq.Len())
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		var pq priorityQueue
		require.Nil(t, pq.popLimited(nil, expiration))
	})
}

func TestBlobHandleMap(t *testing.T) {
	t.Parallel()

	m := blobHandleMap{}
	m.add("h1", "b1", 3)
	m.add("h1", "b2", 4)
	m.addAll("h1", map[string]uint32{"b2": 9, "b3": 1})
	m.addAll("h2", nil)

	assert.Equal(t, map[string]uint32{"b1": 3, "b2": 9, "b3": 1}, m["h1"], "addAll overwrites and merges")
	assert.Contains(t, m, "h2")
	assert.Empty(t, m["h2"])
}

func TestVoterBlobs(t *testing.T) {
	t.Parallel()

	m := voterBlobs{}
	m.add("h", 1, "b1", 3)
	m.add("h", 2, "b1", 5) // shares a blob, and disagrees about where the event is
	m.add("h", 3, "b2", 9)
	m.add("h", 1, "b3", 4) // voter 1 votes again

	assert.Equal(t, map[uint8]blobRef{1: {"b3", 4}, 2: {"b1", 5}, 3: {"b2", 9}}, m["h"], "a voter has one blob, the latest")
	assert.Equal(t, map[string]uint32{"b1": 5, "b3": 4, "b2": 9}, m.handles("h"))
	assert.Empty(t, m.handles("unknown"))

	t.Run("handles picks the lowest index when voters disagree, whatever the order", func(t *testing.T) {
		for range 50 {
			m := voterBlobs{}
			for _, voter := range rand.Perm(10) {
				m.add("h", uint8(voter), "shared", uint32(100+voter))
			}
			require.Equal(t, map[string]uint32{"shared": 100}, m.handles("h"))
		}
	})
}

// A voter that votes again, round after round, must not make the blobs of an event grow: that would let one oracle
// fill the KV state, and have every other event pruned from it, and give the event a dead blob to try for each round.
func TestAddObservedEvents_OneBlobPerVoter(t *testing.T) {
	t.Parallel()
	const f = 2 // so that 2 or 3 votes stay in the observed set (it takes 2f+1=5 to promote)

	s := newKVState()
	s.observedEvents["h"] = obsEvent("wf", "e", 100, 1, 2)
	s.observedBlobs.add("h", 1, "blob-1", 0)
	s.observedBlobs.add("h", 2, "blob-2", 0)

	for round := range 200 {
		blobs := voterBlobs{}
		blobs.add("h", 2, fmt.Sprintf("fresh-blob-%d", round), uint32(round)) // node 2 keeps voting, with a new blob each time
		s.addObservedEvents(f, map[string]*pb.ConsensusQueueObservedEvent{"h": obsEvent("wf", "e", 100, 2)}, blobs, func(*pb.ConsensusQueueEvent) {})
	}

	require.Contains(t, s.observedEvents, "h")
	assert.Equal(t, []byte{1, 2}, s.observedEvents["h"].NodeIds)
	assert.Equal(t, map[uint8]blobRef{1: {"blob-1", 0}, 2: {"fresh-blob-199", 199}}, s.observedBlobs["h"], "still one blob each")

	// and so it is on disk
	s.prepareWriteObservedEvents()
	assert.Len(t, s.observedEventEntries[0].Event.Blobs, 2)
	assert.Len(t, s.observedBlobHandles, 2)
}

func TestPruneExpiredObservations(t *testing.T) {
	t.Parallel()

	expiration := time.Unix(100, 0)

	t.Run("all expired deletes entry", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		s.observedEvents["gone"] = obsEvent("wf", "e1", 10, 1, 2)
		s.pruneExpiredObservations(expiration)
		assert.Empty(t, s.observedEvents)
	})

	t.Run("none expired keeps entry intact", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		s.observedEvents["keep"] = obsEvent("wf", "e1", 200, 1, 2)
		s.pruneExpiredObservations(expiration)
		require.Contains(t, s.observedEvents, "keep")
		assert.Equal(t, []byte{1, 2}, s.observedEvents["keep"].NodeIds)
	})

	t.Run("partially expired removes only expired and keeps slices aligned", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		e := obsEvent("wf", "e1", 0, 1, 2, 3, 4)
		// nodes 1 & 3 expired, 2 & 4 live
		e.Deadlines = []*timestamppb.Timestamp{ts(50), ts(150), ts(60), ts(160)}
		s.observedEvents["h"] = e
		s.pruneExpiredObservations(expiration)

		got := s.observedEvents["h"]
		require.NotNil(t, got)
		assert.Equal(t, []byte{2, 4}, got.NodeIds)
		require.Len(t, got.Deadlines, 2)
		require.Len(t, got.ObservedAts, 2)
		assert.Equal(t, int64(150), got.Deadlines[0].Seconds)
		assert.Equal(t, int64(160), got.Deadlines[1].Seconds)
		assert.Equal(t, int64(1), got.ObservedAts[0].Seconds, "observedAts stay paired with node 2")
		assert.Equal(t, int64(3), got.ObservedAts[1].Seconds, "observedAts stay paired with node 4")
	})
}

func TestAddObservedEvents(t *testing.T) {
	t.Parallel()
	const f = 1

	collect := func(into *[]*pb.ConsensusQueueEvent) func(*pb.ConsensusQueueEvent) {
		return func(e *pb.ConsensusQueueEvent) { *into = append(*into, e) }
	}

	t.Run("fewer than f+1 is ignored", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		var queued []*pb.ConsensusQueueEvent
		s.addObservedEvents(f, map[string]*pb.ConsensusQueueObservedEvent{"h": obsEvent("wf", "e", 100, 1)}, voterBlobs{}, collect(&queued))
		assert.Empty(t, s.observedEvents)
		assert.Empty(t, s.queuedEvents)
		assert.Empty(t, queued)
	})

	t.Run("f+1 goes to observed set", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		var queued []*pb.ConsensusQueueEvent
		blobs := voterBlobs{}
		blobs.add("h", 1, "blob-1", 0)
		blobs.add("h", 2, "blob-2", 3)
		s.addObservedEvents(f, map[string]*pb.ConsensusQueueObservedEvent{"h": obsEvent("wf", "e", 100, 1, 2)}, blobs, collect(&queued))
		require.Contains(t, s.observedEvents, "h")
		assert.Empty(t, s.queuedEvents)
		assert.Empty(t, queued)
		assert.Equal(t, map[uint8]blobRef{1: {"blob-1", 0}, 2: {"blob-2", 3}}, s.observedBlobs["h"], "each voter's blob")
	})

	t.Run("2f+1 goes straight to the priority queue", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		var queued []*pb.ConsensusQueueEvent
		blobs := voterBlobs{}
		for voter := range uint8(3) {
			blobs.add("h", voter+1, "blob", 0) // the voters share a blob
		}
		s.addObservedEvents(f, map[string]*pb.ConsensusQueueObservedEvent{"h": obsEvent("wf", "e", 100, 1, 2, 3)}, blobs, collect(&queued))
		assert.Empty(t, s.observedEvents)
		require.Len(t, s.queuedEvents, 1)
		require.Len(t, queued, 1)
		assert.Equal(t, "h", string(queued[0].TimelessHash))
		assert.Equal(t, int64(102), queued[0].Deadline.Seconds, "median of 101,102,103")
		assert.Equal(t, map[string]uint32{"blob": 0}, s.queuedBlobs["h"])
	})

	t.Run("update promotes at 2f+1 and carries blobs", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		var queued []*pb.ConsensusQueueEvent
		// previously observed by nodes 1 & 2 (with their blobs)
		s.observedEvents["h"] = obsEvent("wf", "e", 100, 1, 2)
		s.observedBlobs.add("h", 1, "blob-1", 1)
		s.observedBlobs.add("h", 2, "old-blob-2", 2)

		newBlobs := voterBlobs{}
		newBlobs.add("h", 2, "new-blob-2", 5) // node 2 votes again, with a new blob
		newBlobs.add("h", 3, "blob-3", 7)
		// node 3 is new; node 2 re-observes with a new deadline
		toAdd := obsEvent("wf", "e", 100, 2, 3)
		toAdd.Deadlines[0] = ts(500)
		s.addObservedEvents(f, map[string]*pb.ConsensusQueueObservedEvent{"h": toAdd}, newBlobs, collect(&queued))

		assert.Empty(t, s.observedEvents, "promoted out of the observed set")
		assert.Empty(t, s.observedBlobs, "observed blobs moved")
		require.Len(t, queued, 1)
		assert.Equal(t, map[string]uint32{"blob-1": 1, "new-blob-2": 5, "blob-3": 7}, s.queuedBlobs["h"], "node 2's old blob is replaced, not kept")
	})

	t.Run("update below 2f+1 stays observed and does not panic", func(t *testing.T) {
		t.Parallel()
		s := newKVState()
		s.observedEvents["h"] = obsEvent("wf", "e", 100, 1, 2)
		toAdd := obsEvent("wf", "e", 200, 2) // existing node, overwrite
		require.NotPanics(t, func() {
			s.addObservedEvents(f, map[string]*pb.ConsensusQueueObservedEvent{"h": toAdd}, voterBlobs{}, func(*pb.ConsensusQueueEvent) {})
		})
		got := s.observedEvents["h"]
		require.NotNil(t, got)
		assert.Equal(t, []byte{1, 2}, got.NodeIds)
		assert.Equal(t, int64(202), got.Deadlines[1].Seconds, "node 2's deadline overwritten")
	})
}

func TestKVState_RoundTrip(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.observedEvents["obs"] = obsEvent("wf", "e-obs", 100, 1, 2)
	s.observedBlobs.add("obs", 1, "blobA", 0)
	s.observedBlobs.add("obs", 2, "blobC", 4)
	s.queuedEvents.push(qEvent("wf", "e-q", "q", 300))
	s.queuedBlobs.add("q", "blobB", 2)
	s.consensusEvents = append(s.consensusEvents, qEvent("wf", "e-c", "c", 200))
	s.queuedBlobs.add("c", "blobB", 7)

	// size* calls prepare the write; mirrors what StateTransition does.
	require.Positive(t, s.sizeOf())

	kv := memKV{}
	require.NoError(t, s.writeTo(&stateWriter{kv: kv}))

	var loaded kvState
	require.NoError(t, loaded.loadFrom(kv))
	require.Contains(t, loaded.observedEvents, "obs")
	assert.Equal(t, map[uint8]blobRef{1: {"blobA", 0}, 2: {"blobC", 4}}, loaded.observedBlobs["obs"], "each voter's blob survives a round trip")
	require.Len(t, loaded.queuedEvents, 1)
	assert.Equal(t, map[string]uint32{"blobB": 2}, loaded.queuedBlobs["q"])

	consensus, handles, err := readConsensusState(kv)
	require.NoError(t, err)
	require.Len(t, consensus, 1)
	require.Len(t, consensus[0].Blobs, 1)
	assert.Equal(t, "blobB", string(handles[consensus[0].Blobs[0].Blob]))
	assert.Equal(t, uint32(7), consensus[0].Blobs[0].Element)
}

// writeTo is called by StateTransition without any size check having run first when no events are popped/pruned.
// It must still persist observed events and blob handles.
func TestKVState_WriteToWithoutSizeCheck(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.observedEvents["obs"] = obsEvent("wf", "e-obs", 100, 1, 2)
	s.observedBlobs.add("obs", 1, "blobA", 0)
	s.observedBlobs.add("obs", 2, "blobA", 0)
	s.queuedEvents.push(qEvent("wf", "e-q", "q", 300))
	s.queuedBlobs.add("q", "blobB", 2)

	kv := memKV{}
	require.NoError(t, s.writeTo(&stateWriter{kv: kv}))

	var loaded kvState
	require.NoError(t, loaded.loadFrom(kv))
	assert.Contains(t, loaded.observedEvents, "obs", "observed events must be persisted")
	require.Len(t, loaded.queuedEvents, 1)
	assert.Equal(t, map[string]uint32{"blobB": 2}, loaded.queuedBlobs["q"], "queued event must keep its blob reference")
}

func TestPrepareWriteObservedEvents_Order(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.observedEvents["few"] = obsEvent("wf", "zzz", 100, 1, 2)
	s.observedEvents["many"] = obsEvent("wf", "b", 100, 1, 2, 3)
	s.observedEvents["tie-x"] = obsEvent("wf", "same", 100, 1, 2)
	s.observedEvents["tie-y"] = obsEvent("wf", "same", 100, 1, 2)

	order := func() (got []string) {
		s.prepareWriteObservedEvents()
		for _, e := range s.observedEventEntries {
			got = append(got, string(e.TimelessHash))
		}
		return
	}
	first := order()
	require.Equal(t, "many", first[0], "most endorsed first, so the last (pruned) is the least endorsed")
	require.Equal(t, []string{"few", "tie-x", "tie-y"}, first[1:], "ties on votes break on hash")
	for range 50 { // map iteration is random; ties on votes and execution ID must still order identically
		require.Equal(t, first, order())
	}
}

// errKV fails reads and/or writes of specific keys.
type errKV struct {
	memKV
	readErrKey, writeErrKey string
}

func (e *errKV) Read(k []byte) ([]byte, error) {
	if string(k) == e.readErrKey {
		return nil, errBoom
	}
	return e.memKV.Read(k)
}

func (e *errKV) Write(k, v []byte) error {
	if string(k) == e.writeErrKey {
		return errBoom
	}
	return e.memKV.Write(k, v)
}

var allKeys = []string{observedEventsKey, queuedEventsKey, consensusEventsKey, blobHandlesKey}

func TestPriorityQueue_EdgeCases(t *testing.T) {
	t.Parallel()

	var nilPQ *priorityQueue
	assert.Equal(t, 0, nilPQ.Len())

	var pq priorityQueue
	assert.Nil(t, pq.Pop())
	assert.Nil(t, pq.mustPop())
}

func TestQueueObservations(t *testing.T) {
	t.Parallel()

	s := newKVState()
	got := s.queueObservations("hash", obsEvent("wf", "exec", 100, 1, 2, 3))

	assert.Equal(t, "wf", got.WorkflowID)
	assert.Equal(t, "exec", got.ExecutionID)
	assert.Equal(t, []byte("hash"), got.TimelessHash)
	assert.Equal(t, int64(101), got.ObservedAt.Seconds, "median of 100,101,102")
	assert.Equal(t, int64(102), got.Deadline.Seconds, "median of 101,102,103")
	assert.Empty(t, got.Blobs, "blobs are tracked separately and attached on write")
	require.Len(t, s.queuedEvents, 1)
	assert.Same(t, got, s.queuedEvents[0])
}

func TestPopLimited_Sequence(t *testing.T) {
	t.Parallel()

	var pq priorityQueue
	for i, d := range []int64{400, 100, 300, 200} {
		pq.push(qEvent("a", string(rune('a'+i)), string(rune('a'+i)), d))
	}
	caps := map[string]uint32{"a": 2}
	past := time.Unix(0, 0)

	assert.Equal(t, int64(100), pq.popLimited(caps, past).Deadline.Seconds)
	assert.Equal(t, int64(200), pq.popLimited(caps, past).Deadline.Seconds)
	assert.Nil(t, pq.popLimited(caps, past), "workflow capacity exhausted")
	require.Equal(t, 2, pq.Len(), "limited events are restored")

	caps["a"] = 5
	assert.Equal(t, int64(300), pq.popLimited(caps, past).Deadline.Seconds, "restored events keep priority order")
	assert.Equal(t, int64(400), pq.popLimited(caps, past).Deadline.Seconds)
}

func TestReadFunctions(t *testing.T) {
	t.Parallel()

	t.Run("absent keys read as empty", func(t *testing.T) {
		t.Parallel()
		kv := memKV{}

		queued, err := readQueuedEvents(kv)
		require.NoError(t, err)
		assert.Empty(t, queued)

		observed, blobs, err := readObservedEvents(kv)
		require.NoError(t, err)
		assert.NotNil(t, observed, "must be usable as a map")
		assert.NotNil(t, blobs)
		assert.Empty(t, observed)

		consensus, handles, err := readConsensusState(kv)
		require.NoError(t, err)
		assert.Empty(t, consensus)
		assert.Empty(t, handles)

		handles, err = readBlobHandles(kv)
		require.NoError(t, err)
		assert.Empty(t, handles)
	})

	t.Run("corrupt values error", func(t *testing.T) {
		t.Parallel()
		for _, key := range allKeys {
			kv := memKV{key: {0xff, 0xff, 0xff}}
			var err error
			switch key {
			case queuedEventsKey:
				_, err = readQueuedEvents(kv)
			case observedEventsKey:
				_, _, err = readObservedEvents(kv)
			case consensusEventsKey:
				_, _, err = readConsensusState(kv)
			case blobHandlesKey:
				_, err = readBlobHandles(kv)
			}
			require.Error(t, err, key)
		}
	})

	t.Run("read errors are returned", func(t *testing.T) {
		t.Parallel()
		_, err := readQueuedEvents(&errKV{readErrKey: queuedEventsKey})
		require.ErrorIs(t, err, errBoom)
		_, _, err = readObservedEvents(&errKV{readErrKey: observedEventsKey})
		require.ErrorIs(t, err, errBoom)
		_, _, err = readConsensusState(&errKV{readErrKey: consensusEventsKey})
		require.ErrorIs(t, err, errBoom)
		_, _, err = readConsensusState(&errKV{readErrKey: blobHandlesKey}) // blob handles are read as part of consensus state
		require.ErrorIs(t, err, errBoom)
		_, err = readBlobHandles(&errKV{readErrKey: blobHandlesKey})
		require.ErrorIs(t, err, errBoom)
	})
}

func TestLoadFrom_Errors(t *testing.T) {
	t.Parallel()

	for _, key := range []string{observedEventsKey, queuedEventsKey, blobHandlesKey} {
		t.Run("read error "+key, func(t *testing.T) {
			t.Parallel()
			var s kvState
			require.ErrorIs(t, s.loadFrom(&errKV{memKV: memKV{}, readErrKey: key}), errBoom)
		})
		t.Run("corrupt "+key, func(t *testing.T) {
			t.Parallel()
			var s kvState
			require.Error(t, s.loadFrom(memKV{key: {0xff, 0xff, 0xff}}))
		})
	}

	t.Run("consensus events are not loaded", func(t *testing.T) {
		t.Parallel()
		var s kvState
		require.NoError(t, s.loadFrom(memKV{consensusEventsKey: {0xff, 0xff, 0xff}}))
		assert.Empty(t, s.consensusEvents)
	})
}

func TestWriteTo_Errors(t *testing.T) {
	t.Parallel()

	for _, key := range allKeys {
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			s := newKVState()
			err := s.writeTo(&stateWriter{kv: &errKV{memKV: memKV{}, writeErrKey: key}})
			require.ErrorIs(t, err, errBoom)
		})
	}
}

// sizeOf is used to stay within libocr's KV limits, so it must equal what is actually written.
func TestKVState_SizeOfMatchesWrittenBytes(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.observedEvents["o1"] = obsEvent("wf", "e-o1", 100, 1, 2)
	s.observedEvents["o2"] = obsEvent("wf", "e-o2", 100, 1, 2, 3)
	for _, voter := range []uint8{1, 2} {
		s.observedBlobs.add("o1", voter, "blobA", 0)
	}
	for _, voter := range []uint8{1, 2, 3} {
		s.observedBlobs.add("o2", voter, "blobA", 1)
	}
	s.queuedEvents.push(qEvent("wf", "e-q", "q", 300))
	s.queuedBlobs.add("q", "blobB", 2)
	s.consensusEvents = append(s.consensusEvents, qEvent("wf", "e-c", "c", 200))
	s.queuedBlobs.add("c", "blobB", 3)

	size := s.sizeOf()
	kv := memKV{}
	require.NoError(t, s.writeTo(&stateWriter{kv: kv}))

	var written int
	for k, v := range kv {
		written += len(k) + len(v)
	}
	assert.Equal(t, written, size)
}

func TestPrepareWriteBlobHandles(t *testing.T) {
	t.Parallel()

	// resolve returns an event's blob refs as handle -> element
	resolve := func(s *kvState, e *pb.ConsensusQueueEvent) map[string]uint32 {
		out := map[string]uint32{}
		for _, b := range e.Blobs {
			out[string(s.queuedBlobHandles[b.Blob])] = b.Element
		}
		return out
	}

	s := newKVState()
	s.queuedEvents.push(qEvent("wf", "e1", "h1", 10))
	s.queuedEvents.push(qEvent("wf", "e2", "h2", 20))
	s.consensusEvents = append(s.consensusEvents, qEvent("wf", "e3", "h3", 5))
	s.queuedBlobs.add("h1", "blobX", 0)
	s.queuedBlobs.add("h2", "blobX", 1)
	s.queuedBlobs.add("h3", "blobY", 0)
	s.queuedBlobs.add("h3", "blobX", 2)

	s.prepareWriteBlobHandles()

	assert.ElementsMatch(t, [][]byte{[]byte("blobX"), []byte("blobY")}, s.queuedBlobHandles, "shared handles are stored once")
	for _, e := range slices.Concat(s.queuedEvents, s.consensusEvents) {
		assert.Equal(t, s.queuedBlobs[string(e.TimelessHash)], resolve(s, e), e.ExecutionID)
	}

	t.Run("idempotent", func(t *testing.T) {
		handles := slices.Clone(s.queuedBlobHandles)
		s.prepareWriteBlobHandles()
		assert.Equal(t, handles, s.queuedBlobHandles)
		for _, e := range slices.Concat(s.queuedEvents, s.consensusEvents) {
			assert.Equal(t, s.queuedBlobs[string(e.TimelessHash)], resolve(s, e))
			assert.Len(t, e.Blobs, len(s.queuedBlobs[string(e.TimelessHash)]), "blobs are replaced, not appended")
		}
	})

	t.Run("event without blobs", func(t *testing.T) {
		s := newKVState()
		s.queuedEvents.push(qEvent("wf", "e1", "h1", 10))
		require.NotPanics(t, s.prepareWriteBlobHandles)
		assert.Empty(t, s.queuedBlobHandles)
		assert.Empty(t, s.queuedEvents[0].Blobs)
	})
}

func TestPrepareWriteObservedEvents_SharedBlobs(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.observedEvents["h1"] = obsEvent("wf", "e1", 100, 1, 2)
	s.observedEvents["h2"] = obsEvent("wf", "e2", 100, 1, 2)
	s.observedBlobs.add("h1", 1, "blobX", 0)
	s.observedBlobs.add("h1", 2, "blobX", 0)
	s.observedBlobs.add("h2", 1, "blobX", 1)
	s.observedBlobs.add("h2", 2, "blobY", 0) // two observers used different blobs

	s.prepareWriteObservedEvents()

	assert.ElementsMatch(t, [][]byte{[]byte("blobX"), []byte("blobY")}, s.observedBlobHandles, "a shared blob is stored once")
	require.Len(t, s.observedEventEntries, 2)
	for _, entry := range s.observedEventEntries {
		hash := string(entry.TimelessHash)
		require.Len(t, entry.Event.Blobs, len(entry.Event.NodeIds), "one blob per vote, aligned with NodeIds")
		for i, voter := range entry.Event.NodeIds {
			b := entry.Event.Blobs[i]
			assert.Equal(t, s.observedBlobs[hash][voter], blobRef{string(s.observedBlobHandles[b.Blob]), b.Element}, "event %s voter %d", hash, voter)
		}
	}
}

// pruneExpiredObservations and popLimited must agree on whether a deadline equal to the expiration is expired.
// popLimited keeps it (Deadline.Before(expiration) is false); pruning must not delete it either.
func TestPruneExpiredObservations_BoundaryMatchesPopLimited(t *testing.T) {
	t.Parallel()

	expiration := time.Unix(100, 0)

	var pq priorityQueue
	pq.push(qEvent("wf", "e", "h", 100))
	require.NotNil(t, pq.popLimited(map[string]uint32{"wf": 1}, expiration), "popLimited keeps an event whose deadline equals the expiration")

	s := newKVState()
	s.observedEvents["h"] = &pb.ConsensusQueueObservedEvent{
		WorkflowID: "wf", ExecutionID: "e", NodeIds: []byte{1},
		ObservedAts: []*timestamppb.Timestamp{ts(99)}, Deadlines: []*timestamppb.Timestamp{ts(100)}, // deadline == expiration
	}
	s.pruneExpiredObservations(expiration)
	assert.Contains(t, s.observedEvents, "h", "pruning must not delete an observation popLimited would still accept")
}

func TestAddObservedEvents_Order(t *testing.T) {
	t.Parallel()

	toAdd := map[string]*pb.ConsensusQueueObservedEvent{}
	for _, h := range []string{"d", "b", "e", "a", "c"} {
		toAdd[h] = obsEvent("wf", "exec-"+h, 100, 1, 2, 3)
	}
	for range 20 { // map iteration is random; callbacks must not depend on it
		s := newKVState()
		var got []string
		s.addObservedEvents(1, toAdd, voterBlobs{}, func(e *pb.ConsensusQueueEvent) { got = append(got, string(e.TimelessHash)) })
		require.Equal(t, []string{"a", "b", "c", "d", "e"}, got)
	}
}

// An event that already sits in the priority queue must not be promoted again if it is observed again
// (e.g. by a node that was slow to drop it from its local queue).
func TestAddObservedEvents_AlreadyQueued(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.queuedEvents.push(qEvent("wf", "exec", "h", 100))
	s.queuedBlobs.add("h", "blob", 0)

	var queued int
	s.addObservedEvents(1, map[string]*pb.ConsensusQueueObservedEvent{"h": obsEvent("wf", "exec", 100, 1, 2, 3)},
		voterBlobs{}, func(*pb.ConsensusQueueEvent) { queued++ })

	assert.Zero(t, queued, "already queued")
	assert.Len(t, s.queuedEvents, 1, "no duplicate in the priority queue")
}

func TestAddObservedEvents_UpdateDeduplicatesNodes(t *testing.T) {
	t.Parallel()

	s := newKVState()
	s.observedEvents["h"] = obsEvent("wf", "e", 100, 1, 2)
	s.addObservedEvents(2, map[string]*pb.ConsensusQueueObservedEvent{"h": obsEvent("wf", "e", 100, 1, 2)}, voterBlobs{}, func(*pb.ConsensusQueueEvent) {})
	assert.Equal(t, []byte{1, 2}, s.observedEvents["h"].NodeIds, "re-observation by the same nodes adds no votes")
	assert.Empty(t, s.queuedEvents)
}

// revPQ pops the latest deadline first.
type revPQ struct{ *priorityQueue }

func (r *revPQ) Less(i, j int) bool {
	a, b := (*r.priorityQueue)[i], (*r.priorityQueue)[j]
	return b.CompareDeadline(a) < 0
}

// naivePruneObservedToFit is the original algorithm: drop the last event and measure everything again, until it fits.
// (The original read the entries before they were refreshed, so it panicked when the rest of the state alone was over
// the limit and the last observed event was dropped. This one refreshes them first.)
func naivePruneObservedToFit(s *kvState, maxTotalBytes int) {
	for {
		s.prepareWriteObservedEvents()
		if len(s.observedEventEntries) == 0 || (s.sizeOfObservedEvents() <= maxKeyValueValueBytes && s.sizeOf() <= maxTotalBytes) {
			return
		}
		delete(s.observedEvents, string(s.observedEventEntries[len(s.observedEventEntries)-1].TimelessHash))
	}
}

// naivePruneQueuedToFit is the original algorithm: pop the lowest priority new event, drop it, and measure everything again.
func naivePruneQueuedToFit(s *kvState, newlyQueued []string, maxTotalBytes int) {
	fits := func() bool {
		return s.sizeOfQueuedEvents() <= maxKeyValueValueBytes && s.sizeOfBlobHandles() <= maxKeyValueValueBytes && s.sizeOf() <= maxTotalBytes
	}
	if len(newlyQueued) == 0 || fits() {
		return
	}
	set := map[string]struct{}{}
	for _, h := range newlyQueued {
		set[h] = struct{}{}
	}
	var eligible []*pb.ConsensusQueueEvent
	for _, e := range s.queuedEvents {
		if _, ok := set[string(e.TimelessHash)]; ok {
			eligible = append(eligible, e)
		}
	}
	backwards := &revPQ{priorityQueue: new(priorityQueue(eligible))}
	heap.Init(backwards)
	for backwards.Len() > 0 && !fits() {
		toRemove := heap.Pop(backwards).(*pb.ConsensusQueueEvent)
		s.queuedEvents = slices.DeleteFunc(s.queuedEvents, func(e *pb.ConsensusQueueEvent) bool { return e == toRemove })
	}
	heap.Init(&s.queuedEvents)
}

// randomState builds a state with observed and queued events sharing a pool of blob handles.
func randomState(seed int64) (s *kvState, newlyQueued []string) {
	rng := rand.New(rand.NewSource(seed))
	s = newKVState()
	pool := make([]string, 3+rng.Intn(12))
	for i := range pool {
		pool[i] = string(handleBytes(byte(i + 1)))
	}
	for i := range 5 + rng.Intn(120) {
		hash := fmt.Sprintf("observed-%03d", i)
		nodes := make([]byte, 2+rng.Intn(3))
		for v := range nodes {
			nodes[v] = byte(v)
		}
		s.observedEvents[hash] = obsEvent(hexID("wf", i%7), hexID("exec-o", i), int64(100+rng.Intn(1000)), nodes...)
		for _, node := range nodes {
			s.observedBlobs.add(hash, node, pool[rng.Intn(len(pool))], uint32(rng.Intn(1000)))
		}
	}
	for i := range 5 + rng.Intn(200) {
		hash := fmt.Sprintf("queued-%03d", i)
		s.queuedEvents.push(&pb.ConsensusQueueEvent{
			WorkflowID: hexID("wf", i%7), ExecutionID: hexID("exec-q", i), TimelessHash: []byte(hash),
			ObservedAt: ts(int64(rng.Intn(100))), Deadline: ts(int64(100 + rng.Intn(1000))),
		})
		for range 1 + rng.Intn(3) {
			s.queuedBlobs.add(hash, pool[rng.Intn(len(pool))], uint32(rng.Intn(1000)))
		}
		if rng.Intn(2) == 0 {
			newlyQueued = append(newlyQueued, hash)
		}
	}
	return s, newlyQueued
}

// The binary searches must drop exactly the events the one-at-a-time algorithms did, for any state and limit.
func TestPruneToFit_MatchesNaive(t *testing.T) {
	t.Parallel()

	for seed := range int64(12) {
		sizeProbe, _ := randomState(seed)
		full := sizeProbe.sizeOf()
		rng := rand.New(rand.NewSource(seed*7 + 1))
		limits := []int{1, full / 8, full/4 + rng.Intn(full), full + 1} // nothing fits, a lot dropped, some dropped, none dropped
		for _, maxTotal := range limits {
			naive, newlyQueued := randomState(seed)
			fast, _ := randomState(seed)

			naivePruneObservedToFit(naive, maxTotal)
			naivePruneQueuedToFit(naive, newlyQueued, maxTotal)
			fast.pruneObservedToFit(maxTotal)
			fast.pruneQueuedToFit(newlyQueued, maxTotal)

			require.ElementsMatch(t, slices.Collect(maps.Keys(naive.observedEvents)), slices.Collect(maps.Keys(fast.observedEvents)),
				"seed %d limit %d: observed events", seed, maxTotal)
			hashes := func(q priorityQueue) (out []string) {
				for _, e := range q {
					out = append(out, string(e.TimelessHash))
				}
				return
			}
			require.Equal(t, hashes(naive.queuedEvents), hashes(fast.queuedEvents), "seed %d limit %d: queued events (and heap order)", seed, maxTotal)
		}
	}
}

// popLimited is the original way to pop: one event at a time, putting back the skipped events after each.
func (p *priorityQueue) popLimited(capPerWF map[string]uint32, expiration time.Time) (event *pb.ConsensusQueueEvent) {
	var wfLimited []*pb.ConsensusQueueEvent // omitted due to workflow limit
	for {
		event = p.mustPop()
		if event == nil {
			break
		}

		if event.Deadline.AsTime().Before(expiration) {
			// drop expired
			continue
		}

		if capPerWF[event.WorkflowID] == 0 {
			// ignore but restore after
			wfLimited = append(wfLimited, event)
			continue
		}

		capPerWF[event.WorkflowID]--
		break
	}

	// restore omitted for next time
	for _, limited := range wfLimited {
		p.push(limited)
	}

	return
}

func randomQueue(rng *rand.Rand, n int, workflows int) (priorityQueue, map[string]uint32) {
	var pq priorityQueue
	caps := map[string]uint32{}
	for w := range workflows {
		if rng.Intn(3) != 0 { // some workflows have no capacity at all
			caps[hexID("wf", w)] = uint32(rng.Intn(6))
		}
	}
	for i := range n {
		pq.push(&pb.ConsensusQueueEvent{
			WorkflowID: hexID("wf", rng.Intn(workflows)), ExecutionID: hexID("exec", i), TimelessHash: []byte(fmt.Sprintf("h%05d", i)),
			ObservedAt: ts(1), Deadline: ts(int64(50 + rng.Intn(200))), // some are before the expiration
		})
	}
	return pq, caps
}

// Popping in a single pass must pop the same events, in the same order, and leave the same ones queued, as popping
// one at a time did.
func TestPopEligible_MatchesPoppingOneAtATime(t *testing.T) {
	t.Parallel()

	expiration := time.Unix(100, 0)
	for seed := range int64(100) {
		rng := rand.New(rand.NewSource(seed))
		n, workflows, want := rng.Intn(300), 1+rng.Intn(12), rng.Intn(60)

		oneByOne, capsA := randomQueue(rand.New(rand.NewSource(seed*3+1)), n, workflows)
		singlePass, capsB := randomQueue(rand.New(rand.NewSource(seed*3+1)), n, workflows)

		var wantPopped []string
		for len(wantPopped) < want {
			event := oneByOne.popLimited(capsA, expiration)
			if event == nil {
				break
			}
			wantPopped = append(wantPopped, event.ExecutionID)
		}
		var gotPopped []string
		for _, event := range singlePass.popEligible(want, capsB, expiration) {
			gotPopped = append(gotPopped, event.ExecutionID)
		}

		require.Equal(t, wantPopped, gotPopped, "seed %d: popped events", seed)
		require.Equal(t, capsA, capsB, "seed %d: workflow capacities left", seed)
		remaining := func(q priorityQueue) (out []string) {
			for _, e := range q {
				out = append(out, e.ExecutionID)
			}
			slices.Sort(out)
			return
		}
		require.Equal(t, remaining(oneByOne), remaining(singlePass), "seed %d: events left in the queue", seed)
	}
}

func TestPopEligible_Edges(t *testing.T) {
	t.Parallel()

	var pq priorityQueue
	pq.push(qEvent(testWF, "a", "h1", 200))
	assert.Empty(t, pq.popEligible(0, map[string]uint32{testWF: 5}, time.Unix(0, 0)), "nothing wanted")
	assert.Empty(t, pq.popEligible(-3, map[string]uint32{testWF: 5}, time.Unix(0, 0)), "a negative amount is nothing")
	assert.Equal(t, 1, pq.Len())
	assert.Empty(t, (&priorityQueue{}).popEligible(5, nil, time.Unix(0, 0)), "an empty queue")
}

// Events waiting on a workflow without capacity must not slow down popping the rest: it was quadratic in them, with
// 10000 of them taking 1.9s to pop 1000 events for other workflows.
func BenchmarkPopEligible_StuckEvents(b *testing.B) {
	for _, stuck := range []int{0, 1000, 10000} {
		fill := func() priorityQueue {
			var pq priorityQueue
			for i := range stuck {
				pq.push(&pb.ConsensusQueueEvent{WorkflowID: hexID("stuck", 0), ExecutionID: hexID("s", i), TimelessHash: []byte(fmt.Sprintf("s%06d", i)), Deadline: ts(int64(100 + i))})
			}
			for i := range 1000 {
				pq.push(&pb.ConsensusQueueEvent{WorkflowID: hexID("ok", 0), ExecutionID: hexID("o", i), TimelessHash: []byte(fmt.Sprintf("o%06d", i)), Deadline: ts(int64(100000 + i))})
			}
			return pq
		}
		b.Run(fmt.Sprintf("single pass/%d stuck", stuck), func(b *testing.B) {
			for range b.N {
				b.StopTimer()
				pq := fill()
				b.StartTimer()
				pq.popEligible(1000, map[string]uint32{hexID("ok", 0): 1000}, time.Unix(0, 0))
			}
		})
		if stuck <= 1000 { // one at a time is too slow to wait for beyond this
			b.Run(fmt.Sprintf("one at a time/%d stuck", stuck), func(b *testing.B) {
				for range b.N {
					b.StopTimer()
					pq := fill()
					caps := map[string]uint32{hexID("ok", 0): 1000}
					b.StartTimer()
					for range 1000 {
						pq.popLimited(caps, time.Unix(0, 0))
					}
				}
			})
		}
	}
}
