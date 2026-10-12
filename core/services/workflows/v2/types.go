package v2

import (
	"time"

	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

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

	triggers.EventSink
	triggers.Subscriber

	Drainable

	// IsCoordinated is true if the engine does not manage its own
	// trigger registration, trigger dequeuing, execution or acknowledgement.
	// Fixed at construction.
	IsCoordinated() bool
}
