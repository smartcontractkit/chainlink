package triggers

import (
	"context"
	"errors"
)

var (
	// ErrQueueFull is returned by Put when the queue is at capacity.
	ErrQueueFull = errors.New("trigger queue full")

	// ErrSequenceUnavailable is returned by Put when the SequenceSource
	// cannot issue a value. The event is dropped.
	ErrSequenceUnavailable = errors.New("trigger queue sequence counter unavailable")
)

// Queue is the node's staging area for trigger events. It is a buffer only:
// it does not sort, dedupe, expire, or admit. Order is FIFO by SequenceNumber.
type Queue interface {
	// Put never blocks. It re-reads the capacity limit on every call, stamps
	// a monotonic SequenceNumber, and returns ErrQueueFull,
	// or ErrSequenceUnavailable instead of blocking.
	Put(context.Context, CoordinatedEvent) error

	// Observe atomically returns all buffered events in FIFO order
	// (ascending SequenceNumber) and empties the queue. The caller owns the
	// returned events.
	Observe(context.Context) ([]CoordinatedEvent, error)

	// Stats reports current depth, capacity, and cumulative drops.
	Stats(context.Context) (Stats, error)
}

// Stats is a point-in-time view of a Queue.
type Stats struct {
	Depth              uint64
	Capacity           uint64
	DepthByWorkflow    map[string]uint64 // keyed by WorkflowID
	DepthByTriggerType map[string]uint64 // keyed by TriggerCapID

	// DroppedByReason counts cumulative drops, keyed by reason (full,
	// sequence_unavailable) and then by the dropped event's trigger
	// registration ID (RegistrationID(WorkflowID, TriggerIndex)). Sum the inner
	// map for the per-reason total.
	DroppedByReason map[string]map[string]uint64
}

// SequenceSource issues strictly increasing values.
type SequenceSource interface {
	// Next returns the next strictly increasing value, or an error if none
	// can be issued.
	Next(context.Context) (uint64, error)
}
