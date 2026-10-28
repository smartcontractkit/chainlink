package triggers

import (
	"context"
	"errors"
	"maps"
	"sync"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
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

// centralQueue is the Queue implementation. It runs no goroutines. All state is guarded by mu, and gauges are
// recorded under mu so they are never written out of order.
type centralQueue struct {
	limit   limits.BoundLimiter[int]
	metrics *monitoring.CentralTriggerQueueMetricLabeler
	lggr    logger.Logger

	mu                 sync.Mutex
	events             []CoordinatedEvent
	depthByWorkflow    map[string]uint64
	depthByTriggerType map[string]uint64
	droppedByReason    map[string]uint64
}

// NewQueue returns the node's central trigger queue. limit is read on every Put. If it cannot be read, the
// CentralTriggerEventQueueLimit default is used. Metrics are recorded as part of each Put and Drain. The queue
// is safe for concurrent use.
func NewQueue(limit limits.BoundLimiter[int], metrics *monitoring.CentralTriggerQueueMetricLabeler, lggr logger.Logger) *centralQueue {
	return &centralQueue{
		limit:              limit,
		metrics:            metrics,
		lggr:               logger.Named(lggr, "CentralTriggerQueue"),
		depthByWorkflow:    map[string]uint64{},
		depthByTriggerType: map[string]uint64{},
		droppedByReason:    map[string]uint64{},
	}
}

// Put appends e, or returns ErrQueueFull and counts the drop when the queue is at capacity.
func (q *centralQueue) Put(ctx context.Context, e CoordinatedEvent) error {
	capacity := q.capacity(ctx)

	q.mu.Lock()
	defer q.mu.Unlock()
	effective := q.effectiveCapacity(capacity)
	q.metrics.RecordCapacity(ctx, effective)
	if len(q.events) >= effective {
		q.drop(ctx, monitoring.CentralTriggerQueueDropReasonQueueFull)
		return ErrQueueFull
	}
	q.events = append(q.events, e)
	q.depthByWorkflow[e.WorkflowID]++
	q.depthByTriggerType[e.TriggerCapID]++
	q.metrics.IncrementPutsCounter(ctx)
	q.metrics.RecordDepth(ctx, len(q.events))
	return nil
}

// Drain hands the buffered events to the caller and starts a new buffer, so later puts never write into the
// returned slice.
func (q *centralQueue) Drain(ctx context.Context) ([]CoordinatedEvent, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	events := q.events
	q.events = nil
	q.depthByWorkflow = map[string]uint64{}
	q.depthByTriggerType = map[string]uint64{}
	q.metrics.IncrementDrainedEventsCounter(ctx, len(events))
	q.metrics.RecordDepth(ctx, 0)
	return events, nil
}

// Stats returns copies of the queue's maps, so callers may keep or change them.
func (q *centralQueue) Stats(ctx context.Context) (Stats, error) {
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
func (q *centralQueue) capacity(ctx context.Context) queueCapacity {
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
func (q *centralQueue) effectiveCapacity(c queueCapacity) int {
	limit := max(c.limit, 0)
	if !c.ok {
		limit = max(limit, len(q.events))
	}
	return limit
}

// drop counts one dropped event. q.mu must be held.
func (q *centralQueue) drop(ctx context.Context, reason string) {
	q.metrics.IncrementDroppedEventsCounter(ctx, reason, 1)
	q.droppedByReason[reason]++
}

// queueCapacity is a limit read. ok is false when the read failed and limit holds the default.
type queueCapacity struct {
	limit int
	ok    bool
}
