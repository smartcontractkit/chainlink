package trigger

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

// EventSink is how trigger events are delivered to an engine for execution.
type EventSink interface {
	ExecuteTrigger(ctx context.Context, event RoutedTriggerEvent) error
}

// Acknowledger acknowledges a trigger event without the engine owning the
// trigger handle. It is injected into EngineConfig so the engine's ACK
// call sites are decoupled from who holds the handles.
type Acknowledger interface {
	Ack(ctx context.Context, triggerCapID, triggerRegistrationID, eventID string) error
}

// RoutedTriggerEvent is the canonical trigger event type that flows
// through the dispatch path into the engine.
type RoutedTriggerEvent struct {
	WorkflowID   string
	TriggerCapID string
	TriggerIndex int
	ObservedAt   time.Time
	// Deadline     time.Time TODO: will be addressed on CRE-6175
	Lamport uint64 // reserved, always 0 in M1
	Event   capabilities.TriggerResponse
}
