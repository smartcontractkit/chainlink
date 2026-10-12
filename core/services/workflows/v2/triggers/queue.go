package triggers

import (
	"context"
	"errors"
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
	// monitoring.TriggerDropReasonCentralQueue* constants) and then by the
	// dropped event's trigger registration ID
	// (RegistrationID(WorkflowID, TriggerIndex)). Sum the inner map for the
	// per-reason total.
	DroppedByReason map[string]map[string]uint64
}
