package events

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	eventsv2 "github.com/smartcontractkit/chainlink-protos/workflows/go/v2"
)

// Known-answer fault injection (PoC).
//
// A DropDecider is consulted at the emitRawMessage fan-out seam — above BOTH
// transports (beholder/chip-ingress and the durable emitter) — so a drop
// suppresses the event everywhere and no retry layer below the seam can
// resurrect it. Every drop is recorded (structured log + counter) so an
// external detector can reconcile injected-vs-detected and measure recall.
//
// The drop decision is a pure function of (seed, workflowExecutionID). Because
// execution IDs are content-derived and identical on every node of a DON, all
// nodes independently drop the same executions with zero coordination, which
// defeats Byzantine-quorum masking downstream.

// DroppedWorkflowEventMsg is the exact log message emitted for every injected
// drop. It is a CONTRACT with external detectors that parse node logs — do not
// change the message or the field keys (executionID, eventType, workflowOwner)
// without updating the detector.
const DroppedWorkflowEventMsg = "Fault injection: dropped workflow event"

// droppedTotalMetricName is the OTEL counter (exported to prometheus via the
// beholder pipeline) incremented once per injected drop, labeled by event_type.
const droppedTotalMetricName = "workflow_fault_injection_dropped_total"

// Drop levels select which event types are eligible for dropping once an
// execution is selected by the deterministic hash.
const (
	// DropLevelFinishedOnly drops only WorkflowExecutionFinished events: the
	// execution stays visible downstream but never terminates, landing it in
	// the prober's non_terminal bucket (verifies the prober).
	DropLevelFinishedOnly = 1
	// DropLevelAllLifecycle drops all workflow execution lifecycle events: the
	// execution becomes invisible downstream, detectable only by reconciling
	// injection records (verifies the log backstop detector).
	DropLevelAllLifecycle = 2
)

// lifecycleEventTypes is the set of event types eligible for dropping at
// DropLevelAllLifecycle. Deliberately excludes MeteringReport (billing),
// deployment events (WorkflowStatusChanged/Activated/Paused/Deleted/
// ActivationAbandoned) and user logs/metrics.
var lifecycleEventTypes = map[string]struct{}{
	WorkflowExecutionStarted:    {},
	WorkflowExecutionFinished:   {},
	TriggerExecutionStarted:     {},
	CapabilityExecutionStarted:  {},
	CapabilityExecutionFinished: {},
	WorkflowExecutionProfile:    {},
}

// DropIdentity is the workflow identity known at the emit seam. Both fields
// must be non-empty for an event to ever be dropped (fail-closed).
type DropIdentity struct {
	// WorkflowOwner is the workflow owner address — the allowlist key. Always
	// present and stable, unlike orgID (best-effort, v2 engine only).
	WorkflowOwner string
	// WorkflowExecutionID is deterministic across the DON:
	// sha256(workflowID ‖ triggerEventID ‖ triggerIndex).
	WorkflowExecutionID string
}

// DropDecider decides whether a workflow event should be suppressed at the
// emit fan-out seam, before it reaches any transport. Implementations MUST be
// fail-closed (any missing identity => false) and MUST record every drop.
type DropDecider interface {
	// ShouldDrop returns true when the event must be suppressed. eventType is
	// the beholder entity, e.g. "workflows.v2.WorkflowExecutionFinished".
	ShouldDrop(ctx context.Context, id DropIdentity, eventType string) bool
}

type deciderHolder struct{ d DropDecider }

// globalDropDecider defaults to unset: the seam is a no-op passthrough until
// SetDropDecider installs a decider (mirrors durableemitter.SetGlobalEmitter).
var globalDropDecider atomic.Pointer[deciderHolder]

// SetDropDecider installs the package-level DropDecider consulted on every
// workflow event emit. Passing nil restores the default no-op behavior.
func SetDropDecider(d DropDecider) {
	globalDropDecider.Store(&deciderHolder{d: d})
}

// shouldDropEvent is nil-safe: with no decider installed (the default) it
// always returns false and existing emit behavior is unchanged.
func shouldDropEvent(ctx context.Context, id DropIdentity, eventType string) bool {
	h := globalDropDecider.Load()
	if h == nil || h.d == nil {
		return false
	}
	return h.d.ShouldDrop(ctx, id, eventType)
}

// v1Identified matches v1 workflow event protos carrying WorkflowMetadata.
type v1Identified interface {
	GetM() *protoevents.WorkflowMetadata
}

// v2Identified matches v2 workflow event protos carrying a WorkflowKey and an
// execution ID.
type v2Identified interface {
	GetWorkflow() *eventsv2.WorkflowKey
	GetWorkflowExecutionID() string
}

// dropIdentityFor extracts the drop identity from an event proto. Messages
// without both owner and execution ID (e.g. ExecutionProfile, which carries no
// owner, or MeteringReport) yield a partial/zero identity and are therefore
// never dropped (fail-closed).
func dropIdentityFor(msg proto.Message) DropIdentity {
	switch m := msg.(type) {
	case v2Identified:
		return DropIdentity{
			WorkflowOwner:       m.GetWorkflow().GetWorkflowOwner(),
			WorkflowExecutionID: m.GetWorkflowExecutionID(),
		}
	case v1Identified:
		md := m.GetM()
		return DropIdentity{
			WorkflowOwner:       md.GetWorkflowOwner(),
			WorkflowExecutionID: md.GetWorkflowExecutionID(),
		}
	default:
		return DropIdentity{}
	}
}

// FaultInjectorConfig configures a FaultInjector. Mirrors the
// [Telemetry.WorkflowFaultInjection] TOML section.
type FaultInjectorConfig struct {
	// Enabled turns injection on. Hard off by default.
	Enabled bool
	// OwnerAllowlist is the exact set of workflow owner addresses whose events
	// may be dropped. Compared case-insensitively, optional 0x prefix ignored.
	// Empty means nothing is ever dropped.
	OwnerAllowlist []string
	// RateBps is the fraction of allowlisted executions to drop, in basis
	// points of 10000 (100 = 1%). 0 disables dropping.
	RateBps int
	// Seed feeds the deterministic drop hash. Rotate to select a new cohort.
	Seed string
	// Level is DropLevelFinishedOnly (1) or DropLevelAllLifecycle (2).
	Level int
}

// FaultInjector is the default DropDecider implementation: a deterministic,
// fail-closed, known-answer fault injector.
type FaultInjector struct {
	enabled   bool
	allowlist map[string]struct{}
	rateBps   uint64
	seed      string
	level     int
	lggr      logger.Logger
	dropped   metric.Int64Counter
}

var _ DropDecider = (*FaultInjector)(nil)

// NewFaultInjector validates cfg and returns a FaultInjector. Construction
// fails hard on invalid config so a misconfigured injector never runs.
func NewFaultInjector(cfg FaultInjectorConfig, lggr logger.Logger) (*FaultInjector, error) {
	if lggr == nil {
		return nil, errors.New("fault injector requires a logger: the injection record is the detector's source of truth")
	}
	if cfg.RateBps < 0 || cfg.RateBps > 10000 {
		return nil, fmt.Errorf("RateBps must be in [0, 10000], got %d", cfg.RateBps)
	}
	if cfg.Level != DropLevelFinishedOnly && cfg.Level != DropLevelAllLifecycle {
		return nil, fmt.Errorf("Level must be %d (WorkflowExecutionFinished only) or %d (all lifecycle events), got %d",
			DropLevelFinishedOnly, DropLevelAllLifecycle, cfg.Level)
	}
	allowlist := make(map[string]struct{}, len(cfg.OwnerAllowlist))
	for _, owner := range cfg.OwnerAllowlist {
		normalized := normalizeOwner(owner)
		if normalized == "" {
			return nil, fmt.Errorf("OwnerAllowlist entry %q is empty after normalization", owner)
		}
		allowlist[normalized] = struct{}{}
	}
	dropped, err := beholder.GetMeter().Int64Counter(droppedTotalMetricName)
	if err != nil {
		return nil, fmt.Errorf("failed to register %s counter: %w", droppedTotalMetricName, err)
	}
	return &FaultInjector{
		enabled:   cfg.Enabled,
		allowlist: allowlist,
		rateBps:   uint64(cfg.RateBps), //nolint:gosec // G115: validated to be in [0, 10000] above
		seed:      cfg.Seed,
		level:     cfg.Level,
		lggr:      lggr,
		dropped:   dropped,
	}, nil
}

// ShouldDrop implements DropDecider. It drops iff ALL of: enabled, owner
// non-empty and allowlisted, event type matches the configured level, and
// hash(seed ‖ executionID) mod 10000 < rateBps. Any missing identity => never
// drop (fail-closed). Nil-safe: a nil *FaultInjector never drops.
func (f *FaultInjector) ShouldDrop(ctx context.Context, id DropIdentity, eventType string) bool {
	if f == nil || !f.enabled {
		return false
	}
	// Fail-closed: without both owner and execution ID the drop can neither be
	// safely scoped nor reconciled, so the event always flows.
	if id.WorkflowOwner == "" || id.WorkflowExecutionID == "" {
		return false
	}
	if _, ok := f.allowlist[normalizeOwner(id.WorkflowOwner)]; !ok {
		return false
	}
	if !f.levelMatches(eventType) {
		return false
	}
	if !deterministicDrop(f.seed, id.WorkflowExecutionID, f.rateBps) {
		return false
	}

	// Injection record: parsed by an external detector — msg and field keys
	// are a contract (see DroppedWorkflowEventMsg).
	f.lggr.Infow(DroppedWorkflowEventMsg,
		"executionID", id.WorkflowExecutionID,
		"eventType", eventType,
		"workflowOwner", id.WorkflowOwner,
	)
	f.dropped.Add(ctx, 1, metric.WithAttributes(attribute.String("event_type", eventType)))
	return true
}

// levelMatches reports whether eventType (a beholder entity such as
// "workflows.v2.WorkflowExecutionFinished") is droppable at the configured
// level. Matching is on the unqualified type name so v1 and v2 variants of the
// same event are treated identically — dropping only one would inject nothing,
// since the other still reaches downstream.
func (f *FaultInjector) levelMatches(eventType string) bool {
	name := eventType
	if i := strings.LastIndexByte(eventType, '.'); i >= 0 {
		name = eventType[i+1:]
	}
	switch f.level {
	case DropLevelFinishedOnly:
		return name == WorkflowExecutionFinished
	case DropLevelAllLifecycle:
		_, ok := lifecycleEventTypes[name]
		return ok
	default:
		return false
	}
}

// normalizeOwner canonicalizes an owner address for exact-match comparison:
// trim whitespace, strip an optional 0x prefix, lowercase. No globs/prefixes.
func normalizeOwner(owner string) string {
	owner = strings.TrimSpace(owner)
	owner = strings.TrimPrefix(owner, "0x")
	owner = strings.TrimPrefix(owner, "0X")
	return strings.ToLower(owner)
}

// deterministicDrop selects a stable rateBps/10000 slice of executions:
// sha256(seed ‖ executionID) first 8 bytes big-endian mod 10000 < rateBps.
// Pure function of its inputs — every node in a DON computes the identical
// decision independently, and a reconciler can recompute the expected drop set
// offline from (seed, rateBps, executionIDs).
func deterministicDrop(seed, executionID string, rateBps uint64) bool {
	if rateBps == 0 {
		return false
	}
	h := sha256.New()
	h.Write([]byte(seed))
	h.Write([]byte(executionID))
	v := binary.BigEndian.Uint64(h.Sum(nil)[:8])
	return v%10000 < rateBps
}
