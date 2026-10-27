package v2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/smartcontractkit/chainlink-common/pkg/beholder"
	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/sharding"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2/triggers"
)

// cachedExpiry is how long a cached trigger event is kept for potential failover replay.
const cachedExpiry = 10 * time.Minute

const timeoutBasedFailoverTick = time.Second

const failoverAutoExecutionTotalMetric = "failover_auto_execution_total"

// ShardFailoverManager wraps a workflow Engine and handles shard ownership
// decisions externally, keeping the Engine oblivious to sharding. On the
// primary shard it allows all trigger events through. On a secondary shard
// with failover enabled it caches trigger events so they can be replayed if
// the primary reports a system failure.
//
// The ownership check runs wherever the wrapped engine receives events. A legacy
// engine registers its own triggers and calls the OnTriggerAdmission hook before
// enqueueing. A coordinated engine relies on the trigger coordinator, which
// subscribes and delivers through this manager, so ExecuteTrigger runs the check
// before delegating. Either way the manager only decides; whoever holds the
// trigger handles acts on a denial (ACK and drop) or a cached event (no ACK).
//
// Lifecycle:
//   - NewShardFailoverManager(cfg) creates the manager (no engine yet).
//   - WireHooks(engineCfg) wires the admission and status-update hooks on
//     the EngineConfig before the engine is created.
//   - SetEngine(engine) injects the engine after creation.
//   - Start(ctx) wires the communicator and starts the engine.
//
// Shard ownership is resolved on each event (trigger admission, execution
// status update) by querying the ShardResolver — no static role assignment
// or background polling loop. If the orchestrator reassigns ownership at
// runtime, the next event automatically picks up the new role.
type ShardFailoverManager struct {
	services.Service
	eng *services.Engine

	engine v2.WorkflowEngine
	cfg    ShardFailoverManagerConfig

	autoExecCounter metric.Int64Counter

	mu    sync.RWMutex
	cache map[string]cachedEvent
}

type cachedEvent struct {
	event      triggers.CoordinatedEvent
	cachedAt   time.Time
	failoverAt time.Time
}

type ShardFailoverManagerConfig struct {
	ShardingEnabled bool
	MyDONID         uint32
	MyShardIndex    uint32
	WorkflowID      string
	WorkflowOwner   string

	ShardResolver           shardownership.ShardResolver
	ShardOrchestratorClient shardorchestrator.ClientInterface
	ShardRoutingSteady      *shardownership.SteadySignal

	// ShardDonLookup resolves a DON ID to the current DON, e.g. via
	// capRegistry.DONByID. ResolveAllShards always returns real DON IDs
	// (shardownership.manualShardResolver translates any configured shard
	// index to a DON ID internally), so no shard-index remapping happens here.
	ShardDonLookup func(ctx context.Context, donID uint32) *commoncap.DON
	DonSubscriber  capabilities.DonSubscriber

	Failover FailoverConfig

	Logger logger.Logger
}

// FailoverConfig groups the failover machinery: the gates and window that
// decide when a secondary shard takes over an owned workflow, and the
// communicator the shards exchange execution outcomes over.
type FailoverConfig struct {
	ShardingFailoverGate limits.GateLimiter
	// TimeoutBasedFailoverGate opens the timeout-based failover path: a secondary executes
	// a cached trigger event itself if no execution outcome from the primary
	// arrives within TimeoutBasedFailoverWindow (the silent-primary-death case).
	TimeoutBasedFailoverGate limits.GateLimiter
	// TimeoutBasedFailoverWindow is how long a cached trigger event waits for the
	// primary's execution outcome before the secondary executes it itself.
	// A nil limiter (or closed gate) disables timeout-based failover.
	TimeoutBasedFailoverWindow limits.TimeLimiter
	Communicator               *sharding.ShardFailoverCommunicator
}

// NewShardFailoverManager creates a manager that is not yet wired to an engine.
// Call WireHooks before engine creation and SetEngine after.
func NewShardFailoverManager(cfg ShardFailoverManagerConfig) *ShardFailoverManager {
	m := &ShardFailoverManager{
		cfg:   cfg,
		cache: make(map[string]cachedEvent),
	}
	autoExecCounter, err := beholder.GetMeter().Int64Counter(failoverAutoExecutionTotalMetric)
	if err != nil {
		m.cfg.Logger.Warnw("shard failover: failed to register failover_auto_execution_total counter", "err", err)
	} else {
		m.autoExecCounter = autoExecCounter
	}
	m.Service, m.eng = services.Config{
		Name:  "ShardFailoverManager",
		Start: m.start,
		Close: m.close,
	}.NewServiceEngine(cfg.Logger)
	return m
}

// WireHooks sets the admission and execution-status-update hooks on the
// EngineConfig so the engine delegates shard decisions to this manager.
// Must be called before v2.NewEngine.
func (m *ShardFailoverManager) WireHooks(cfg *v2.EngineConfig) {
	cfg.Hooks.OnTriggerAdmission = m.admissionCheck
	cfg.Hooks.OnExecutionStatusUpdate = m.forwardExecutionStatus
}

// SetEngine injects the engine after it has been created. Required before Start.
func (m *ShardFailoverManager) SetEngine(engine v2.WorkflowEngine) {
	m.engine = engine
}

func (m *ShardFailoverManager) start(ctx context.Context) error {
	if m.engine == nil {
		return errors.New("engine not set, call SetEngine before Start")
	}

	if m.cfg.Failover.Communicator != nil {
		if err := m.wireFailover(ctx); err != nil {
			return fmt.Errorf("failed to wire failover: %w", err)
		}
	}

	if err := m.engine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}

	m.eng.GoCtx(ctx, m.pruneLoop)
	m.eng.GoCtx(ctx, m.timeoutBasedFailoverLoop)
	return nil
}

func (m *ShardFailoverManager) close() error {
	if m.cfg.Failover.Communicator != nil {
		m.cfg.Failover.Communicator.UnregisterHandler(m.cfg.WorkflowID)
	}
	if m.engine != nil {
		return m.engine.Close()
	}
	return nil
}

// admissionCheck decides whether the engine should process a trigger event,
// checking shard ownership dynamically per-trigger. It is wired as a legacy
// engine's OnTriggerAdmission hook and called by ExecuteTrigger for a
// coordinated engine.
func (m *ShardFailoverManager) admissionCheck(ctx context.Context, event triggers.CoordinatedEvent) error {
	if !m.cfg.ShardingEnabled {
		return nil
	}

	verdict := m.checkShardOwnership(ctx)

	switch verdict {
	case shardownership.Allow:
		return nil
	case shardownership.DenyNotOwner:
		if m.cfg.Failover.ShardingFailoverGate != nil && m.cfg.Failover.ShardingFailoverGate.AllowErr(ctx) == nil {
			m.cacheEvent(ctx, event)
			return v2.ErrAdmissionCache
		}
		return v2.ErrShardDeniedNotOwner
	case shardownership.DenyOrchestratorError:
		return v2.ErrShardDeniedOrchestrator
	default:
		return nil
	}
}

// forwardExecutionStatus is wired as the engine's OnExecutionStatusUpdate
// hook. It resolves shard ownership on each call: on the primary shard it
// sends the status to the secondary via the communicator; on the secondary
// it is a no-op. This naturally handles role changes — if the orchestrator
// reassigns ownership, the next execution picks up the new role.
func (m *ShardFailoverManager) forwardExecutionStatus(workflowID string, executionID string, triggerEventID string, triggerIndex int, status string, errClass events.ErrorClassification) {
	if m.cfg.Failover.Communicator == nil {
		return
	}
	if !m.isPrimaryShard(context.Background(), m.cfg.WorkflowID, m.cfg.WorkflowOwner) {
		return
	}
	// Position 1 in ResolveAllShards' result is the secondary DON ID for this
	// workflow; ShardDonLookup resolves it to the current DON.
	secondaryDon := m.resolveDon(context.Background(), m.cfg.WorkflowID, m.cfg.WorkflowOwner, 1)
	if secondaryDon == nil {
		m.cfg.Logger.Warnw("failover: primary cannot resolve secondary DON, skipping status update",
			"workflowID", workflowID, "myShardIndex", m.cfg.MyShardIndex)
		return
	}
	execStatus := mapExecutionStatus(status, errClass)
	m.cfg.Failover.Communicator.Send(context.Background(), &ringpb.ExecutionStatusUpdate{
		WorkflowId:     workflowID,
		ExecutionId:    executionID,
		TriggerEventId: triggerEventID,
		TriggerIndex:   uint32(triggerIndex), //nolint:gosec // G115: triggerIndex is small
		Status:         execStatus,
		PrimaryDonId:   m.cfg.MyDONID,
	}, *secondaryDon)
}

// HandleExecutionStatusUpdate is called on the secondary shard when the
// primary reports an execution outcome. On system failure it replays the
// cached trigger event through the engine. Replay calls the wrapped engine
// directly, not ExecuteTrigger: this node is still not the owner, so the
// admission check would only cache the event again.
func (m *ShardFailoverManager) HandleExecutionStatusUpdate(msg *ringpb.ExecutionStatusUpdate) {
	switch msg.Status {
	case ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED:
	case ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS,
		ringpb.ExecutionStatus_EXECUTION_STATUS_USER_ERROR:
		m.removeCachedEvent(msg.TriggerEventId)

	case ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR:
		event, ok := m.popCachedEvent(msg.TriggerEventId, monitoring.ShardFailoverCacheExitReplayed)
		if !ok {
			m.cfg.Logger.Warnw("failover: no cached event for triggerEventID, cannot replay",
				"triggerEventID", msg.TriggerEventId,
				"workflowID", msg.WorkflowId,
				"executionID", msg.ExecutionId)
			shardFailoverMetrics().RecordReplay(context.Background(), monitoring.ShardFailoverReplayMiss)
			return
		}
		m.cfg.Logger.Infow("failover: replaying cached trigger event after system error",
			"triggerEventID", msg.TriggerEventId,
			"workflowID", msg.WorkflowId,
			"executionID", msg.ExecutionId)
		replayStart := time.Now()
		err := m.engine.ExecuteTrigger(m.replayExecCtx(context.Background()), event)
		outcome := monitoring.ShardFailoverReplaySucceeded
		if err != nil {
			outcome = monitoring.ShardFailoverReplayFailed
			m.cfg.Logger.Errorw("failover: failed to replay cached trigger event",
				"triggerEventID", msg.TriggerEventId,
				"workflowID", msg.WorkflowId,
				"executionID", msg.ExecutionId,
				"err", err)
		}
		fm := shardFailoverMetrics()
		fm.RecordReplay(context.Background(), outcome)
		fm.RecordReplayDuration(context.Background(), outcome, time.Since(replayStart))
	}
}

// replayExecCtx stamps the workflow's tenant identity onto the replay context,
// matching the engine's queue path (capability calls fail without it).
func (m *ShardFailoverManager) replayExecCtx(parent context.Context) context.Context {
	return contexts.WithCRE(parent, m.engine.Tenant())
}

// checkShardOwnership performs a dynamic per-trigger shard ownership check.
func (m *ShardFailoverManager) checkShardOwnership(ctx context.Context) shardownership.Verdict {
	needShardOwnerCheck := m.cfg.ShardRoutingSteady == nil || !m.cfg.ShardRoutingSteady.SkipCommittedOwnerCheck()
	if !needShardOwnerCheck {
		return shardownership.Allow
	}

	switch {
	case m.cfg.ShardResolver != nil:
		donID, found, err := m.cfg.ShardResolver.ResolveShard(ctx, m.cfg.WorkflowID, m.cfg.WorkflowOwner)
		if err != nil {
			return shardownership.DenyOrchestratorError
		}
		if !found || donID != m.cfg.MyDONID {
			return shardownership.DenyNotOwner
		}
		return shardownership.Allow
	case m.cfg.ShardOrchestratorClient != nil:
		verdict, _, _ := shardownership.CheckCommittedOwner(ctx, m.cfg.ShardOrchestratorClient, m.cfg.WorkflowID, m.cfg.MyDONID)
		return verdict
	default:
		return shardownership.Allow
	}
}

// shardFailoverMetrics returns the process-global failover instruments, or nil
// when they failed to register (all methods are nil-receiver safe).
func shardFailoverMetrics() *monitoring.ShardFailoverMetrics {
	m, _ := monitoring.GlobalShardFailoverMetrics()
	return m
}

func (m *ShardFailoverManager) cacheEvent(ctx context.Context, event triggers.CoordinatedEvent) {
	eventID := event.Event.Event.ID
	failoverAt := m.timeoutBasedFailoverDeadline(ctx)
	m.mu.Lock()
	m.cache[eventID] = cachedEvent{event: event, cachedAt: time.Now(), failoverAt: failoverAt}
	cached := int64(len(m.cache))
	m.mu.Unlock()
	fm := shardFailoverMetrics()
	fm.RecordTriggerCached(ctx)
	fm.UpdateCachedEvents(ctx, cached)
	m.cfg.Logger.Infow("secondary shard: cached trigger event for failover",
		"eventID", eventID,
		"myShardIndex", m.cfg.MyShardIndex,
		"workflowID", m.cfg.WorkflowID,
		"triggerIndex", event.TriggerIndex,
		"failoverAt", failoverAt)
}

// timeoutBasedFailoverDeadline returns the instant after which this node may execute a
// just-cached event itself, or the zero time if automatic failover is disabled.
func (m *ShardFailoverManager) timeoutBasedFailoverDeadline(ctx context.Context) time.Time {
	if m.cfg.Failover.TimeoutBasedFailoverGate == nil || m.cfg.Failover.TimeoutBasedFailoverWindow == nil {
		return time.Time{}
	}
	if m.cfg.Failover.TimeoutBasedFailoverGate.AllowErr(ctx) != nil {
		return time.Time{}
	}
	window, err := m.cfg.Failover.TimeoutBasedFailoverWindow.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			m.cfg.Logger.Warnw("secondary shard: failed to read failover auto window, automatic failover disabled for this event", "err", err)
			return time.Time{}
		}
		m.cfg.Logger.Errorw("secondary shard: failed to read failover auto window; continuing with the value the limiter returned", "err", err)
	}
	if window <= 0 {
		return time.Time{}
	}
	return time.Now().Add(window)
}

// timeoutBasedFailoverLoop periodically executes cached trigger events whose failover
// window elapsed without an execution outcome from the primary shard.
func (m *ShardFailoverManager) timeoutBasedFailoverLoop(ctx context.Context) {
	if m.cfg.Failover.TimeoutBasedFailoverGate == nil || m.cfg.Failover.TimeoutBasedFailoverWindow == nil {
		return
	}
	ticker := time.NewTicker(timeoutBasedFailoverTick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.executeDueTimeoutBasedFailovers(ctx)
		}
	}
}

// executeDueTimeoutBasedFailovers runs one auto-failover pass over cached events whose
// failover deadline has elapsed.
func (m *ShardFailoverManager) executeDueTimeoutBasedFailovers(ctx context.Context) {
	if m.engine == nil {
		return
	}

	now := time.Now()
	type candidate struct {
		eventID string
		event   triggers.CoordinatedEvent
	}
	var due []candidate
	m.mu.RLock()
	for id, ce := range m.cache {
		if ce.failoverAt.IsZero() || now.Before(ce.failoverAt) {
			continue
		}
		due = append(due, candidate{eventID: id, event: ce.event})
	}
	m.mu.RUnlock()

	for _, c := range due {
		if m.cfg.Failover.TimeoutBasedFailoverGate.AllowErr(ctx) != nil {
			return
		}

		switch m.checkShardOwnership(ctx) {
		case shardownership.DenyNotOwner:
			event, ok := m.popCachedEvent(c.eventID, monitoring.ShardFailoverCacheExitReplayed)
			if !ok {
				continue
			}
			m.cfg.Logger.Infow("secondary shard: auto failover executed cached trigger event",
				"eventID", c.eventID,
				"myShardIndex", m.cfg.MyShardIndex,
				"workflowID", m.cfg.WorkflowID,
				"triggerIndex", event.TriggerIndex)
			if m.autoExecCounter != nil {
				m.autoExecCounter.Add(ctx, 1,
					metric.WithAttributes(
						attribute.String(platform.KeyWorkflowID, m.cfg.WorkflowID),
						attribute.String(platform.KeyOrganizationID, m.engine.Tenant().Org),
					))
			}
			// Async: ExecuteTrigger blocks for the execution's duration, and
			// serializing due events would prevent concurrent executions of
			// this workflow.
			m.eng.GoCtx(ctx, func(ctx context.Context) {
				if err := m.engine.ExecuteTrigger(m.replayExecCtx(ctx), event); err != nil {
					m.cfg.Logger.Errorw("secondary shard: auto failover failed to execute cached trigger event",
						"eventID", c.eventID,
						"workflowID", m.cfg.WorkflowID,
						"triggerIndex", event.TriggerIndex,
						"err", err)
				}
			})
		case shardownership.Allow:
			if _, ok := m.popCachedEvent(c.eventID, monitoring.ShardFailoverCacheExitDrained); ok {
				m.cfg.Logger.Infow("secondary shard: dropped cached trigger event, shard is owner again",
					"eventID", c.eventID,
					"workflowID", m.cfg.WorkflowID)
			}
		case shardownership.DenyOrchestratorError:
		}
	}
}

// popCachedEvent removes eventID from the cache for replay. exit is the
// ShardFailoverCacheExit* constant describing why the entry left the cache.
func (m *ShardFailoverManager) popCachedEvent(eventID string, exit string) (triggers.CoordinatedEvent, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ce, ok := m.cache[eventID]
	if !ok {
		return triggers.CoordinatedEvent{}, false
	}
	delete(m.cache, eventID)
	fm := shardFailoverMetrics()
	fm.RecordCacheExit(context.Background(), exit, 1)
	fm.UpdateCachedEvents(context.Background(), int64(len(m.cache)))
	return ce.event, true
}

func (m *ShardFailoverManager) removeCachedEvent(eventID string) {
	m.mu.Lock()
	_, cached := m.cache[eventID]
	if cached {
		delete(m.cache, eventID)
	}
	size := int64(len(m.cache))
	m.mu.Unlock()
	if !cached {
		return
	}
	fm := shardFailoverMetrics()
	fm.RecordCacheExit(context.Background(), monitoring.ShardFailoverCacheExitDrained, 1)
	fm.UpdateCachedEvents(context.Background(), size)
}

func (m *ShardFailoverManager) pruneLoop(ctx context.Context) {
	ticker := time.NewTicker(cachedExpiry / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.pruneExpiredCachedEvents(ctx)
		}
	}
}

// pruneExpiredCachedEvents drops entries older than cachedExpiry: the primary
// never delivered a verdict for them, so they exit the funnel as expired.
func (m *ShardFailoverManager) pruneExpiredCachedEvents(ctx context.Context) {
	m.mu.Lock()
	now := time.Now()
	expired := 0
	for id, ce := range m.cache {
		if now.Sub(ce.cachedAt) > cachedExpiry {
			delete(m.cache, id)
			expired++
		}
	}
	size := int64(len(m.cache))
	m.mu.Unlock()
	if expired == 0 {
		return
	}
	fm := shardFailoverMetrics()
	fm.RecordCacheExit(ctx, monitoring.ShardFailoverCacheExitExpired, int64(expired))
	fm.UpdateCachedEvents(ctx, size)
}

// wireFailover sets up the communicator for ExecutionStatusUpdate messages.
// It resolves the primary shard DON and registers a per-workflow handler
// (with the primary DON for quorum validation) on the shared communicator.
// The handler is always registered — when this node is primary, nobody
// sends to us so the handler is simply never invoked. Shard ownership is
// then resolved per-event in admissionCheck and forwardExecutionStatus.
func (m *ShardFailoverManager) wireFailover(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	sub, unsub, err := m.cfg.DonSubscriber.Subscribe(ctx)
	if err != nil {
		m.cfg.Logger.Warnw("shard failover: failed to subscribe to DON updates, skipping hook wiring", "err", err)
		return nil
	}
	defer unsub()

	select {
	case <-sub:
	case <-ctx.Done():
		m.cfg.Logger.Warnw("shard failover: timed out waiting for DON info, skipping hook wiring")
		return nil
	}

	primaryDon := m.resolveDon(ctx, m.cfg.WorkflowID, m.cfg.WorkflowOwner, 0)

	var primaryDonVal commoncap.DON
	if primaryDon != nil {
		primaryDonVal = *primaryDon
	}

	m.cfg.Failover.Communicator.RegisterHandler(m.cfg.WorkflowID, primaryDonVal, m.HandleExecutionStatusUpdate)
	m.cfg.Logger.Infow("shard failover: wired communicator",
		"myShardIndex", m.cfg.MyShardIndex,
		"workflowID", m.cfg.WorkflowID,
		"primaryDonID", primaryDonIDOrZero(primaryDon))

	if err := m.cfg.Failover.Communicator.Start(ctx); err != nil {
		return fmt.Errorf("failed to start communicator: %w", err)
	}
	return nil
}

func primaryDonIDOrZero(d *commoncap.DON) uint32 {
	if d == nil {
		return 0
	}
	return d.ID
}

func (m *ShardFailoverManager) isPrimaryShard(ctx context.Context, workflowID, owner string) bool {
	if m.cfg.ShardResolver == nil {
		return m.cfg.MyShardIndex == 0
	}
	allResolver, ok := m.cfg.ShardResolver.(shardownership.AllShardsResolver)
	if !ok {
		return m.cfg.MyShardIndex == 0
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	shards, found, err := allResolver.ResolveAllShards(ctx, workflowID, owner)
	if err != nil || !found || len(shards) == 0 {
		return m.cfg.MyShardIndex == 0
	}
	return shards[0] == m.cfg.MyDONID
}

func (m *ShardFailoverManager) resolveDon(ctx context.Context, workflowID, owner string, shardIndex int) *commoncap.DON {
	if m.cfg.ShardDonLookup == nil {
		return nil
	}
	allResolver, ok := m.cfg.ShardResolver.(shardownership.AllShardsResolver)
	if !ok {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	shards, found, err := allResolver.ResolveAllShards(ctx, workflowID, owner)
	if err != nil || !found || len(shards) <= shardIndex {
		return nil
	}
	return m.cfg.ShardDonLookup(ctx, shards[shardIndex])
}

func mapExecutionStatus(status string, errClass events.ErrorClassification) ringpb.ExecutionStatus {
	switch status {
	case store.StatusCompleted, store.StatusCompletedEarlyExit:
		return ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS
	case store.StatusErrored, store.StatusTimeout:
		if errClass == events.ErrorClassificationUser {
			return ringpb.ExecutionStatus_EXECUTION_STATUS_USER_ERROR
		}
		return ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR
	default:
		return ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED
	}
}

// Drain delegates to the wrapped engine.
func (m *ShardFailoverManager) Drain() bool { return m.engine.Drain() }

// ActiveExecutions delegates to the wrapped engine.
func (m *ShardFailoverManager) ActiveExecutions() int32 { return m.engine.ActiveExecutions() }

// DrainStartedAt delegates to the wrapped engine.
func (m *ShardFailoverManager) DrainStartedAt() (time.Time, bool) { return m.engine.DrainStartedAt() }

// ExecuteTrigger is how the trigger coordinator delivers events to a
// coordinated engine on a sharded node. The coordinated engine never calls
// OnTriggerAdmission, so the shard ownership check runs here first. An
// admission error (triggers.ErrAdmissionDenied or triggers.ErrAdmissionCached)
// is returned for the coordinator to act on; otherwise the event is delegated.
func (m *ShardFailoverManager) ExecuteTrigger(ctx context.Context, event triggers.CoordinatedEvent) error {
	if err := m.admissionCheck(ctx, event); err != nil {
		return err
	}
	return m.engine.ExecuteTrigger(ctx, event)
}

// Subscribe delegates to the wrapped engine. The trigger coordinator calls it
// while registering triggers for a coordinated engine.
func (m *ShardFailoverManager) Subscribe(ctx context.Context) ([]*sdkpb.TriggerSubscription, error) {
	return m.engine.Subscribe(ctx)
}

// Tenant delegates to the wrapped engine.
func (m *ShardFailoverManager) Tenant() contexts.CRE { return m.engine.Tenant() }

// SdkName delegates to the wrapped engine, same reasoning as Tenant.
func (m *ShardFailoverManager) SdkName() string { return m.engine.SdkName() }

// IsCoordinated delegates to the wrapped engine, so the registry classifies the
// manager as whichever engine it wraps.
func (m *ShardFailoverManager) IsCoordinated() bool { return m.engine.IsCoordinated() }

var (
	_ DrainableService  = (*ShardFailoverManager)(nil)
	_ v2.WorkflowEngine = (*ShardFailoverManager)(nil)
)
