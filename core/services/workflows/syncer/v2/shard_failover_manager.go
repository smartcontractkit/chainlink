package v2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	commoncap "github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	ringpb "github.com/smartcontractkit/chainlink-protos/ring/go"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/remote/sharding"
	"github.com/smartcontractkit/chainlink/v2/core/services/shardorchestrator"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/shardownership"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	v2 "github.com/smartcontractkit/chainlink/v2/core/services/workflows/v2"
)

// cachedExpiry is how long a cached trigger event is kept for potential failover replay.
const cachedExpiry = 10 * time.Minute

// ShardFailoverManager wraps a workflow Engine and handles shard ownership
// decisions externally, keeping the Engine oblivious to sharding. On the
// primary shard it allows all trigger events through. On a secondary shard
// with failover enabled it caches trigger events so they can be replayed if
// the primary reports a system failure.
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

	engine *v2.Engine
	cfg    ShardFailoverManagerConfig

	mu    sync.RWMutex
	cache map[string]cachedEvent
}

type cachedEvent struct {
	event    v2.RoutedTriggerEvent
	cachedAt time.Time
}

type ShardFailoverManagerConfig struct {
	ShardingEnabled bool
	MyShardID       uint32
	WorkflowID      string
	WorkflowOwner   string

	ShardResolver           shardownership.ShardResolver
	ShardOrchestratorClient shardorchestrator.ClientInterface
	ShardRoutingSteady      *shardownership.SteadySignal

	FailoverGate   limits.GateLimiter
	Communicator   *sharding.ShardFailoverCommunicator
	ShardDonLookup func(ctx context.Context, shardID uint32) *commoncap.DON
	DonSubscriber  capabilities.DonSubscriber

	Logger logger.Logger
}

// NewShardFailoverManager creates a manager that is not yet wired to an engine.
// Call WireHooks before engine creation and SetEngine after.
func NewShardFailoverManager(cfg ShardFailoverManagerConfig) *ShardFailoverManager {
	m := &ShardFailoverManager{
		cfg:   cfg,
		cache: make(map[string]cachedEvent),
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
func (m *ShardFailoverManager) SetEngine(engine *v2.Engine) {
	m.engine = engine
}

func (m *ShardFailoverManager) start(ctx context.Context) error {
	if m.engine == nil {
		return errors.New("engine not set, call SetEngine before Start")
	}

	if m.cfg.Communicator != nil {
		if err := m.wireFailover(ctx); err != nil {
			return fmt.Errorf("failed to wire failover: %w", err)
		}
	}

	if err := m.engine.Start(ctx); err != nil {
		return fmt.Errorf("failed to start engine: %w", err)
	}

	m.eng.GoCtx(ctx, m.pruneLoop)
	return nil
}

func (m *ShardFailoverManager) close() error {
	if m.cfg.Communicator != nil {
		m.cfg.Communicator.UnregisterHandler(m.cfg.WorkflowID)
	}
	if m.engine != nil {
		return m.engine.Close()
	}
	return nil
}

// admissionCheck is wired as the engine's OnTriggerAdmission hook. It
// checks shard ownership dynamically per-trigger and decides whether the
// engine should process the event.
func (m *ShardFailoverManager) admissionCheck(ctx context.Context, event v2.RoutedTriggerEvent) error {
	if !m.cfg.ShardingEnabled {
		return nil
	}

	verdict := m.checkShardOwnership(ctx)

	switch verdict {
	case shardownership.Allow:
		return nil
	case shardownership.DenyNotOwner:
		if m.cfg.FailoverGate != nil && m.cfg.FailoverGate.AllowErr(ctx) == nil {
			m.cacheEvent(event)
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
// it is a no-op. This naturally handles role changes without a loop — if
// the orchestrator reassigns ownership, the next execution picks up the
// new role.
func (m *ShardFailoverManager) forwardExecutionStatus(workflowID string, executionID string, triggerEventID string, triggerIndex int, status string, errClass events.ErrorClassification) {
	if m.cfg.Communicator == nil {
		return
	}
	if !m.isPrimaryShard(context.Background(), m.cfg.WorkflowID, m.cfg.WorkflowOwner) {
		return
	}
	execStatus := mapExecutionStatus(status, errClass)
	m.cfg.Communicator.Send(context.Background(), &ringpb.ExecutionStatusUpdate{
		WorkflowId:     workflowID,
		ExecutionId:    executionID,
		TriggerEventId: triggerEventID,
		TriggerIndex:   uint32(triggerIndex), //nolint:gosec // G115: triggerIndex is small
		Status:         execStatus,
		PrimaryDonId:   m.cfg.MyShardID,
	})
}

// HandleExecutionStatusUpdate is called on the secondary shard when the
// primary reports an execution outcome. On system failure it replays the
// cached trigger event through the engine.
func (m *ShardFailoverManager) HandleExecutionStatusUpdate(msg *ringpb.ExecutionStatusUpdate) {
	switch msg.Status {
	case ringpb.ExecutionStatus_EXECUTION_STATUS_UNSPECIFIED:
	case ringpb.ExecutionStatus_EXECUTION_STATUS_SUCCESS,
		ringpb.ExecutionStatus_EXECUTION_STATUS_USER_ERROR:
		m.removeCachedEvent(msg.TriggerEventId)

	case ringpb.ExecutionStatus_EXECUTION_STATUS_SYSTEM_ERROR:
		event, ok := m.popCachedEvent(msg.TriggerEventId)
		if !ok {
			m.cfg.Logger.Warnw("failover: no cached event for triggerEventID, cannot replay",
				"triggerEventID", msg.TriggerEventId,
				"workflowID", msg.WorkflowId,
				"executionID", msg.ExecutionId)
			return
		}
		m.cfg.Logger.Infow("failover: replaying cached trigger event after system error",
			"triggerEventID", msg.TriggerEventId,
			"workflowID", msg.WorkflowId,
			"executionID", msg.ExecutionId)
		if err := m.engine.ExecuteTrigger(context.Background(), event); err != nil {
			m.cfg.Logger.Errorw("failover: failed to replay cached trigger event",
				"triggerEventID", msg.TriggerEventId,
				"workflowID", msg.WorkflowId,
				"executionID", msg.ExecutionId,
				"err", err)
		}
	}
}

// checkShardOwnership performs a dynamic per-trigger shard ownership check.
func (m *ShardFailoverManager) checkShardOwnership(ctx context.Context) shardownership.Verdict {
	needShardOwnerCheck := m.cfg.ShardRoutingSteady == nil || !m.cfg.ShardRoutingSteady.SkipCommittedOwnerCheck()
	if !needShardOwnerCheck {
		return shardownership.Allow
	}

	switch {
	case m.cfg.ShardResolver != nil:
		shardID, found, err := m.cfg.ShardResolver.ResolveShard(ctx, m.cfg.WorkflowID, m.cfg.WorkflowOwner)
		if err != nil {
			return shardownership.DenyOrchestratorError
		}
		if !found || shardID != m.cfg.MyShardID {
			return shardownership.DenyNotOwner
		}
		return shardownership.Allow
	case m.cfg.ShardOrchestratorClient != nil:
		verdict, _, _ := shardownership.CheckCommittedOwner(ctx, m.cfg.ShardOrchestratorClient, m.cfg.WorkflowID, m.cfg.MyShardID)
		return verdict
	default:
		return shardownership.Allow
	}
}

func (m *ShardFailoverManager) cacheEvent(event v2.RoutedTriggerEvent) {
	eventID := event.Event.Event.ID
	m.mu.Lock()
	m.cache[eventID] = cachedEvent{event: event, cachedAt: time.Now()}
	m.mu.Unlock()
	m.cfg.Logger.Infow("secondary shard: cached trigger event for failover",
		"eventID", eventID,
		"myShardID", m.cfg.MyShardID,
		"workflowID", m.cfg.WorkflowID,
		"triggerIndex", event.TriggerIndex)
}

func (m *ShardFailoverManager) popCachedEvent(eventID string) (v2.RoutedTriggerEvent, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	ce, ok := m.cache[eventID]
	if !ok {
		return v2.RoutedTriggerEvent{}, false
	}
	delete(m.cache, eventID)
	return ce.event, true
}

func (m *ShardFailoverManager) removeCachedEvent(eventID string) {
	m.mu.Lock()
	delete(m.cache, eventID)
	m.mu.Unlock()
}

func (m *ShardFailoverManager) pruneLoop(ctx context.Context) {
	ticker := time.NewTicker(cachedExpiry / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.mu.Lock()
			now := time.Now()
			for id, ce := range m.cache {
				if now.Sub(ce.cachedAt) > cachedExpiry {
					delete(m.cache, id)
				}
			}
			m.mu.Unlock()
		}
	}
}

// wireFailover sets up the communicator for ExecutionStatusUpdate messages.
// It resolves both the primary and secondary shard DONs and configures the
// communicator with them. The handler is always registered (harmless when
// primary — nobody sends to us). Shard ownership is then resolved per-event
// in admissionCheck and forwardExecutionStatus, so no re-wiring is needed
// if the orchestrator reassigns roles at runtime.
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
	secondaryDon := m.resolveDon(ctx, m.cfg.WorkflowID, m.cfg.WorkflowOwner, 1)

	if primaryDon == nil && secondaryDon == nil {
		m.cfg.Logger.Warnw("shard failover: no shard DONs found, skipping communicator wiring", "myShardID", m.cfg.MyShardID)
		return nil
	}

	if primaryDon != nil && secondaryDon != nil {
		m.cfg.Communicator.SetShardDons(*primaryDon, *secondaryDon)
	} else if primaryDon != nil {
		m.cfg.Communicator.SetShardDons(*primaryDon, commoncap.DON{})
	} else {
		m.cfg.Communicator.SetShardDons(commoncap.DON{}, *secondaryDon)
	}

	m.cfg.Communicator.RegisterHandler(m.cfg.WorkflowID, m.HandleExecutionStatusUpdate)
	m.cfg.Logger.Infow("shard failover: wired communicator",
		"myShardID", m.cfg.MyShardID,
		"primaryDonID", primaryDonIDOrZero(primaryDon),
		"secondaryDonID", secondaryDonIDOrZero(secondaryDon))

	if err := m.cfg.Communicator.Start(ctx); err != nil {
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

func secondaryDonIDOrZero(d *commoncap.DON) uint32 {
	if d == nil {
		return 0
	}
	return d.ID
}

func (m *ShardFailoverManager) isPrimaryShard(ctx context.Context, workflowID, owner string) bool {
	if m.cfg.ShardResolver == nil {
		return m.cfg.MyShardID == 0
	}
	allResolver, ok := m.cfg.ShardResolver.(shardownership.AllShardsResolver)
	if !ok {
		return m.cfg.MyShardID == 0
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	shards, found, err := allResolver.ResolveAllShards(ctx, workflowID, owner)
	if err != nil || !found || len(shards) == 0 {
		return m.cfg.MyShardID == 0
	}
	return shards[0] == m.cfg.MyShardID
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

var _ DrainableService = (*ShardFailoverManager)(nil)
