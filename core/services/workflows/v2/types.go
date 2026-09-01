package v2

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

// EventSink is how trigger events are delivered to an engine for execution.
// ExecuteTrigger performs no admission — it runs the
// execution-intrinsic gates (dedup, shard ownership, metering) and then
// executes the workflow.
//
// Lifecycle hooks invoked during execution (in order):
//   - OnExecutionFinished(executionID, status) — always called, exactly once,
//     via defer. Status is one of: "completed", "errored", "timeout".
//   - OnExecutionError(msg) — called after OnExecutionFinished if the
//     execution returned an error (WASM failure, user error, or timeout).
//   - OnResultReceived(result) — called only on successful completion,
//     after OnExecutionFinished.
//
// Expected errors:
//   - ErrDuplicateExecution — the event was already executed (dedup gate).
//     The engine ACKs the duplicate internally before returning.
//   - ErrShardDeniedNotOwner — this node is not the shard owner.
//     The engine ACKs the event internally before returning.
//   - ErrShardDeniedOrchestrator — shard ownership check failed due to
//     orchestrator error. The engine ACKs the event internally before returning.
//   - ErrMeteringReserveFailed — metering report reservation failed.
//     No ACK is sent; the caller may retry.
//
// WASM execution errors (module failure, user workflow error, timeout) are
// NOT returned as errors. They are captured by OnExecutionError and
// OnExecutionFinished hooks. ExecuteTrigger returns nil in these cases.
type EventSink interface {
	ExecuteTrigger(ctx context.Context, event RoutedTriggerEvent) error
}

// Acknowledger acknowledges a trigger event without the engine owning the
// trigger handle. It is injected into EngineConfig so the engine's ACK
// call sites are decoupled from who holds the handles.
//
// The engine calls Ack in three situations:
//   - Duplicate execution — the event was already executed; the engine
//     re-ACKs to prevent redelivery.
//   - Shard ownership denial — this node is not the shard owner; the engine
//     ACKs to signal the event was processed (skipped).
//   - Normal execution start — the engine ACKs after the execution begins
//     (not shown in the current code path; reserved for M2 dispatcher).
//
// Ack is idempotent: calling it multiple times for the same event is safe.
// The implementation is responsible for looking up the trigger handle by
// triggerRegistrationID and calling AckEvent on it.
type Acknowledger interface {
	Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error
}

// RoutedTriggerEvent is the canonical trigger event type that flows
// through the dispatch path into the engine.
type RoutedTriggerEvent struct {
	WorkflowID   string
	TriggerCapID string
	TriggerIndex int

	// ObservedAt is the time the RoutedTriggerEvent was constructed by the
	// dispatcher. It is used for skew metrics (queue wait time)
	// and deadline enforcement.
	ObservedAt time.Time

	// Deadline is the expiry of this event in the dispatch queue,
	// stamped once at dispatch as ObservedAt + TriggerEventQueueTimeout.
	Deadline time.Time

	// SequenceNumber determines the execution order of trigger events across the DON. In M1 it is always 0 (no consensus ordering).
	SequenceNumber uint64
	Event          capabilities.TriggerResponse
}
