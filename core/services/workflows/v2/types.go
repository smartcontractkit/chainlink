package v2

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
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
//     (not shown in the current code path; reserved for M2 coordinator).
//
// Ack is idempotent: calling it multiple times for the same event is safe.
// The implementation is responsible for looking up the trigger handle by
// triggerRegistrationID and calling AckEvent on it.
type Acknowledger interface {
	Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error
}

// Subscriber is how a caller obtains an engine's trigger subscriptions on
// demand. Subscribe issues the WASM Subscribe call directly (no caching): the
// engine holds no subscription state of its own, so every call is a fresh
// WASM round trip and callers are responsible for calling it exactly once
// per registration. Tenant identifies the tenant the subscriptions belong to.
type Subscriber interface {
	Subscribe(ctx context.Context) ([]*sdkpb.TriggerSubscription, error)
	Tenant() contexts.CRE
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
	Subscriber

	// IsCoordinated is true if the engine does not manage its own
	// trigger registration, trigger dequeuing, execution or acknowledgement.
	// Fixed at construction.
	IsCoordinated() bool
}
