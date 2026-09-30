package v2

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
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
//   - ErrMeteringReserveFailed — metering report reservation failed.
//     No ACK is sent; the caller may retry.
//
// WASM execution errors (module failure, user workflow error, timeout) are
// NOT returned as errors. They are captured by OnExecutionError and
// OnExecutionFinished hooks. ExecuteTrigger returns nil in these cases.
type EventSink interface {
	ExecuteTrigger(ctx context.Context, event triggers.CoordinatedEvent) error
}

// Drainable is the graceful-shutdown contract. The syncer has a structurally
// identical local interface (syncer/v2.DrainableService); both are satisfied by
// the same methods, so no cross-package dependency is introduced.
type Drainable interface {
	Drain() bool
	ActiveExecutions() int32
	DrainStartedAt() (time.Time, bool)
}

// WorkflowEngine is the contract every workflow engine implementation satisfies,
// independent of which component owns trigger registration and acknowledgement.
type WorkflowEngine interface {
	services.Service
	EventSink
	Drainable
	triggers.Subscriber

	// IsCoordinated is true if the engine does not manage its own
	// trigger registration, trigger dequeuing, execution or acknowledgement.
	// Fixed at construction.
	IsCoordinated() bool
}
