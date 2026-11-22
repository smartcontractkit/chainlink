package triggers

import (
	"bytes"
	"cmp"
	"container/heap"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"iter"
	"maps"
	"slices"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
)

// ErrQueueFull is returned by Put when the queue is at capacity.
var ErrQueueFull = errors.New("trigger queue full")

// Queue is the node's staging area for trigger events. It is a buffer only:
// it does not sort, dedupe, expire, admit, or stamp events. Order is FIFO by
// insertion.
type Queue interface {
	// Put never blocks. It re-reads the capacity limit on every call and
	// returns ErrQueueFull instead of blocking.
	Put(context.Context, CoordinatedEvent) error

	// Drain atomically returns all buffered events in FIFO (insertion) order
	// and empties the queue. The caller owns the returned events.
	Drain(context.Context) ([]CoordinatedEvent, error)

	// Stats reports current length, capacity, and cumulative drops.
	Stats(context.Context) (Stats, error)

	// Sorted returns a sorted iterator over events, up to limit. Events remain but can be removed via Remove.
	// Events are filtered as they are pulled: those whose encoded size (without timestamps, i.e. as an element of a
	// blob) exceeds maxEventBytes can never be observed, so they are skipped (not counting toward limit), and dropped
	// from the queue (counted and logged) when iteration ends. maxEventBytes <= 0 disables this.
	Sorted(ctx context.Context, limit, maxEventBytes int) iter.Seq2[[]byte, *pb.CoordinatedEvent]

	// Remove removes events from the queue by execution ID.
	//TODO can't ack with trigger event ID - not unique!
	Remove(ctx context.Context, triggerEventIDs iter.Seq[string])
}

// Stats is a point-in-time view of a Queue.
type Stats struct {
	// Len is the number of events in the queue
	Len uint64

	// Capacity is the current limit or max number of events the queue can hold
	Capacity           uint64
	DepthByWorkflow    map[string]uint64 // keyed by WorkflowID
	DepthByTriggerType map[string]uint64 // keyed by TriggerCapID

	// DroppedByReason counts cumulative drops, keyed by reason (see the
	// monitoring.CentralTriggerQueueDropReason* constants). It is not keyed by
	// workflow or registration, so it stays bounded as workflows come and go.
	DroppedByReason map[string]uint64
}

// CentralQueue is the Queue implementation. It runs no goroutines. All state is guarded by mu, and gauges are
// recorded under mu so they are never written out of order.
type CentralQueue struct {
	limit   limits.BoundLimiter[int]
	metrics *monitoring.CentralTriggerQueueMetricLabeler
	lggr    logger.Logger

	mu                 sync.Mutex
	events             []centralQueueEvent
	depthByWorkflow    map[string]uint64
	depthByTriggerType map[string]uint64
	droppedByReason    map[string]uint64
}

type centralQueueEvent struct {
	CoordinatedEvent
	pb           *pb.CoordinatedEvent
	timelessHash []byte
	executionID  string
	blobSize     int // encoded size of pb without timestamps: its size as an element of a blob
}

func (q *centralQueueEvent) Compare(other centralQueueEvent) int {
	return cmp.Or(q.Deadline.Compare(other.Deadline), bytes.Compare(q.timelessHash, other.timelessHash))
}

// NewQueue returns the node's central trigger queue. limit is read on every Put. If it cannot be read, the
// CentralTriggerEventQueueLimit default is used. Metrics are recorded as part of each Put and Drain. The queue
// is safe for concurrent use.
func NewQueue(limit limits.BoundLimiter[int], metrics *monitoring.CentralTriggerQueueMetricLabeler, lggr logger.Logger) *CentralQueue {
	return &CentralQueue{
		limit:              limit,
		metrics:            metrics,
		lggr:               logger.Named(lggr, "CentralTriggerQueue"),
		depthByWorkflow:    map[string]uint64{},
		depthByTriggerType: map[string]uint64{},
		droppedByReason:    map[string]uint64{},
	}
}

// Put appends e, or returns ErrQueueFull and counts the drop when the queue is at capacity.
func (q *CentralQueue) Put(ctx context.Context, e CoordinatedEvent) error {
	capacity := q.capacity(ctx)

	q.mu.Lock()
	defer q.mu.Unlock()
	effective := q.effectiveCapacity(capacity)
	q.metrics.RecordCapacity(ctx, effective)
	if len(q.events) >= effective {
		q.drop(ctx, monitoring.CentralTriggerQueueDropReasonQueueFull)
		return ErrQueueFull
	}
	pbe := e.ToProto()
	timelessBytes, err := pbe.MarshalTimeless() // as it appears in a blob: hashed, and sized for the oversize filter
	if err != nil {
		return fmt.Errorf("failed to marshal event without timestamps: %w", err)
	}
	timelessHash := sha256.Sum256(timelessBytes)
	execID, err := workflows.GenerateExecutionIDWithTriggerIndex(e.WorkflowID, e.Event.Event.ID, e.TriggerIndex)
	if err != nil {
		return fmt.Errorf("failed to generate execution ID: %w", err)
	}
	q.events = append(q.events, centralQueueEvent{
		CoordinatedEvent: e,
		pb:               pbe,
		timelessHash:     timelessHash[:],
		executionID:      execID,
		blobSize:         len(timelessBytes),
	})
	q.depthByWorkflow[e.WorkflowID]++
	q.depthByTriggerType[e.TriggerCapID]++
	q.metrics.IncrementPutsCounter(ctx)
	q.metrics.RecordDepth(ctx, len(q.events))
	return nil
}

// Drain hands the buffered events to the caller and starts a new buffer, so later puts never write into the
// returned slice. It also records capacity, because Drain runs on a fixed interval even when nothing is put, so
// the gauge follows a changed limit on an idle queue.
func (q *CentralQueue) Drain(ctx context.Context) ([]CoordinatedEvent, error) {
	capacity := q.capacity(ctx)

	q.mu.Lock()
	defer q.mu.Unlock()
	events := make([]CoordinatedEvent, len(q.events))
	for i, pbe := range q.events {
		events[i] = pbe.CoordinatedEvent
	}
	q.events = nil
	q.depthByWorkflow = map[string]uint64{}
	q.depthByTriggerType = map[string]uint64{}
	q.metrics.IncrementDrainedEventsCounter(ctx, len(events))
	q.metrics.RecordDepth(ctx, 0)
	q.metrics.RecordCapacity(ctx, q.effectiveCapacity(capacity))
	return events, nil
}

func (q *CentralQueue) Remove(ctx context.Context, executionIDs iter.Seq[string]) {
	execIDSet := map[string]struct{}{}
	for triggerEventID := range executionIDs {
		execIDSet[triggerEventID] = struct{}{}
	}

	capacity := q.capacity(ctx)

	q.mu.Lock()
	defer q.mu.Unlock()

	var deleted int
	//TODO opt: map so we don't have to iterate over them all
	q.events = slices.DeleteFunc(q.events, func(cqe centralQueueEvent) bool {
		_, ok := execIDSet[cqe.executionID]
		if ok {
			deleted++
			q.depthByWorkflow[cqe.WorkflowID]--
			q.depthByTriggerType[cqe.TriggerCapID]--
			//TODO emit metric for time.Since(e.ObservedAt)
		}
		return ok
	})
	q.metrics.IncrementDrainedEventsCounter(ctx, deleted)
	q.metrics.RecordDepth(ctx, len(q.events))
	q.metrics.RecordCapacity(ctx, q.effectiveCapacity(capacity))
	return
}

func (q *CentralQueue) Sorted(ctx context.Context, limit, maxEventBytes int) iter.Seq2[[]byte, *pb.CoordinatedEvent] {
	q.mu.Lock()
	events := priorityQueue(slices.Clone(q.events))
	q.mu.Unlock()

	heap.Init(&events)

	return func(yield func([]byte, *pb.CoordinatedEvent) bool) {
		// Events are only examined as they are pulled. Those too large to ever be observed are skipped (they do not
		// count toward limit), and then dropped from the queue once the caller is done pulling.
		var oversize []centralQueueEvent
		defer func() { q.dropOversize(ctx, oversize, maxEventBytes) }()

		for count := 0; events.Len() > 0 && count < limit; {
			cqe := heap.Pop(&events).(centralQueueEvent)
			if maxEventBytes > 0 && cqe.blobSize > maxEventBytes {
				oversize = append(oversize, cqe)
				continue
			}
			count++
			if !yield(cqe.timelessHash, proto.CloneOf(cqe.pb)) {
				return
			}
		}
	}
}

// dropOversize removes the events from the queue, if they are still there.
// TODO should we fall back to executing instead of just dropping?
func (q *CentralQueue) dropOversize(ctx context.Context, oversize []centralQueueEvent, maxEventBytes int) {
	if len(oversize) == 0 {
		return
	}
	executionIDs := make(map[string]struct{}, len(oversize))
	for _, cqe := range oversize {
		executionIDs[cqe.executionID] = struct{}{}
	}

	q.mu.Lock()
	defer q.mu.Unlock()
	var dropped int
	q.events = slices.DeleteFunc(q.events, func(cqe centralQueueEvent) bool {
		if _, ok := executionIDs[cqe.executionID]; !ok {
			return false
		}
		dropped++
		q.depthByWorkflow[cqe.WorkflowID]--
		q.depthByTriggerType[cqe.TriggerCapID]--
		q.lggr.Errorw("Dropping event too large to be observed", "workflowID", cqe.WorkflowID, "executionID", cqe.executionID,
			"size", cqe.blobSize, "maxSize", maxEventBytes)
		q.drop(ctx, monitoring.CentralTriggerQueueDropReasonOversize)
		return true
	})
	if dropped > 0 {
		q.metrics.RecordDepth(ctx, len(q.events))
	}
}

// Stats returns copies of the queue's maps, so callers may keep or change them.
func (q *CentralQueue) Stats(ctx context.Context) (Stats, error) {
	capacity := q.capacity(ctx)

	q.mu.Lock()
	defer q.mu.Unlock()
	return Stats{
		Len:                uint64(len(q.events)),
		Capacity:           uint64(q.effectiveCapacity(capacity)), //nolint:gosec // never negative
		DepthByWorkflow:    maps.Clone(q.depthByWorkflow),
		DepthByTriggerType: maps.Clone(q.depthByTriggerType),
		DroppedByReason:    maps.Clone(q.droppedByReason),
	}, nil
}

// capacity reads the limit, falling back to the default when the read fails.
func (q *CentralQueue) capacity(ctx context.Context) queueCapacity {
	limit, err := q.limit.Limit(ctx)
	if err != nil {
		q.lggr.Warnw("Failed to read the central trigger queue limit, using the default", "err", err)
		return queueCapacity{limit: cresettings.Default.CentralTriggerEventQueueLimit.DefaultValue}
	}
	return queueCapacity{limit: limit, ok: true}
}

// effectiveCapacity is the capacity to enforce and report. With the default in use it never reads below the
// current length, so a node holding more than the default after a failed read reports full rather than
// over capacity. q.mu must be held.
func (q *CentralQueue) effectiveCapacity(c queueCapacity) int {
	limit := max(c.limit, 0)
	if !c.ok {
		limit = max(limit, len(q.events))
	}
	return limit
}

// drop counts one dropped event. q.mu must be held.
func (q *CentralQueue) drop(ctx context.Context, reason string) {
	q.metrics.IncrementDroppedEventsCounter(ctx, reason, 1)
	q.droppedByReason[reason]++
}

// queueCapacity is a limit read. ok is false when the read failed and limit holds the default.
type queueCapacity struct {
	limit int
	ok    bool
}

var _ heap.Interface = (*priorityQueue)(nil)

// priorityQueue is a heap of coordinated events, ordered by deadline and then trigger event ID.
type priorityQueue []centralQueueEvent

func (p *priorityQueue) Len() int {
	if p == nil {
		return 0
	}
	return len(*p)
}

func (p *priorityQueue) Less(i, j int) bool {
	a, b := (*p)[i], (*p)[j]
	return a.Compare(b) < 0
}

func (p *priorityQueue) Swap(i, j int) {
	(*p)[i], (*p)[j] = (*p)[j], (*p)[i]
}

func (p *priorityQueue) Push(x any) {
	*p = append(*p, x.(centralQueueEvent))
}

func (p *priorityQueue) Pop() any {
	if p.Len() == 0 {
		return nil
	}
	e := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	return e
}
