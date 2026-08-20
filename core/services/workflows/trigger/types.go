package trigger

import (
	"context"
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
)

type TriggerEventSink interface {
	HandleTriggerEvent(ctx context.Context, event RoutedTriggerEvent) error
}

type RoutedTriggerEvent struct {
	WorkflowID   string
	TriggerCapID string
	TriggerIndex int
	ObservedAt   time.Time
	// Deadline     time.Time TODO: will be addressed on CRE-6175
	Lamport uint64 // reserved, always 0 in M1
	Event   capabilities.TriggerResponse
}
