package v2

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shopspring/decimal"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/smartcontractkit/chainlink-common/pkg/aggregation"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/services/orgresolver"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/cresettings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"
	sdkpb "github.com/smartcontractkit/chainlink-protos/cre/go/sdk"
	protoevents "github.com/smartcontractkit/chainlink-protos/workflows/go/events"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/safe"
)

var executingWorkflows atomic.Int64

var (
	ErrDuplicateExecution      = errors.New("duplicate execution")
	ErrShardDeniedNotOwner     = errors.New("shard ownership denied: not owner")
	ErrShardDeniedOrchestrator = errors.New("shard ownership denied: orchestrator error")
	ErrMeteringReserveFailed   = errors.New("metering reserve failed")
	ErrAdmissionCache          = errors.New("admission: event cached for failover")

	ErrEngineDraining = errors.New("engine is draining")
	ErrQueueFull      = errors.New("trigger event queue is full")
	ErrEnqueueFailed  = errors.New("failed to enqueue trigger event")
)

// Pin config version to 1 to avoid updating forwarder contracts on every single config update.
// Config Version set in CapabilitiesRegistry is included in every report but is irrelevant
// to validation on the forwarder side. What matters is DON ID and the set of signer public keys.
const pinnedWorkflowDonConfigVersion = 1

// baseEngine holds the execution machinery shared by every workflow engine:
// module execution, metering, secrets, labels, the heartbeat, drain state, and
// DON sync.
//
// It is never used directly. Engine (legacy, trigger-owning) and ExecutionEngine
// (execution-only) each embed it and supply their own start/init/close. The
// single services.Engine for a workflow engine lives here — see attachService.
type baseEngine struct {
	services.Service
	srvcEng *services.Engine

	cfg *EngineConfig

	// lggr is the engine's logger. It is protected by lggrMu to allow safe updates when DON configuration changes.
	// IMPORTANT: Do NOT access this field directly. Always use the logger() method to ensure thread-safety.
	// Direct access will cause data races when the logger is updated in localNodeSync().
	lggr logger.SugaredLogger

	// lggrMu protects lggr during dynamic updates when DON configuration changes.
	// Write access (Lock): Only in setLogger() called from localNodeSync()
	// Read access (RLock): In logger() method called from everywhere else
	lggrMu sync.RWMutex

	loggerLabels atomic.Pointer[map[string]string]
	localNode    atomic.Pointer[capabilities.Node]

	capCallsSemaphore limits.ResourcePoolLimiter[int]

	meterReports *metering.Reports

	metrics *monitoring.WorkflowsMetricLabeler

	// tracer is the OTel tracer for this engine. It's a noop tracer when DebugMode is false.
	tracer trace.Tracer

	orgID string
	// orgIDMissingReason records why orgID is empty ("resolver_nil",
	// "resolver_error", or "empty_response"). It is set once in startWith() and
	// read by startExecution() to label the org_id_missing counter. It is only
	// meaningful when orgID == "".
	orgIDMissingReason string

	draining         atomic.Bool
	activeExecutions atomic.Int32
	drainStartedAtNs atomic.Int64
}

func TriggerRegistrationID(workflowID string, triggerIndex int) string {
	return fmt.Sprintf("trigger_reg_%s_%d", workflowID, triggerIndex)
}

// buildLabels creates the label slice for the beholder logger based on config and localNode state.
// This is used both during engine creation and when updating labels after a DON configuration change.
func (e *baseEngine) buildLabels(localNode *capabilities.Node) []any {
	return []any{
		platform.KeyWorkflowID, e.cfg.WorkflowID,
		platform.KeyWorkflowOwner, e.cfg.WorkflowOwner,
		platform.KeyWorkflowName, e.cfg.WorkflowName.String(),
		platform.KeyWorkflowVersion, platform.ValueWorkflowVersionV2,
		platform.KeyDonID, strconv.Itoa(int(localNode.WorkflowDON.ID)),
		platform.KeyDonF, strconv.Itoa(int(localNode.WorkflowDON.F)),
		platform.KeyDonN, strconv.Itoa(len(localNode.WorkflowDON.Members)),
		platform.KeyDonQ, strconv.Itoa(aggregation.ByzantineQuorum(
			len(localNode.WorkflowDON.Members),
			int(localNode.WorkflowDON.F),
		)),
		platform.KeyP2PID, localNode.PeerID.String(),
		platform.WorkflowRegistryAddress, e.cfg.WorkflowRegistryAddress,
		platform.WorkflowRegistryChainSelector, e.cfg.WorkflowRegistryChainSelector,
		platform.EngineVersion, platform.ValueWorkflowVersionV2,
		platform.DonVersion, strconv.FormatUint(uint64(pinnedWorkflowDonConfigVersion), 10),
		platform.KeySDK, e.cfg.SdkName,
	}
}

// eventLabels returns a copy of the current label map with org ID applied.
func (e *baseEngine) eventLabels() map[string]string {
	return maps.Clone(*e.loggerLabels.Load())
}

// storeLoggerLabels persists base labels and always merges in the resolved org ID.
func (e *baseEngine) storeLoggerLabels(base map[string]string) {
	labels := maps.Clone(base)
	if e.orgID != "" {
		labels[platform.KeyOrganizationID] = e.orgID
	}
	e.loggerLabels.Store(&labels)
}

// logger returns the current logger in a thread-safe manner.
// This method should be used instead of accessing e.lggr directly to avoid race conditions
// when the logger is dynamically updated (e.g., when DON configuration changes).
func (e *baseEngine) logger() logger.SugaredLogger {
	e.lggrMu.RLock()
	defer e.lggrMu.RUnlock()
	return e.lggr
}

// setLogger updates the logger in a thread-safe manner.
// This is called when the DON configuration changes and we need to update the platform.DonVersion label.
func (e *baseEngine) setLogger(lggr logger.SugaredLogger) {
	e.lggrMu.Lock()
	defer e.lggrMu.Unlock()
	e.lggr = lggr
}

// Drain marks the engine as draining and prevents new executions from starting.
// In-flight executions continue to run to completion.
// It returns true only on the first transition to draining.
func (e *baseEngine) Drain() bool {
	started := e.draining.CompareAndSwap(false, true)
	if started {
		e.drainStartedAtNs.CompareAndSwap(0, time.Now().UnixNano())
	}
	e.srvcEng.SetHealthCond("draining", errors.New("engine is draining, pending deletion"))
	return started
}

func (e *baseEngine) ActiveExecutions() int32 {
	return e.activeExecutions.Load()
}

func (e *baseEngine) DrainStartedAt() (time.Time, bool) {
	ns := e.drainStartedAtNs.Load()
	if ns == 0 {
		return time.Time{}, false
	}

	return time.Unix(0, ns), true
}

// newBaseEngine builds the execution machinery shared by every workflow engine.
// It installs no service: the caller must call attachService before the engine
// is started.
//
// Start and Close are supplied by the outer type, because the lifecycle differs
// per engine, but srvcEng itself lives here — base methods spawn the execution
// goroutines on it, and Drain sets its health condition.
func newBaseEngine(cfg *EngineConfig) (*baseEngine, logger.SugaredLogger, error) {
	if err := cfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("invalid config: %w", err)
	}

	em, err := monitoring.InitMonitoringResources()
	if err != nil {
		return nil, nil, fmt.Errorf("could not initialize monitoring resources: %w", err)
	}

	// LocalNode() is expected to be non-blocking at this stage (i.e. the registry is already synced)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.LocalLimits.LocalNodeTimeoutMs)*time.Millisecond)
	defer cancel()
	localNode, err := cfg.CapRegistry.LocalNode(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get local node state: %w", err)
	}

	// Create engine first so we can use the buildLabels method
	engine := &baseEngine{
		cfg:               cfg,
		capCallsSemaphore: cfg.LocalLimiters.CapabilityConcurrency,
	}

	// Build labels using the helper method
	labels := engine.buildLabels(&localNode)

	beholderLogger := logger.Sugared(custmsg.NewBeholderLogger(cfg.Lggr, cfg.BeholderEmitter).Named("WorkflowEngine").With(labels...))
	baseLabels := []string{
		platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String(),
		platform.KeySDK, cfg.SdkName,
	}
	metricsLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em).With(baseLabels...)
	labelsMap := make(map[string]string, len(labels)/2)
	for i := 0; i < len(labels); i += 2 {
		labelsMap[labels[i].(string)] = labels[i+1].(string)
	}

	if cfg.DebugMode {
		beholderLogger.Errorw("WARNING: Debug mode is enabled, this is not suitable for production")
		engine.tracer = otel.Tracer("workflow_engine_v2")
	} else {
		engine.tracer = noop.NewTracerProvider().Tracer("")
	}

	// Store logger and other fields
	engine.setLogger(beholderLogger)
	engine.meterReports = metering.NewReports(cfg.BillingClient, cfg.WorkflowOwner, cfg.WorkflowID, beholderLogger, labelsMap, metricsLabeler, cfg.WorkflowRegistryAddress, cfg.WorkflowRegistryChainSelector, metering.EngineVersionV2)
	engine.metrics = metricsLabeler
	engine.loggerLabels.Store(&labelsMap)
	engine.localNode.Store(&localNode)

	return engine, beholderLogger, nil
}

// attachService installs the single services.Engine for this workflow engine.
// start and close belong to the outer type that owns the lifecycle.
func (e *baseEngine) attachService(lggr logger.SugaredLogger, start func(context.Context) error, closeFn func() error) {
	e.Service, e.srvcEng = services.Config{
		Name:  "WorkflowEngineV2",
		Start: start,
		Close: closeFn,
	}.NewServiceEngine(lggr)
}

// resolvedOrg holds the result of an organization ID resolution attempt.
type resolvedOrg struct {
	// ID is the resolved organization ID, or empty if resolution failed.
	ID string
	// Err is the error returned by the OrgResolver, or nil if resolution
	// succeeded or the resolver was not configured.
	Err error
	// Reason explains why ID is empty: "resolver_nil" (OrgResolver not
	// configured), "resolver_error" (Get returned an error), or
	// "empty_response" (Get returned an empty string). Empty on success.
	Reason string
}

// ExecuteTrigger is the engine's single execution entry point. It performs no admission control, the caller is responsible for those.
func (e *baseEngine) ExecuteTrigger(ctx context.Context, event RoutedTriggerEvent) error {
	e.activeExecutions.Add(1)
	defer e.activeExecutions.Add(-1)

	eventID := event.Event.Event.ID
	e.logger().Debugw("Scheduling a trigger event for execution", "eventID", eventID)
	creCtx := contexts.CREValue(ctx)
	// Tracer is no-op if DebugMode is false
	ctx, span := e.tracer.Start(ctx, "workflow_execution",
		trace.WithAttributes(
			attribute.String("workflow_name", e.cfg.WorkflowName.String()),
			attribute.String("version", "v2"),
			attribute.String("org_id", creCtx.Org),
			attribute.String("owner_id", creCtx.Owner),
			attribute.String("workflow_id", creCtx.Workflow),
		))
	defer span.End()

	return e.startExecution(ctx, event)
}

// Subscribe issues the WASM Subscribe request and returns the validated trigger subscriptions.
func (e *baseEngine) Subscribe(ctx context.Context) ([]*sdkpb.TriggerSubscription, error) {
	// call into the workflow to get trigger subscriptions
	subCtx, subCancel, err := e.cfg.LocalLimiters.TriggerSubscriptionTime.WithTimeout(ctx)
	if err != nil {
		return nil, err
	}
	defer subCancel()

	maxUserLogEventsPerExecution, err := e.cfg.LocalLimiters.LogEvent.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			return nil, err
		}
		e.logger().Errorw("Failed to get log event limit; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.LogEventLimit.Key)
	}
	userLogChan := make(chan *protoevents.LogLine, maxUserLogEventsPerExecution)
	defer close(userLogChan)
	e.srvcEng.GoCtx(subCtx, func(ctx context.Context) {
		e.emitUserLogs(ctx, userLogChan, e.cfg.WorkflowID, e.eventLabels())
	})

	var timeProvider TimeProvider = &types.LocalTimeProvider{}
	if !e.cfg.UseLocalTimeProvider {
		timeProvider = NewDonTimeProvider(e.cfg.DonTimeStore, e.cfg.WorkflowID, e.donTimeRequestTimeout(subCtx, e.cfg.LocalLimiters.DONTimeRequestTimeout), e.logger(), e.metrics, e.cfg.Clock)
	}

	moduleExecuteMaxResponseSizeBytes, err := e.cfg.LocalLimiters.ExecutionResponse.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			return nil, err
		}
		e.logger().Errorw("Failed to get execution response size limit; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.ExecutionResponseLimit.Key)
	}
	if moduleExecuteMaxResponseSizeBytes < 0 {
		return nil, fmt.Errorf("invalid moduleExecuteMaxResponseSizeBytes; must not be negative: %d", moduleExecuteMaxResponseSizeBytes)
	}
	result, err := e.cfg.Module.Execute(subCtx, &sdkpb.ExecuteRequest{
		Request:         &sdkpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(moduleExecuteMaxResponseSizeBytes),
		Config:          e.cfg.WorkflowConfig,
	}, NewDisallowedExecutionHelper(e.logger(), userLogChan, timeProvider, e.secretsFetcher(e.cfg.WorkflowID)))
	if err != nil {
		return nil, fmt.Errorf("failed to execute subscribe: %w", err)
	}
	if result.GetError() != "" {
		return nil, fmt.Errorf("failed to execute subscribe: %s", result.GetError())
	}
	subs := result.GetTriggerSubscriptions()
	if subs == nil {
		return nil, errors.New("subscribe result is nil")
	}
	err = e.cfg.LocalLimiters.TriggerSubscription.Check(ctx, len(subs.Subscriptions))
	if err != nil {
		return nil, err
	}

	return subs.Subscriptions, nil
}

// Draining returns true if the engine has been marked for deletion and is no longer accepting new trigger events.
func (e *baseEngine) Draining() bool {
	return e.draining.Load()
}

// resolveOrgID resolves the organization ID for the given workflow owner.
// If resolution fails, the returned ID is empty and Reason explains why.
// The original error from the resolver (if any) is preserved in Err and
// logged via the provided logger.
func resolveOrgID(ctx context.Context, resolver orgresolver.OrgResolver, workflowOwner string, lggr logger.SugaredLogger) resolvedOrg {
	if resolver == nil {
		return resolvedOrg{Reason: "resolver_nil"}
	}
	orgID, err := resolver.Get(ctx, workflowOwner)
	if err != nil {
		lggr.Warnw("Failed to resolve organization ID, continuing without it", "workflowOwner", workflowOwner, "err", err)
		return resolvedOrg{Err: err, Reason: "resolver_error"}
	}
	if orgID == "" {
		return resolvedOrg{Reason: "empty_response"}
	}
	return resolvedOrg{ID: orgID}
}

// startWith performs the startup shared by every engine and spawns initFn as the
// initialization goroutine. Each engine passes its own init.
//
// triggerLoopFn is Engine's queue-draining loop (handleAllTriggerEvents). It is
// legacy-only: ExecutionEngine has no queue to drain, since its future
// coordinator calls ExecuteTrigger directly instead of going through Put. Pass
// nil to skip it.
func (e *baseEngine) startWith(ctx context.Context, initFn func(context.Context), triggerLoopFn func(context.Context)) error {
	e.cfg.Module.Start()
	ctx = context.WithoutCancel(ctx)

	// Resolve the workflow owner's org once at engine startup and treat it as stable
	// for the lifetime of this engine instance. If org membership/linking changes, the
	// workflow must be restarted to pick up the new org mapping.
	resolved := resolveOrgID(ctx, e.cfg.OrgResolver, e.cfg.WorkflowOwner, e.logger())
	e.orgID = resolved.ID
	e.orgIDMissingReason = resolved.Reason
	e.storeLoggerLabels(e.eventLabels())

	e.metrics = e.metrics.With(platform.KeyOrganizationID, e.orgID)

	ctx = contexts.WithCRE(ctx, contexts.CRE{Org: e.orgID, Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID})
	e.srvcEng.GoCtx(ctx, e.heartbeatLoop)
	e.srvcEng.GoCtx(ctx, initFn)
	if triggerLoopFn != nil {
		e.srvcEng.GoCtx(ctx, triggerLoopFn)
	}
	return nil
}

// initDONSubscribe subscribes to the DON notifier and starts the local-node sync
// loop. A returned error has already been logged; the caller passes it to OnInitialized.
func (e *baseEngine) initDONSubscribe(ctx context.Context) error {
	donSubCh, cleanup, err := e.cfg.DonSubscriber.Subscribe(ctx)
	if err != nil {
		e.logger().Errorw("failed to subscribe to DON notifier", "error", err)
		return fmt.Errorf("failed to subscribe to DON notifier: %w", err)
	}

	// start loop to sync local node state each time a DON is received on the
	// subscribed channel
	e.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
		defer cleanup()
		for {
			select {
			case <-ctx.Done():
				return
			case _, open := <-donSubCh:
				if !open {
					return
				}
				e.localNodeSync(ctx)
			}
		}
	})
	return nil
}

// initSubscriptions runs the WASM Subscribe call and hands the validated
// subscriptions to the OnSubscriptionsReady hook. A returned error has already
// been logged; the caller passes it to OnInitialized.
func (e *baseEngine) initSubscriptions(ctx context.Context) ([]*sdkpb.TriggerSubscription, error) {
	subscriptions, err := e.Subscribe(ctx)
	if err != nil {
		e.logger().Errorw("failed to subscribe to triggers", "err", err)
		return nil, err
	}

	cre := contexts.CRE{Org: e.orgID, Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID}
	if err = e.cfg.Hooks.OnSubscriptionsReady(subscriptions, cre); err != nil {
		e.logger().Errorw("OnSubscriptionsReady hook failed", "err", err)
		return nil, err
	}
	return subscriptions, nil
}

// initDone records a successful initialization and fires OnInitialized(nil).
func (e *baseEngine) initDone(ctx context.Context) {
	e.logger().Info("Workflow Engine initialized")
	e.metrics.IncrementWorkflowInitializationCounter(ctx)
	e.cfg.Hooks.OnInitialized(nil)
}

// shutdownCtx builds the close context: bounded by the shutdown timeout and
// carrying the workflow's tenant identity.
func (e *baseEngine) shutdownCtx() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(e.cfg.LocalLimits.ShutdownTimeoutMs))
	return contexts.WithCRE(ctx, contexts.CRE{Org: e.orgID, Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID}), cancel
}

// closeCommon is the teardown shared by every engine.
func (e *baseEngine) closeCommon(ctx context.Context) {
	if err := e.cfg.ExecutionsStore.DeleteByWorkflowID(ctx, e.cfg.WorkflowID); err != nil {
		e.logger().Errorw("Failed to purge executions on close", "err", err)
	}

	e.cfg.Module.Close()

	if e.cfg.LocalLimiters != nil {
		if err := e.cfg.LocalLimiters.EvictWorkflow(e.cfg.WorkflowID); err != nil {
			e.logger().Errorw("Failed to evict workflow from scoped limiters", "err", err)
		}
	}

	// Encourage the Go runtime to release memory back to the OS after tearing
	// down the WASM module and execution state.  Without this, freed heap pages
	// stay resident (MADV_FREE) and CGo/wasmtime freed pages remain in the C
	// allocator's free-list, so RSS never drops even though the memory is unused.
	runtime.GC()
	debug.FreeOSMemory()

	// reset metering mode metric so that a positive value does not persist
	e.metrics.UpdateWorkflowMeteringModeGauge(ctx, false)
}

func (e *baseEngine) localNodeSync(ctx context.Context) {
	to := time.Duration(e.cfg.LocalLimits.LocalNodeTimeoutMs) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, to)
	defer cancel()

	localNode, err := e.cfg.CapRegistry.LocalNode(ctx)
	if err != nil {
		e.cfg.Lggr.Errorf("could not get local node state: %s", err)
		e.cfg.Hooks.OnNodeSynced(localNode, err)
		return
	}

	// ignore any reads that do not update the config version
	if e.localNode.Load().WorkflowDON.ConfigVersion == localNode.WorkflowDON.ConfigVersion {
		return
	}

	e.cfg.Lggr.Debugw("Setting local node state",
		"Workflow DON ID", localNode.WorkflowDON.ID,
		"Workflow DON Families", localNode.WorkflowDON.Families,
		"Workflow DON Config Version (onchain)", localNode.WorkflowDON.ConfigVersion,
		"Workflow DON Config Version (pinned)", pinnedWorkflowDonConfigVersion,
	)

	// Publish the new node before updating logger state so concurrent executions
	// observe the synced DON while labels are rebuilt.
	e.localNode.Store(&localNode)

	// Recreate the beholder logger with updated labels to reflect the new DON version.
	labels := e.buildLabels(&localNode)
	newLogger := logger.Sugared(
		custmsg.NewBeholderLogger(e.cfg.Lggr, e.cfg.BeholderEmitter).
			Named("WorkflowEngine").
			With(labels...),
	)
	e.setLogger(newLogger)

	labelsMap := make(map[string]string, len(labels)/2)
	for i := 0; i < len(labels); i += 2 {
		labelsMap[labels[i].(string)] = labels[i+1].(string)
	}
	e.storeLoggerLabels(labelsMap)

	e.cfg.Hooks.OnNodeSynced(localNode, nil)
}

// startExecution initiates a new workflow execution, blocking until completed
func (e *baseEngine) startExecution(ctx context.Context, event RoutedTriggerEvent) error {
	triggerDrop := func(reason string) {
		e.metrics.With(platform.KeyTriggerID, event.TriggerCapID).IncrementTriggerEventDroppedTotal(ctx, reason)
	}

	executionID, err := workflows.GenerateExecutionIDWithTriggerIndex(e.cfg.WorkflowID, event.Event.Event.ID, event.TriggerIndex)
	if err != nil {
		e.logger().Errorw("Failed to generate execution ID", "err", err, "triggerID", event.TriggerCapID)
		triggerDrop(monitoring.TriggerDropReasonExecutionIDGenerationFailed)
		return err
	}
	e.metrics.IncrementExecutionIDFullCounter(ctx)
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("execution_id", executionID))

	loggerLabels := e.eventLabels()
	lggr := e.logger().With(platform.KeyOrganizationID, e.orgID)

	var executionTimestamp time.Time
	var executionDonTimeProvider TimeProvider
	if tsErr := e.cfg.LocalLimiters.ExecutionTimestampsEnabled.AllowErr(ctx); tsErr == nil {
		executionDonTimeProvider = NewDonTimeProvider(e.cfg.DonTimeStore, executionID, e.donTimeRequestTimeout(ctx, e.cfg.LocalLimiters.DONTimeRequestTimeout), lggr, e.metrics, e.cfg.Clock)
		donTime, dtErr := executionDonTimeProvider.GetDONTime()
		if dtErr != nil {
			executionTimestamp = e.cfg.Clock.Now()
			lggr.Warnw("Failed to get DON time for execution timestamp, falling back to local time", "err", dtErr, "executionTimestamp", executionTimestamp)
			e.metrics.IncrementExecutionTimestampFallbackCounter(ctx)
		} else {
			executionTimestamp = donTime
			lggr.Debugw("Execution timestamp assigned", "executionTimestamp", executionTimestamp)
			e.metrics.IncrementExecutionTimestampAssignedCounter(ctx)
		}
	}

	triggerEvent := event.Event.Event

	// disallow duplicate executions
	_, addErr := e.cfg.ExecutionsStore.Add(ctx, nil, executionID, e.cfg.WorkflowID, store.StatusStarted)
	if addErr != nil {
		if errors.Is(addErr, store.ErrDuplicateExecution) {
			lggr.Infow("Skipping duplicate execution", "executionID", executionID, "triggerID", event.TriggerCapID, "triggerIndex", event.TriggerIndex)
			tm := e.metrics.With(platform.KeyTriggerID, event.TriggerCapID)
			tm.IncrementTriggerExecutionDeduplicatedCounter(ctx)
			tm.IncrementWorkflowTriggerEventErrorCounter(ctx)
			tm.IncrementTriggerEventDroppedTotal(ctx, monitoring.TriggerDropReasonDuplicateExecution)
			registrationID := TriggerRegistrationID(e.cfg.WorkflowID, event.TriggerIndex)
			ackErr := e.cfg.TriggerAcknowledger.Ack(ctx, event.TriggerCapID, registrationID, triggerEvent.ID)
			if ackErr != nil {
				e.lggr.Errorw("failed to re-ACK trigger event", "eventID", triggerEvent.ID, "err", ackErr)
			}
			return ErrDuplicateExecution
		}
		lggr.Errorw("Failed to register execution in store, proceeding anyway", "executionID", executionID, "err", addErr)
	}

	var executionStatus string
	defer func() {
		if executionStatus == "" {
			executionStatus = store.StatusErrored
		}
		if _, finishErr := e.cfg.ExecutionsStore.FinishExecution(ctx, executionID, executionStatus); finishErr != nil {
			lggr.Errorw("Failed to finish execution in store", "executionID", executionID, "status", executionStatus, "err", finishErr)
		}
	}()

	// emitDroppedExecution publishes the Started/Finished pair for an execution abandoned
	// before the normal Started/Finished emit points below, so the failure reaches the UI
	// instead of vanishing.
	emitDroppedExecution := func(cause error, class events.ErrorClassification) {
		executionStatus = store.StatusErrored
		_ = events.EmitExecutionStartedEvent(ctx, loggerLabels, triggerEvent.ID, executionID)
		_ = events.EmitExecutionFinishedEvent(ctx, loggerLabels, store.StatusErrored, executionID, cause, class, lggr)
		e.metrics.IncrementWorkflowExecutionFinishedCounter(ctx, store.StatusErrored)
	}

	e.metrics.UpdateTotalWorkflowsGauge(ctx, executingWorkflows.Add(1))
	defer e.metrics.UpdateTotalWorkflowsGauge(ctx, executingWorkflows.Add(-1))

	// TODO(CAPPL-911): add rate-limiting

	meteringReport, meteringErr := e.meterReports.Start(ctx, executionID)
	if meteringErr != nil {
		lggr.Errorw("could start metering workflow execution. continuing without metering", "err", meteringErr)
	}

	isMetering := meteringErr == nil
	if isMetering {
		mrErr := meteringReport.Reserve(ctx)
		if mrErr != nil {
			lggr.Errorw("could not reserve metering", "err", mrErr)
			triggerDrop(monitoring.TriggerDropReasonMeteringReserveFailed)
			return ErrMeteringReserveFailed
		}

		e.deductStandardBalances(ctx, meteringReport)
	}

	// WithTimeout returns a usable ctx/cancel even on a read failure; err is advisory.
	execCtx, execCancel, err := e.cfg.LocalLimiters.ExecutionTime.WithTimeout(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			lggr.Errorw("Failed to get execution time limit with no usable value", "err", err)
			triggerDrop(monitoring.TriggerDropReasonExecutionTimeLimitReadFailed)
			emitDroppedExecution(err, events.ErrorClassificationSystem)
			if execCancel != nil {
				execCancel() // WithTimeout may still have built a context
			}
			return err
		}
		lggr.Errorw("Failed to get execution time limit; continuing with the timeout the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.ExecutionTimeout.Key)
		if execCtx == nil { // only nil when the limiter is closed and no ctx was built
			execCtx, execCancel = context.WithTimeout(ctx, cresettings.Default.PerWorkflow.ExecutionTimeout.DefaultValue)
		}
	}
	defer execCancel()
	triggerCapID := event.TriggerCapID
	skewRec := &monitoring.TriggerSkewRecorder{
		EnqueueTime: event.ObservedAt,
		Record: func(ctx context.Context, seconds float64, source string) {
			e.metrics.With(platform.KeyTriggerID, triggerCapID).RecordTriggerQueueToExecutionStartSeconds(ctx, seconds, source)
		},
	}
	execCtx = monitoring.ContextWithTriggerSkewRecorder(execCtx, skewRec)
	executionLogger := logger.With(lggr, "executionID", executionID, "triggerID", event.TriggerCapID,
		"triggerIndex", event.TriggerIndex, "eventID", triggerEvent.ID)

	// This is only a peek to size the user-log channel's burst buffer; LogEvent.Check
	// (called per log line in emitUserLogs) is what actually enforces the cap.
	maxUserLogEventsPerExecution, err := e.cfg.LocalLimiters.LogEvent.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			lggr.Errorw("Failed to get log event limit with no usable value", "err", err)
			triggerDrop(monitoring.TriggerDropReasonLogEventLimitReadFailed)
			emitDroppedExecution(err, events.ErrorClassificationSystem)
			return err
		}
		lggr.Errorw("Failed to get log event limit; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.LogEventLimit.Key)
	}
	userLogChan := make(chan *protoevents.LogLine, maxUserLogEventsPerExecution)
	defer close(userLogChan)
	e.srvcEng.GoCtx(execCtx, func(ctx context.Context) {
		e.emitUserLogs(ctx, userLogChan, executionID, loggerLabels)
	})

	tid, err := safe.IntToUint64(event.TriggerIndex)
	if err != nil {
		executionLogger.Errorw("Failed to convert trigger index to uint64", "err", err)
		triggerDrop(monitoring.TriggerDropReasonTriggerIndexInvalid)
		emitDroppedExecution(err, events.ErrorClassificationSystem)
		return err
	}

	startTime := e.cfg.Clock.Now()
	executionLogger.Infow("Workflow execution starting ...")
	if e.orgID == "" {
		e.metrics.IncrementOrgIDMissingCounter(ctx, e.orgIDMissingReason)
	}
	_ = events.EmitExecutionStartedEvent(ctx, loggerLabels, triggerEvent.ID, executionID)

	registrationID := TriggerRegistrationID(e.cfg.WorkflowID, event.TriggerIndex)
	err = e.cfg.TriggerAcknowledger.Ack(ctx, event.TriggerCapID, registrationID, triggerEvent.ID)
	if err != nil {
		e.lggr.Errorf("failed to ACK trigger event (eventID=%s): %v", triggerEvent.ID, err)
	}
	e.metrics.With("workflowID", e.cfg.WorkflowID, "workflowName", e.cfg.WorkflowName.String()).IncrementWorkflowExecutionStartedCounter(ctx)

	// Track execution error (and its user/system classification) for deferred
	// event emission. Set alongside executionStatus at each outcome site below.
	var execErr error
	execErrClass := events.ErrorClassificationUnspecified
	var execHelper *ExecutionHelper
	defer func() {
		_ = events.EmitExecutionFinishedEvent(ctx, loggerLabels, executionStatus, executionID, execErr, execErrClass, lggr)
		if execHelper != nil {
			endTime := e.cfg.Clock.Now()
			profile, emitErr := events.EmitExecutionProfile(
				ctx,
				e.cfg.WorkflowID,
				executionID,
				startTime,
				endTime,
				executionStatus,
				execHelper.executionProfile.stepInputs(),
			)
			if emitErr != nil {
				lggr.Errorw("Failed to emit execution profile", "err", emitErr)
			}
			if profile != nil {
				profileJSON, jsonErr := protojson.Marshal(profile)
				if jsonErr != nil {
					lggr.Errorw("Failed to marshal execution profile to JSON", "err", jsonErr)
				} else {
					lggr.Infow("Workflow execution profile", "executionProfile", string(profileJSON))
				}
			}
		}
		e.cfg.Hooks.OnExecutionFinished(executionID, executionStatus)
		e.cfg.Hooks.OnExecutionStatusUpdate(e.cfg.WorkflowID, executionID, triggerEvent.ID, event.TriggerIndex, executionStatus, execErrClass)
		if execErr != nil {
			e.cfg.Hooks.OnExecutionError(execErr.Error())
		}
	}()

	var timeProvider TimeProvider = &types.LocalTimeProvider{}
	if !e.cfg.UseLocalTimeProvider {
		if executionDonTimeProvider != nil {
			timeProvider = executionDonTimeProvider
		} else {
			lggr.Warnw("ExecutionTimestampsEnabled is false - creating a new DON time provider")
			timeProvider = NewDonTimeProvider(e.cfg.DonTimeStore, executionID, e.donTimeRequestTimeout(execCtx, e.cfg.LocalLimiters.DONTimeRequestTimeout), lggr, e.metrics, e.cfg.Clock)
		}
	}

	// Track time the guest spends blocked in DON-time host calls so it can be
	// excluded from metered compute (see compute_metering.go). The wrapper is
	// always installed; whether the recorded wait is subtracted is gated by
	// e.cfg.MeterComputeExcludingHostWait below.
	suspension := &suspensionTracker{}
	timeProvider = newMeasuredTimeProvider(timeProvider, e.cfg.Clock, suspension)

	// Limit is always usable even on a read failure; err is advisory.
	moduleExecuteMaxResponseSizeBytes, err := e.cfg.LocalLimiters.ExecutionResponse.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			execErr = fmt.Errorf("failed to get execution response size limit with no usable value: %w", err)
			lggr.Errorw(execErr.Error())
			executionStatus = store.StatusErrored
			execErrClass = events.ErrorClassificationSystem
			triggerDrop(monitoring.TriggerDropReasonExecutionResponseLimitReadFailed)
			return execErr
		}
		lggr.Errorw("Failed to get execution response size limit; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.ExecutionResponseLimit.Key)
	}
	if moduleExecuteMaxResponseSizeBytes < 0 {
		execErr = fmt.Errorf("invalid moduleExecuteMaxResponseSizeBytes; must not be negative: %d", moduleExecuteMaxResponseSizeBytes)
		lggr.Errorw(execErr.Error())
		executionStatus = store.StatusErrored
		execErrClass = events.ErrorClassificationSystem
		triggerDrop(monitoring.TriggerDropReasonExecutionResponseLimitInvalid)
		return execErr
	}
	execHelper = &ExecutionHelper{
		baseEngine: e, WorkflowExecutionID: executionID, ExecutionTimestamp: executionTimestamp,
		UserLogChan: userLogChan, TimeProvider: timeProvider, SecretsFetcher: e.secretsFetcher(executionID),
		executionProfile: newExecutionProfileCollector(),
		suspension:       suspension,
	}

	execHelper.initLimiters(e.cfg.LocalLimiters)
	e.metrics.With(platform.KeyTriggerID, event.TriggerCapID).RecordTriggerPayloadBytes(ctx, int64(proto.Size(triggerEvent.Payload)))
	var result *sdkpb.ExecutionResult
	result, execErr = e.cfg.Module.Execute(execCtx, &sdkpb.ExecuteRequest{
		Request: &sdkpb.ExecuteRequest_Trigger{
			Trigger: &sdkpb.Trigger{
				Id:      tid,
				Payload: triggerEvent.Payload,
			},
		},
		MaxResponseSize: uint64(moduleExecuteMaxResponseSizeBytes),
		Config:          e.cfg.WorkflowConfig,
	}, execHelper.PossiblyWithRawSecrets())
	// Non-evictable modules do not record skew; label those as direct.
	skewRec.RecordReady(execCtx, monitoring.ModuleLoadSourceDirect)

	endTime := e.cfg.Clock.Now()
	executionDuration := endTime.Sub(startTime)

	// computeDuration is what we meter as RESOURCE_TYPE_COMPUTE. We start from the
	// end-to-end wall-clock and subtract time the guest spent blocked in DON-time
	// host calls: that wait is not compute and is non-deterministic across nodes
	// (a node whose DON-time request lands in the current consensus round returns
	// in microseconds; one whose request slips to the next round blocks ~a full
	// round). Operational latency metrics below intentionally keep using
	// executionDuration.
	computeDuration := executionDuration
	if hostWait := suspension.total(); hostWait > 0 {
		computeDuration -= hostWait
		if computeDuration < 0 {
			computeDuration = 0
		}
		executionLogger.Debugw("Excluding DON-time wait from metered compute",
			"wallClockMs", executionDuration.Milliseconds(),
			"hostWaitMs", hostWait.Milliseconds(),
			"computeMs", computeDuration.Milliseconds())
	}

	if isMetering {
		computeUnit := billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_COMPUTE)]
		mrErr := meteringReport.Settle(computeUnit,
			capabilities.ResponseMetadata{
				Metering: []capabilities.MeteringNodeDetail{{
					Peer2PeerID: e.localNode.Load().PeerID.String(),
					SpendUnit:   computeUnit,
					SpendValue:  strconv.Itoa(int(computeDuration.Milliseconds())),
				}},
				CapDON_N: 1,
			},
		)
		if mrErr != nil {
			lggr.Errorw("could not set metering for compute", "err", mrErr)
		}
		mrErr = e.meterReports.End(ctx, executionID)
		if mrErr != nil {
			lggr.Errorw("could not end metering report", "err", mrErr)
		}
	}

	if execErr != nil {
		executionStatus = store.StatusErrored
		// Module/host and timeout failures are platform errors by default, but a
		// user-origin caperrors.Error propagating from a capability or the guest
		// is attributed to the user.
		execErrClass = events.ClassifyError(execErr, events.ErrorClassificationSystem)
		if errors.Is(execErr, context.DeadlineExceeded) {
			executionStatus = store.StatusTimeout
			e.metrics.UpdateWorkflowTimeoutDurationHistogram(ctx, int64(executionDuration.Seconds()))
		} else {
			e.metrics.UpdateWorkflowErrorDurationHistogram(ctx, int64(executionDuration.Seconds()))
		}
		e.metrics.IncrementWorkflowExecutionFinishedCounter(ctx, executionStatus)
		executionLogger.Errorw("Workflow execution failed with module execution error", "status", executionStatus, "durationMs", executionDuration.Milliseconds(), "err", execErr)
		return nil // not an error from the caller's perspective. Execution ran, just failed. This is already captured by the deferred lifecycle hook e.cfg.Hooks.OnExecutionError(execErr.Error())
	}

	if e.cfg.DebugMode {
		lggr.Debugw("User workflow execution result", "result", result.GetValue(), "err", result.GetError())
	}

	if len(result.GetError()) > 0 {
		executionStatus = store.StatusErrored
		execErr = errors.New(result.GetError())
		// The user's workflow ran and returned an error: a user failure.
		execErrClass = events.ErrorClassificationUser
		e.metrics.UpdateWorkflowErrorDurationHistogram(ctx, int64(executionDuration.Seconds()))
		e.metrics.With("workflowID", e.cfg.WorkflowID, "workflowName", e.cfg.WorkflowName.String()).IncrementWorkflowExecutionFailedCounter(ctx)
		e.metrics.IncrementWorkflowExecutionFinishedCounter(ctx, executionStatus)
		executionLogger.Errorw("Workflow execution failed", "status", executionStatus, "durationMs", executionDuration.Milliseconds(), "error", result.GetError())
		return nil // not an error from the caller's perspective. Execution ran, just returned an error. This is already captured by the deferred lifecycle hook e.cfg.Hooks.OnExecutionError(execErr.Error())
	}

	executionStatus = store.StatusCompleted
	executionLogger.Infow("Workflow execution finished successfully", "durationMs", executionDuration.Milliseconds())
	e.metrics.UpdateWorkflowCompletedDurationHistogram(ctx, int64(executionDuration.Seconds()))
	e.metrics.With("workflowID", e.cfg.WorkflowID, "workflowName", e.cfg.WorkflowName.String()).IncrementWorkflowExecutionSucceededCounter(ctx)
	e.metrics.IncrementWorkflowExecutionFinishedCounter(ctx, executionStatus)
	e.cfg.Hooks.OnResultReceived(result)
	return nil
}

func (e *baseEngine) secretsFetcher(phaseID string) SecretsFetcher {
	if e.cfg.SecretsFetcher != nil {
		return e.cfg.SecretsFetcher
	}

	return NewSecretsFetcher(
		e.metrics,
		e.cfg.CapRegistry,
		e.logger(),
		e.cfg.LocalLimiters.SecretsConcurrency,
		e.cfg.LocalLimiters.SecretsCalls,
		e.cfg.LocalLimiters.Settings,
		e.orgID,
		e.cfg.WorkflowOwner,
		e.cfg.WorkflowName.String(),
		e.cfg.WorkflowID,
		// phaseID is the executionID if called during an execution,
		// or the workflowID if called during trigger subscription
		phaseID,
		e.cfg.WorkflowEncryptionKey,
		e.cfg.OverrideFetcher,
	)
}

func (e *baseEngine) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(e.cfg.LocalLimits.HeartbeatFrequencyMs) * time.Millisecond)
	defer ticker.Stop()
	e.logger().Info("Starting heartbeat loop")
	e.metrics.EngineHeartbeatGauge(ctx, 1)

	for {
		select {
		case <-ctx.Done():
			e.metrics.EngineHeartbeatGauge(ctx, 0)
			e.logger().Info("Shutting down heartbeat")
			return
		case <-ticker.C:
			e.logger().Debugw("Engine heartbeat tick", "time", e.cfg.Clock.Now().Format(time.RFC3339))
			e.metrics.IncrementEngineHeartbeatCounter(ctx)
		}
	}
}

func (e *baseEngine) deductStandardBalances(ctx context.Context, meteringReport *metering.Report) {
	// V2Engine runs the entirety of a module's execution as compute. Ensure that the max execution time can run.
	// Add an extra second of metering padding for context cancel propagation
	ctxCancelPadding := (time.Millisecond * 1000).Milliseconds()
	// Limit is always usable even on a read failure; err is advisory.
	workflowExecutionTimeout, err := e.cfg.LocalLimiters.ExecutionTime.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			e.logger().Errorw("Failed to get execution time limit with no usable value; skipping compute deduction", "err", err)
			return
		}
		e.logger().Errorw("Failed to get execution time limit; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.ExecutionTimeout.Key)
	}
	compMs := decimal.NewFromInt(workflowExecutionTimeout.Milliseconds() + ctxCancelPadding)
	computeUnit := billing.ResourceType_RESOURCE_TYPE_COMPUTE.String()

	if _, err := meteringReport.Deduct(
		computeUnit,
		metering.ByResource(computeUnit, "v2-standard-deduction-compute", compMs),
	); err != nil {
		e.logger().Errorw("could not deduct balance for capability request", "capReq", "standard-deduction-compute", "err", err)
	}
}

const emitUserLogsTimeout = 30 * time.Second

// separate call for each workflow execution
func (e *baseEngine) emitUserLogs(ctx context.Context, userLogChan chan *protoevents.LogLine, executionID string, executionLabels map[string]string) {
	e.logger().Debugw("Listening for user logs ...")
	count := 0
	defer func() { e.logger().Debugw("Listening for user logs done.", "processedLogLines", count) }()

	processLogLine := func(emitCtx context.Context, logLine *protoevents.LogLine) bool {
		if e.cfg.DebugMode {
			e.logger().Debugf("User log: <<<%s>>>, local node timestamp: %s", logLine.Message, logLine.NodeTimestamp)
		}
		if err := e.cfg.LocalLimiters.LogEvent.Check(emitCtx, count); err != nil {
			if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[int]](err); ok {
				e.logger().Warnw("Max user log events per execution reached, dropping event", "maxEvents", errBoundLimited.Limit, "err", err)
				return false
			}
			// A settings read failure should not stop the drain. Fail open instead.
			if limits.IsErrRecoverable(err) {
				e.logger().Errorw("Failed to check user log event limit; emitting anyway", "err", err)
			} else {
				e.logger().Errorw("User log event limit could not be evaluated; emitting anyway", "err", err)
			}
			e.metrics.IncrementLimitCheckUnenforcedCounter(emitCtx, cresettings.Default.PerWorkflow.LogEventLimit.Key)
		}
		if err := e.cfg.LocalLimiters.LogLine.Check(emitCtx, config.Size(len(logLine.Message))); err != nil {
			if errBoundLimited, ok := errors.AsType[limits.ErrorBoundLimited[config.Size]](err); ok {
				logLine.Message = logLine.Message[:errBoundLimited.Limit] + " ...(truncated)"
			} else {
				if limits.IsErrRecoverable(err) {
					e.logger().Errorw("Failed to check user log line limit; emitting untruncated", "err", err)
				} else {
					e.logger().Errorw("User log line limit could not be evaluated; emitting untruncated", "err", err)
				}
				e.metrics.IncrementLimitCheckUnenforcedCounter(emitCtx, cresettings.Default.PerWorkflow.LogLineLimit.Key)
			}
		}

		if err := events.EmitUserLogs(emitCtx, executionLabels, []*protoevents.LogLine{logLine}, executionID); err != nil {
			e.logger().Errorw("Failed to emit user logs", "err", err)
		}
		count++
		return true
	}

	for {
		select {
		case <-ctx.Done():
			// Ensure we don't hang shutdown if the user logs channel is not drained.
			drainCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), emitUserLogsTimeout)
			for {
				select {
				case <-drainCtx.Done():
					e.logger().Warnw("Timeout reached while draining user logs")
					cancel()
					return
				case logLine, ok := <-userLogChan:
					if !ok || !processLogLine(drainCtx, logLine) {
						cancel()
						return
					}
				}
			}
		case logLine, ok := <-userLogChan:
			if !ok {
				return
			}
			// The execution context is cancelled the moment the execution completes,
			// but this goroutine outlives it and may still have buffered log lines.
			// Emit on a context detached from the execution lifecycle.
			emitCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), emitUserLogsTimeout)
			ok = processLogLine(emitCtx, logLine)
			cancel()
			if !ok {
				return
			}
		}
	}
}

func (e *baseEngine) donTimeRequestTimeout(ctx context.Context, limiter limits.TimeLimiter) time.Duration {
	if limiter == nil {
		return cresettings.Default.PerWorkflow.DONTime.RequestTimeout.DefaultValue
	}
	// A zero timeout is a valid, explicitly configured limit and is honoured as-is.
	limit, err := limiter.Limit(ctx)
	if err != nil {
		if !limits.IsErrRecoverable(err) {
			e.logger().Errorw("Failed to get DON time request timeout with no usable value; using the compiled default", "err", err)
			return cresettings.Default.PerWorkflow.DONTime.RequestTimeout.DefaultValue
		}
		e.logger().Errorw("Failed to get DON time request timeout; continuing with the value the limiter returned", "err", err)
		e.metrics.IncrementLimitReadFallbackCounter(ctx, cresettings.Default.PerWorkflow.DONTime.RequestTimeout.Key)
	}
	return limit
}
