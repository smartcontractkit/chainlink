package v2

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/aggregation"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/contexts"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/settings"
	"github.com/smartcontractkit/chainlink-common/pkg/settings/limits"
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

type Engine struct {
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

	// registration ID -> trigger capability
	triggers map[string]*triggerCapability
	// used to separate registration and unregistration phases
	triggersRegMu sync.Mutex

	allTriggerEventsQueueCh limits.QueueLimiter[enqueuedTriggerEvent]
	executionsSemaphore     limits.ResourcePoolLimiter[int]
	capCallsSemaphore       limits.ResourcePoolLimiter[int]

	meterReports *metering.Reports

	metrics *monitoring.WorkflowsMetricLabeler
}

type triggerCapability struct {
	capabilities.TriggerCapability
	payload *anypb.Any
	method  string
}

type enqueuedTriggerEvent struct {
	triggerCapID string
	triggerIndex int
	timestamp    time.Time
	event        capabilities.TriggerResponse
}

// buildLabels creates the label slice for the beholder logger based on config and localNode state.
// This is used both during engine creation and when updating labels after a DON configuration change.
func (e *Engine) buildLabels(localNode *capabilities.Node) []any {
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
		platform.DonVersion, strconv.FormatUint(uint64(localNode.WorkflowDON.ConfigVersion), 10),
	}
}

// logger returns the current logger in a thread-safe manner.
// This method should be used instead of accessing e.lggr directly to avoid race conditions
// when the logger is dynamically updated (e.g., when DON configuration changes).
func (e *Engine) logger() logger.SugaredLogger {
	e.lggrMu.RLock()
	defer e.lggrMu.RUnlock()
	return e.lggr
}

// setLogger updates the logger in a thread-safe manner.
// This is called when the DON configuration changes and we need to update the platform.DonVersion label.
func (e *Engine) setLogger(lggr logger.SugaredLogger) {
	e.lggrMu.Lock()
	defer e.lggrMu.Unlock()
	e.lggr = lggr
}

func NewEngine(cfg *EngineConfig) (*Engine, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	em, err := monitoring.InitMonitoringResources()
	if err != nil {
		return nil, fmt.Errorf("could not initialize monitoring resources: %w", err)
	}

	// LocalNode() is expected to be non-blocking at this stage (i.e. the registry is already synced)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.LocalLimits.LocalNodeTimeoutMs)*time.Millisecond)
	defer cancel()
	localNode, err := cfg.CapRegistry.LocalNode(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get local node state: %w", err)
	}

	// Create engine first so we can use the buildLabels method
	engine := &Engine{
		cfg:                     cfg,
		triggers:                make(map[string]*triggerCapability),
		allTriggerEventsQueueCh: cfg.LocalLimiters.TriggerEventQueue,
		executionsSemaphore:     cfg.LocalLimiters.ExecutionConcurrency,
		capCallsSemaphore:       cfg.LocalLimiters.CapabilityConcurrency,
	}

	// Build labels using the helper method
	labels := engine.buildLabels(&localNode)

	beholderLogger := logger.Sugared(custmsg.NewBeholderLogger(cfg.Lggr, cfg.BeholderEmitter).Named("WorkflowEngine").With(labels...))
	metricsLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em).With(
		platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String())
	labelsMap := make(map[string]string, len(labels)/2)
	for i := 0; i < len(labels); i += 2 {
		labelsMap[labels[i].(string)] = labels[i+1].(string)
	}

	if cfg.DebugMode {
		beholderLogger.Errorw("WARNING: Debug mode is enabled, this is not suitable for production")
	}

	// Store logger and other fields
	engine.setLogger(beholderLogger)
	engine.meterReports = metering.NewReports(cfg.BillingClient, cfg.WorkflowOwner, cfg.WorkflowID, beholderLogger, labelsMap, metricsLabeler, cfg.WorkflowRegistryAddress, cfg.WorkflowRegistryChainSelector, metering.EngineVersionV2)
	engine.metrics = metricsLabeler
	engine.loggerLabels.Store(&labelsMap)
	engine.localNode.Store(&localNode)
	engine.Service, engine.srvcEng = services.Config{
		Name:  "WorkflowEngineV2",
		Start: engine.start,
		Close: engine.close,
	}.NewServiceEngine(beholderLogger)
	return engine, nil
}

func (e *Engine) start(ctx context.Context) error {
	e.cfg.Module.Start()
	ctx = context.WithoutCancel(ctx)

	// Fetch organization ID for this owner
	organizationID := ""
	if e.cfg.OrgResolver != nil {
		orgID, gerr := e.cfg.OrgResolver.Get(ctx, e.cfg.WorkflowOwner)
		if gerr != nil {
			e.logger().Warnw("Failed to resolve organization ID, continuing without it", "workflowOwner", e.cfg.WorkflowOwner, "err", gerr)
		} else {
			organizationID = orgID
		}
	}
	loggerLabels := maps.Clone(*e.loggerLabels.Load())
	loggerLabels[platform.KeyOrganizationID] = organizationID
	e.loggerLabels.Store(&loggerLabels)

	ctx = contexts.WithCRE(ctx, contexts.CRE{Org: organizationID, Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID})
	e.srvcEng.GoCtx(ctx, e.heartbeatLoop)
	e.srvcEng.GoCtx(ctx, e.init)
	e.srvcEng.GoCtx(ctx, e.handleAllTriggerEvents)
	return nil
}

func (e *Engine) init(ctx context.Context) {
	// apply global engine instance limits
	// TODO(CAPPL-794): consider moving this outside of the engine, into the Syncer
	err := e.cfg.GlobalExecutionConcurrencyLimiter.Use(ctx, 1)
	if err != nil {
		var errLimited limits.ErrorResourceLimited[int]
		if errors.As(err, &errLimited) {
			switch errLimited.Scope {
			case settings.ScopeOwner:
				e.logger().Info("Per owner workflow count limit reached", "err", err)
				e.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
				e.cfg.Hooks.OnInitialized(types.ErrPerOwnerWorkflowCountLimitReached)
			case settings.ScopeGlobal:
				e.logger().Info("Global workflow count limit reached", "err", err)
				e.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
				e.cfg.Hooks.OnInitialized(types.ErrGlobalWorkflowCountLimitReached)
			default:
				e.logger().Errorw("Workflow count limit reached for unexpected scope", "scope", errLimited.Scope, "err", err)
				e.cfg.Hooks.OnInitialized(err)
			}
		} else {
			e.cfg.Hooks.OnInitialized(err)
		}
		return
	}

	donSubCh, cleanup, err := e.cfg.DonSubscriber.Subscribe(ctx)
	if err != nil {
		e.logger().Errorw("failed to subscribe to DON notifier", "error", err)
		e.cfg.Hooks.OnInitialized(fmt.Errorf("failed to subscribe to DON notifier: %w", err))
		return
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

	err = e.runTriggerSubscriptionPhase(ctx)
	if err != nil {
		e.logger().Errorw("Workflow Engine initialization failed", "err", err)
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	e.logger().Info("Workflow Engine initialized")
	e.metrics.IncrementWorkflowInitializationCounter(ctx)
	e.cfg.Hooks.OnInitialized(nil)
}

func (e *Engine) localNodeSync(ctx context.Context) {
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
		"Workflow DON Config Version", localNode.WorkflowDON.ConfigVersion,
	)

	// Recreate the beholder logger with updated labels to reflect the new DON version
	labels := e.buildLabels(&localNode)
	newLogger := logger.Sugared(
		custmsg.NewBeholderLogger(e.cfg.Lggr, e.cfg.BeholderEmitter).
			Named("WorkflowEngine").
			With(labels...),
	)
	e.setLogger(newLogger)

	// Update loggerLabels map for metrics
	labelsMap := make(map[string]string, len(labels)/2)
	for i := 0; i < len(labels); i += 2 {
		labelsMap[labels[i].(string)] = labels[i+1].(string)
	}
	e.loggerLabels.Store(&labelsMap)

	e.cfg.Hooks.OnNodeSynced(localNode, nil)
	e.localNode.Store(&localNode)
}

func (e *Engine) runTriggerSubscriptionPhase(ctx context.Context) error {
	// call into the workflow to get trigger subscriptions
	subCtx, subCancel, err := e.cfg.LocalLimiters.TriggerSubscriptionTime.WithTimeout(ctx)
	if err != nil {
		return err
	}
	defer subCancel()

	maxUserLogEventsPerExecution, err := e.cfg.LocalLimiters.LogEvent.Limit(ctx)
	if err != nil {
		return err
	}
	userLogChan := make(chan *protoevents.LogLine, maxUserLogEventsPerExecution)
	defer close(userLogChan)
	e.srvcEng.Go(func(_ context.Context) {
		e.emitUserLogs(subCtx, userLogChan, e.cfg.WorkflowID, *e.loggerLabels.Load())
	})

	var timeProvider TimeProvider = &types.LocalTimeProvider{}
	if !e.cfg.UseLocalTimeProvider {
		timeProvider = NewDonTimeProvider(e.cfg.DonTimeStore, e.cfg.WorkflowID, e.logger())
	}

	moduleExecuteMaxResponseSizeBytes, err := e.cfg.LocalLimiters.ExecutionResponse.Limit(ctx)
	if err != nil {
		return err
	}
	if moduleExecuteMaxResponseSizeBytes < 0 {
		return fmt.Errorf("invalid moduleExecuteMaxResponseSizeBytes; must not be negative: %d", moduleExecuteMaxResponseSizeBytes)
	}
	result, err := e.cfg.Module.Execute(subCtx, &sdkpb.ExecuteRequest{
		Request:         &sdkpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(moduleExecuteMaxResponseSizeBytes), //nolint:gosec // G115
		Config:          e.cfg.WorkflowConfig,
	}, NewDisallowedExecutionHelper(e.logger(), userLogChan, timeProvider, e.secretsFetcher(e.cfg.WorkflowID)))
	if err != nil {
		return fmt.Errorf("failed to execute subscribe: %w", err)
	}
	if result.GetError() != "" {
		return fmt.Errorf("failed to execute subscribe: %s", result.GetError())
	}
	subs := result.GetTriggerSubscriptions()
	if subs == nil {
		return errors.New("subscribe result is nil")
	}
	err = e.cfg.LocalLimiters.TriggerSubscription.Check(ctx, len(subs.Subscriptions))
	if err != nil {
		return err
	}

	// check if all requested triggers exist in the registry
	triggers := make([]capabilities.TriggerCapability, 0, len(subs.Subscriptions))
	for _, sub := range subs.Subscriptions {
		_, labels, _ := capabilities.ParseID(sub.Id)
		chainSelector, err2 := capabilities.ChainSelectorLabel(labels)
		if err2 != nil {
			return fmt.Errorf("invalid chain selector for ID %s: %w", sub.Id, err2)
		}
		if chainSelector != nil {
			err2 := e.cfg.LocalLimiters.ChainAllowed.AllowErr(contexts.WithChainSelector(ctx, *chainSelector))
			if err2 != nil {
				if errors.Is(err2, limits.ErrorNotAllowed{}) {
					return fmt.Errorf("unable to subscribe to capability %s: ChainSelector %d: %w", sub.Id, *chainSelector, err2)
				}
				return fmt.Errorf("failed to check access for ChainSelector %d: %w", *chainSelector, err2)
			}
		}
		triggerCap, triggerErr := e.cfg.CapRegistry.GetTrigger(ctx, sub.Id)
		if triggerErr != nil {
			return fmt.Errorf("trigger capability not found: %w", triggerErr)
		}
		triggers = append(triggers, triggerCap)
	}

	// register to all triggers concurrently
	regCtx, regCancel, err := e.cfg.LocalLimiters.TriggerRegistrationsTime.WithTimeout(ctx)
	if err != nil {
		return err
	}
	defer regCancel()

	// trigger registration results for use in concurrent trigger subscriptions
	type triggerRegResult struct {
		index          int
		registrationID string
		triggerCap     capabilities.TriggerCapability
		eventCh        <-chan capabilities.TriggerResponse
		payload        *anypb.Any
		method         string
		triggerCapID   string
	}

	resultsCh := make(chan triggerRegResult, len(subs.Subscriptions))
	g, gCtx := errgroup.WithContext(regCtx)

	// Launch concurrent trigger registrations
	for i, sub := range subs.Subscriptions {
		triggerCap := triggers[i]
		g.Go(func() error {
			registrationID := fmt.Sprintf("trigger_reg_%s_%d", e.cfg.WorkflowID, i)
			e.logger().Debugw("Registering trigger", "triggerID", sub.Id, "method", sub.Method)
			triggerEventCh, regErr := triggerCap.RegisterTrigger(gCtx, capabilities.TriggerRegistrationRequest{
				TriggerID: registrationID,
				Metadata: capabilities.RequestMetadata{
					WorkflowID:                    e.cfg.WorkflowID,
					WorkflowOwner:                 e.cfg.WorkflowOwner,
					WorkflowName:                  e.cfg.WorkflowName.Hex(),
					WorkflowTag:                   e.cfg.WorkflowTag,
					DecodedWorkflowName:           e.cfg.WorkflowName.String(),
					WorkflowDonID:                 e.localNode.Load().WorkflowDON.ID,
					WorkflowDonConfigVersion:      e.localNode.Load().WorkflowDON.ConfigVersion,
					ReferenceID:                   fmt.Sprintf("trigger_%d", i),
					WorkflowRegistryChainSelector: e.cfg.WorkflowRegistryChainSelector,
					WorkflowRegistryAddress:       e.cfg.WorkflowRegistryAddress,
					EngineVersion:                 platform.ValueWorkflowVersionV2,
					// no WorkflowExecutionID needed (or available at this stage)
				},
				Payload: sub.Payload,
				Method:  sub.Method,
				// no Config needed - NoDAG uses Payload
			})
			if regErr != nil {
				e.logger().Errorw("Trigger registration failed", "triggerID", sub.Id, "err", regErr)
				e.metrics.With(platform.KeyTriggerID, sub.Id).IncrementRegisterTriggerFailureCounter(gCtx)
				return fmt.Errorf("failed to register trigger %s: %w", sub.Id, regErr)
			}
			// Send successful result
			resultsCh <- triggerRegResult{
				index:          i,
				registrationID: registrationID,
				triggerCap:     triggerCap,
				eventCh:        triggerEventCh,
				payload:        sub.Payload,
				method:         sub.Method,
				triggerCapID:   sub.Id,
			}
			return nil
		})
	}

	// wait for all registrations to complete.
	// returns first non-nil error.
	registrationErr := g.Wait()
	close(resultsCh)

	// Collect results into e.triggers map
	e.triggersRegMu.Lock()
	defer e.triggersRegMu.Unlock()

	eventChans := make([]<-chan capabilities.TriggerResponse, len(subs.Subscriptions))
	triggerCapIDs := make([]string, len(subs.Subscriptions))

	for result := range resultsCh {
		e.triggers[result.registrationID] = &triggerCapability{
			TriggerCapability: result.triggerCap,
			payload:           result.payload,
			method:            result.method,
		}
		eventChans[result.index] = result.eventCh
		triggerCapIDs[result.index] = result.triggerCapID
	}

	// If any registration failed, unregister successful ones and return error
	if registrationErr != nil {
		e.logger().Errorw("One or more trigger registrations failed - reverting all", "err", registrationErr)
		e.unregisterAllTriggers(ctx) // needs to be called under e.triggersRegMu lock
		return registrationErr
	}

	// start listening for trigger events only if all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		e.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					return
				case event, isOpen := <-triggerEventCh:
					if !isOpen {
						return
					}
					triggerID := subs.Subscriptions[idx].Id
					eventID := event.Event.ID
					e.logger().Debugw("Processing trigger event", "triggerID", triggerID, "eventID", eventID)
					if event.Err != nil {
						e.logger().Errorw("Received a trigger event with error, dropping", "triggerID", triggerID, "err", event.Err)
						e.metrics.With(platform.KeyTriggerID, triggerID).IncrementWorkflowTriggerEventErrorCounter(ctx)
						continue
					}
					if err := e.allTriggerEventsQueueCh.Put(ctx, enqueuedTriggerEvent{
						triggerCapID: triggerID,
						triggerIndex: idx,
						timestamp:    e.cfg.Clock.Now(),
						event:        event,
					}); err != nil {
						var errFull limits.ErrorQueueFull
						if errors.As(err, &errFull) {
							// queue full, drop the event
							e.logger().Errorw("Trigger event queue is full, dropping event", "triggerID", triggerID, "triggerIndex", idx, "err", err)
							e.metrics.With(platform.KeyTriggerID, triggerID).IncrementWorkflowTriggerEventQueueFullCounter(ctx)
						}
						e.logger().Errorw("Failed to enqueue trigger event", "triggerID", triggerID, "triggerIndex", idx, "err", err)
						e.metrics.With(platform.KeyTriggerID, triggerID).IncrementWorkflowTriggerEventErrorCounter(ctx)
						continue
					}
					e.logger().Debugw("Enqueued trigger event", "triggerID", triggerID, "eventID", eventID)
				}
			}
		})
	}
	e.logger().Infow("All triggers registered successfully", "numTriggers", len(subs.Subscriptions), "triggerIDs", triggerCapIDs)
	e.metrics.IncrementWorkflowRegisteredCounter(ctx)
	e.cfg.Hooks.OnSubscribedToTriggers(triggerCapIDs)
	return nil
}

func (e *Engine) handleAllTriggerEvents(ctx context.Context) {
	for {
		queueHead, err := e.allTriggerEventsQueueCh.Wait(ctx)
		if err != nil {
			return
		}
		eventAge := queueHead.timestamp.Sub(e.cfg.Clock.Now())
		eventID := queueHead.event.Event.ID
		e.logger().Debugw("Popped a trigger event from the queue", "eventID", eventID, "eventAgeMs", eventAge.Milliseconds())
		triggerEventMaxAge, err := e.cfg.LocalLimiters.TriggerEventQueueTime.Limit(ctx)
		if err != nil {
			e.logger().Errorw("Failed to get trigger event queue time limit", "err", err)
			continue
		}
		if eventAge > triggerEventMaxAge {
			e.logger().Warnw("Trigger event is too old, skipping execution", "triggerID", queueHead.triggerCapID, "eventID", eventID, "eventAgeMs", eventAge.Milliseconds())
			continue
		}
		free, err := e.executionsSemaphore.Wait(ctx, 1) // block if too many concurrent workflow executions
		if err != nil {
			e.logger().Errorw("Failed to acquire executions semaphore", "err", err)
			continue
		}
		e.logger().Debugw("Scheduling a trigger event for execution", "eventID", eventID)
		e.srvcEng.GoCtx(context.WithoutCancel(ctx), func(ctx context.Context) {
			defer free()
			e.startExecution(ctx, queueHead)
		})
	}
}

// startExecution initiates a new workflow execution, blocking until completed
func (e *Engine) startExecution(ctx context.Context, wrappedTriggerEvent enqueuedTriggerEvent) {
	triggerEvent := wrappedTriggerEvent.event.Event
	executionID, err := events.GenerateExecutionID(e.cfg.WorkflowID, triggerEvent.ID)
	if err != nil {
		e.logger().Errorw("Failed to generate execution ID", "err", err, "triggerID", wrappedTriggerEvent.triggerCapID)
		return
	}

	// Fetch organization ID for this execution
	organizationID := contexts.CREValue(ctx).Org
	if e.cfg.OrgResolver != nil {
		orgID, gerr := e.cfg.OrgResolver.Get(ctx, e.cfg.WorkflowOwner)
		if gerr != nil {
			e.logger().Warnw("Failed to resolve organization ID, continuing without it", "workflowOwner", e.cfg.WorkflowOwner, "err", gerr)
		} else {
			organizationID = orgID

			creCtx := contexts.CREValue(ctx)
			creCtx.Org = organizationID
			ctx = contexts.WithCRE(ctx, creCtx)
		}
	}
	loggerLabels := maps.Clone(*e.loggerLabels.Load())
	loggerLabels[platform.KeyOrganizationID] = organizationID
	e.loggerLabels.Store(&loggerLabels)
	lggr := e.logger().With(platform.KeyOrganizationID, organizationID)

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
			return
		}

		e.deductStandardBalances(ctx, meteringReport)
	}

	execCtx, execCancel, err := e.cfg.LocalLimiters.ExecutionTime.WithTimeout(ctx)
	if err != nil {
		lggr.Errorw("Failed to get execution time limit", "err", err)
		return
	}
	defer execCancel()
	executionLogger := logger.With(lggr, "executionID", executionID, "triggerID", wrappedTriggerEvent.triggerCapID, "triggerIndex", wrappedTriggerEvent.triggerIndex)

	maxUserLogEventsPerExecution, err := e.cfg.LocalLimiters.LogEvent.Limit(ctx)
	if err != nil {
		lggr.Errorw("Failed to get log event limit", "err", err)
		return
	}
	userLogChan := make(chan *protoevents.LogLine, maxUserLogEventsPerExecution)
	defer close(userLogChan)
	e.srvcEng.Go(func(_ context.Context) {
		e.emitUserLogs(execCtx, userLogChan, executionID, loggerLabels)
	})

	tid, err := safe.IntToUint64(wrappedTriggerEvent.triggerIndex)
	if err != nil {
		executionLogger.Errorw("Failed to convert trigger index to uint64", "err", err)
		return
	}

	startTime := e.cfg.Clock.Now()
	executionLogger.Infow("Workflow execution starting ...")
	_ = events.EmitExecutionStartedEvent(ctx, loggerLabels, triggerEvent.ID, executionID)
	e.metrics.With("workflowID", e.cfg.WorkflowID, "workflowName", e.cfg.WorkflowName.String()).IncrementWorkflowExecutionStartedCounter(ctx)
	var executionStatus string // store.StatusStarted

	var timeProvider TimeProvider = &types.LocalTimeProvider{}
	if !e.cfg.UseLocalTimeProvider {
		timeProvider = NewDonTimeProvider(e.cfg.DonTimeStore, e.cfg.WorkflowID, lggr)
	}

	moduleExecuteMaxResponseSizeBytes, err := e.cfg.LocalLimiters.ExecutionResponse.Limit(ctx)
	if err != nil {
		lggr.Errorw("Failed to get execution response size limit", "err", err)
		return
	}
	if moduleExecuteMaxResponseSizeBytes < 0 {
		lggr.Errorf("invalid moduleExecuteMaxResponseSizeBytes; must not be negative: %d", moduleExecuteMaxResponseSizeBytes)
		return
	}
	execHelper := &ExecutionHelper{
		Engine: e, WorkflowExecutionID: executionID, UserLogChan: userLogChan,
		TimeProvider: timeProvider, SecretsFetcher: e.secretsFetcher(executionID),
	}
	execHelper.initLimiters(e.cfg.LocalLimiters)
	result, execErr := e.cfg.Module.Execute(execCtx, &sdkpb.ExecuteRequest{
		Request: &sdkpb.ExecuteRequest_Trigger{
			Trigger: &sdkpb.Trigger{
				Id:      tid,
				Payload: triggerEvent.Payload,
			},
		},
		MaxResponseSize: uint64(moduleExecuteMaxResponseSizeBytes), //nolint:gosec // G115
		Config:          e.cfg.WorkflowConfig,
	}, execHelper)

	endTime := e.cfg.Clock.Now()
	executionDuration := endTime.Sub(startTime)

	if isMetering {
		computeUnit := billing.ResourceType_name[int32(billing.ResourceType_RESOURCE_TYPE_COMPUTE)]
		mrErr := meteringReport.Settle(computeUnit,
			capabilities.ResponseMetadata{
				Metering: []capabilities.MeteringNodeDetail{{
					Peer2PeerID: e.localNode.Load().PeerID.String(),
					SpendUnit:   computeUnit,
					SpendValue:  strconv.Itoa(int(executionDuration.Milliseconds())),
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
		if errors.Is(execErr, context.DeadlineExceeded) {
			executionStatus = store.StatusTimeout
			e.metrics.UpdateWorkflowTimeoutDurationHistogram(ctx, int64(executionDuration.Seconds()))
		} else {
			e.metrics.UpdateWorkflowErrorDurationHistogram(ctx, int64(executionDuration.Seconds()))
		}

		executionLogger.Errorw("Workflow execution failed with module execution error", "status", executionStatus, "durationMs", executionDuration.Milliseconds(), "err", execErr)
		_ = events.EmitExecutionFinishedEvent(ctx, loggerLabels, executionStatus, executionID, lggr)
		e.cfg.Hooks.OnExecutionFinished(executionID, executionStatus)
		e.cfg.Hooks.OnExecutionError(execErr.Error())
		return
	}

	if e.cfg.DebugMode {
		lggr.Debugw("User workflow execution result", "result", result.GetValue(), "err", result.GetError())
	}

	if len(result.GetError()) > 0 {
		executionStatus = store.StatusErrored
		e.metrics.UpdateWorkflowErrorDurationHistogram(ctx, int64(executionDuration.Seconds()))
		e.metrics.With("workflowID", e.cfg.WorkflowID, "workflowName", e.cfg.WorkflowName.String()).IncrementWorkflowExecutionFailedCounter(ctx)
		executionLogger.Errorw("Workflow execution failed", "status", executionStatus, "durationMs", executionDuration.Milliseconds(), "error", result.GetError())
		_ = events.EmitExecutionFinishedEvent(ctx, loggerLabels, executionStatus, executionID, lggr)
		e.cfg.Hooks.OnExecutionFinished(executionID, executionStatus)
		e.cfg.Hooks.OnExecutionError(result.GetError())
		return
	}

	executionStatus = store.StatusCompleted
	executionLogger.Infow("Workflow execution finished successfully", "durationMs", executionDuration.Milliseconds())
	_ = events.EmitExecutionFinishedEvent(ctx, loggerLabels, executionStatus, executionID, lggr)
	e.metrics.UpdateWorkflowCompletedDurationHistogram(ctx, int64(executionDuration.Seconds()))
	e.metrics.With("workflowID", e.cfg.WorkflowID, "workflowName", e.cfg.WorkflowName.String()).IncrementWorkflowExecutionSucceededCounter(ctx)
	e.cfg.Hooks.OnResultReceived(result)
	e.cfg.Hooks.OnExecutionFinished(executionID, executionStatus)
}

func (e *Engine) secretsFetcher(phaseID string) SecretsFetcher {
	if e.cfg.SecretsFetcher != nil {
		return e.cfg.SecretsFetcher
	}

	return NewSecretsFetcher(
		e.metrics,
		e.cfg.CapRegistry,
		e.logger(),
		e.cfg.LocalLimiters.SecretsConcurrency,
		e.cfg.WorkflowOwner,
		e.cfg.WorkflowName.String(),
		e.cfg.WorkflowID,
		// phaseID is the executionID if called during an execution,
		// or the workflowID if called during trigger subscription
		phaseID,
		e.cfg.WorkflowEncryptionKey,
	)
}

func (e *Engine) close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(e.cfg.LocalLimits.ShutdownTimeoutMs))
	defer cancel()
	ctx = contexts.WithCRE(ctx, contexts.CRE{Owner: e.cfg.WorkflowOwner, Workflow: e.cfg.WorkflowID}) // TODO org?
	e.triggersRegMu.Lock()
	e.unregisterAllTriggers(ctx)
	e.triggersRegMu.Unlock()
	e.metrics.IncrementWorkflowUnregisteredCounter(ctx)

	e.cfg.Module.Close()

	// reset metering mode metric so that a positive value does not persist
	e.metrics.UpdateWorkflowMeteringModeGauge(ctx, false)

	return e.cfg.GlobalExecutionConcurrencyLimiter.Free(ctx, 1)
}

// NOTE: needs to be called under the triggersRegMu lock
func (e *Engine) unregisterAllTriggers(ctx context.Context) {
	failCount := 0
	for registrationID, trigger := range e.triggers {
		err := trigger.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    e.cfg.WorkflowID,
				WorkflowDonID: e.localNode.Load().WorkflowDON.ID,
			},
			Payload: trigger.payload,
			Method:  trigger.method,
		})
		if err != nil {
			e.logger().Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
			failCount++
		}
	}
	e.logger().Infow("All triggers unregistered", "numTriggers", len(e.triggers), "failed", failCount)
	e.triggers = make(map[string]*triggerCapability)
}

func (e *Engine) heartbeatLoop(ctx context.Context) {
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

func (e *Engine) deductStandardBalances(ctx context.Context, meteringReport *metering.Report) {
	// V2Engine runs the entirety of a module's execution as compute. Ensure that the max execution time can run.
	// Add an extra second of metering padding for context cancel propagation
	ctxCancelPadding := (time.Millisecond * 1000).Milliseconds()
	workflowExecutionTimeout, err := e.cfg.LocalLimiters.ExecutionTime.Limit(ctx)
	if err != nil {
		e.logger().Errorw("Failed to get execution time limit", "err", err)
		return
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

// separate call for each workflow execution
func (e *Engine) emitUserLogs(ctx context.Context, userLogChan chan *protoevents.LogLine, executionID string, executionLabels map[string]string) {
	e.logger().Debugw("Listening for user logs ...")
	count := 0
	defer func() { e.logger().Debugw("Listening for user logs done.", "processedLogLines", count) }()
	for {
		select {
		case <-ctx.Done():
			return
		case logLine, ok := <-userLogChan:
			if !ok {
				return
			}
			if e.cfg.DebugMode {
				e.logger().Debugf("User log: <<<%s>>>, local node timestamp: %s", logLine.Message, logLine.NodeTimestamp)
			}
			err := e.cfg.LocalLimiters.LogEvent.Check(ctx, count)
			if err != nil {
				var errBoundLimited limits.ErrorBoundLimited[int]
				if errors.As(err, &errBoundLimited) {
					e.logger().Warnw("Max user log events per execution reached, dropping event", "maxEvents", errBoundLimited.Limit)
					return
				}
				e.logger().Errorw("Failed to get user log event limit", "err", err)
				return
			}
			maxUserLogLength, err := e.cfg.LocalLimiters.LogLine.Limit(ctx)
			if err != nil {
				e.logger().Errorw("Failed to get user log line limit", "err", err)
				return
			}
			if len(logLine.Message) > int(maxUserLogLength) {
				logLine.Message = logLine.Message[:maxUserLogLength] + " ...(truncated)"
			}

			if err := events.EmitUserLogs(ctx, executionLabels, []*protoevents.LogLine{logLine}, executionID); err != nil {
				e.logger().Errorw("Failed to emit user logs", "err", err)
			}
			count++
		}
	}
}
