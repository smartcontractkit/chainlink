package queue

import (
	"cmp"
	"container/heap"

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
