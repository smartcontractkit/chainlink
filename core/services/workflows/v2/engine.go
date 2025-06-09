package v2

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"google.golang.org/protobuf/types/known/anypb"

	"github.com/smartcontractkit/chainlink-common/pkg/aggregation"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2/pb"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/monitoring"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/safe"
)

type Engine struct {
	services.Service
	srvcEng *services.Engine

	cfg          *EngineConfig
	lggr         logger.Logger
	loggerLabels map[string]string
	localNode    capabilities.Node

	// registration ID -> trigger capability
	triggers map[string]*triggerCapability
	// used to separate registration and unregistration phases
	triggersRegMu sync.Mutex

	allTriggerEventsQueueCh chan enqueuedTriggerEvent
	executionsSemaphore     chan struct{}
	capCallsSemaphore       chan struct{}

	meterReports *metering.Reports

	metrics *monitoring.WorkflowsMetricLabeler
}

type triggerCapability struct {
	capabilities.TriggerCapability
	payload *anypb.Any
}

type enqueuedTriggerEvent struct {
	triggerCapID string
	triggerIndex int
	timestamp    time.Time
	event        capabilities.TriggerResponse
}

func NewEngine(ctx context.Context, cfg *EngineConfig) (*Engine, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	em, err := monitoring.InitMonitoringResources()
	if err != nil {
		return nil, fmt.Errorf("could not initialize monitoring resources: %w", err)
	}
	localNode, err := cfg.CapRegistry.LocalNode(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get local node state: %w", err)
	}

	labels := []any{
		platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String(),
		platform.KeyWorkflowVersion, platform.ValueWorkflowVersionV2,
		platform.KeyDonID, strconv.Itoa(int(localNode.WorkflowDON.ID)),
		platform.KeyDonF, strconv.Itoa(int(localNode.WorkflowDON.F)),
		platform.KeyDonN, strconv.Itoa(len(localNode.WorkflowDON.Members)),
		platform.KeyDonQ, strconv.Itoa(aggregation.ByzantineQuorum(
			len(localNode.WorkflowDON.Members),
			int(localNode.WorkflowDON.F),
		)),
		platform.KeyP2PID, localNode.PeerID.String(),
	}

	beholderLogger := custmsg.NewBeholderLogger(cfg.Lggr, cfg.BeholderEmitter).Named("WorkflowEngine").With(labels...)
	metricsLabeler := monitoring.NewWorkflowsMetricLabeler(metrics.NewLabeler(), em).With(
		platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String())
	labelsMap := make(map[string]string, len(labels)/2)
	for i := 0; i < len(labels); i += 2 {
		labelsMap[labels[i].(string)] = labels[i+1].(string)
	}

	engine := &Engine{
		cfg:                     cfg,
		lggr:                    beholderLogger,
		loggerLabels:            labelsMap,
		localNode:               localNode,
		triggers:                make(map[string]*triggerCapability),
		allTriggerEventsQueueCh: make(chan enqueuedTriggerEvent, cfg.LocalLimits.TriggerEventQueueSize),
		executionsSemaphore:     make(chan struct{}, cfg.LocalLimits.MaxConcurrentWorkflowExecutions),
		capCallsSemaphore:       make(chan struct{}, cfg.LocalLimits.MaxConcurrentCapabilityCallsPerWorkflow),
		meterReports:            metering.NewReports(),
		metrics:                 metricsLabeler,
	}
	engine.Service, engine.srvcEng = services.Config{
		Name:  "WorkflowEngineV2",
		Start: engine.start,
		Close: engine.close,
	}.NewServiceEngine(beholderLogger)
	return engine, nil
}

func (e *Engine) start(_ context.Context) error {
	e.cfg.Module.Start()
	e.srvcEng.Go(e.heartbeatLoop)
	e.srvcEng.Go(e.init)
	e.srvcEng.Go(e.handleAllTriggerEvents)
	return nil
}

func (e *Engine) init(ctx context.Context) {
	// apply global engine instance limits
	// TODO(CAPPL-794): consider moving this outside of the engine, into the Syncer
	ownerAllow, globalAllow := e.cfg.GlobalLimits.Allow(e.cfg.WorkflowOwner)
	if !globalAllow {
		e.lggr.Info("Global workflow count limit reached")
		e.metrics.IncrementWorkflowLimitGlobalCounter(ctx)
		e.cfg.Hooks.OnInitialized(types.ErrGlobalWorkflowCountLimitReached)
		return
	}
	if !ownerAllow {
		e.lggr.Info("Per owner workflow count limit reached")
		e.metrics.IncrementWorkflowLimitPerOwnerCounter(ctx)
		e.cfg.Hooks.OnInitialized(types.ErrPerOwnerWorkflowCountLimitReached)
		return
	}

	err := e.runTriggerSubscriptionPhase(ctx)
	if err != nil {
		e.lggr.Errorw("Workflow Engine initialization failed", "err", err)
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	e.lggr.Info("Workflow Engine initialized")
	e.metrics.IncrementWorkflowInitializationCounter(ctx)
	e.cfg.Hooks.OnInitialized(nil)
}

func (e *Engine) runTriggerSubscriptionPhase(ctx context.Context) error {
	// call into the workflow to get trigger subscriptions
	subCtx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.TriggerSubscriptionRequestTimeoutMs))
	defer cancel()
	result, err := e.cfg.Module.Execute(subCtx, &wasmpb.ExecuteRequest{
		Request:         &wasmpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(e.cfg.LocalLimits.ModuleExecuteMaxResponseSizeBytes),
		// no Config needed
	}, DisallowedCapabilityExecutor{})
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
	if len(subs.Subscriptions) > int(e.cfg.LocalLimits.MaxTriggerSubscriptions) {
		return fmt.Errorf("too many trigger subscriptions: %d", len(subs.Subscriptions))
	}

	// check if all requested triggers exist in the registry
	triggers := make([]capabilities.TriggerCapability, 0, len(subs.Subscriptions))
	for _, sub := range subs.Subscriptions {
		triggerCap, err := e.cfg.CapRegistry.GetTrigger(ctx, sub.Id)
		if err != nil {
			return fmt.Errorf("trigger capability not found: %w", err)
		}
		triggers = append(triggers, triggerCap)
	}

	// register to all triggers
	e.triggersRegMu.Lock()
	defer e.triggersRegMu.Unlock()
	eventChans := make([]<-chan capabilities.TriggerResponse, len(subs.Subscriptions))
	triggerCapIDs := make([]string, len(subs.Subscriptions))
	for i, sub := range subs.Subscriptions {
		triggerCap := triggers[i]
		registrationID := fmt.Sprintf("trigger_reg_%s_%d", e.cfg.WorkflowID, i)
		// TODO(CAPPL-737): run with a timeout
		e.lggr.Debugw("Registering trigger", "triggerID", sub.Id, "method", sub.Method)
		triggerEventCh, err := triggerCap.RegisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:               e.cfg.WorkflowID,
				WorkflowOwner:            e.cfg.WorkflowOwner,
				WorkflowName:             e.cfg.WorkflowName.Hex(),
				DecodedWorkflowName:      e.cfg.WorkflowName.String(),
				WorkflowDonID:            e.localNode.WorkflowDON.ID,
				WorkflowDonConfigVersion: e.localNode.WorkflowDON.ConfigVersion,
				ReferenceID:              fmt.Sprintf("trigger_%d", i),
				// no WorkflowExecutionID needed (or available at this stage)
			},
			Payload: sub.Payload,
			Method:  sub.Method,
			// no Config needed - NoDAG uses Payload
		})
		if err != nil {
			e.lggr.Errorw("One of trigger registrations failed - reverting all", "triggerID", sub.Id, "err", err)
			e.metrics.With(platform.KeyTriggerID, sub.Id).IncrementRegisterTriggerFailureCounter(ctx)
			e.unregisterAllTriggers(ctx)
			return fmt.Errorf("failed to register trigger: %w", err)
		}
		e.triggers[registrationID] = &triggerCapability{
			TriggerCapability: triggerCap,
			payload:           sub.Payload,
		}
		eventChans[i] = triggerEventCh
		triggerCapIDs[i] = sub.Id
	}

	// start listening for trigger events only if all registrations succeeded
	for idx, triggerEventCh := range eventChans {
		e.srvcEng.Go(func(srvcCtx context.Context) {
			for {
				select {
				case <-srvcCtx.Done():
					return
				case event, isOpen := <-triggerEventCh:
					if !isOpen {
						return
					}
					select {
					case e.allTriggerEventsQueueCh <- enqueuedTriggerEvent{
						triggerCapID: subs.Subscriptions[idx].Id,
						triggerIndex: idx,
						timestamp:    e.cfg.Clock.Now(),
						event:        event,
					}:
					default: // queue full, drop the event
						e.lggr.Errorw("Trigger event queue is full, dropping event", "triggerID", subs.Subscriptions[idx].Id, "triggerIndex", idx)
					}
				}
			}
		})
	}
	e.lggr.Infow("All triggers registered successfully", "numTriggers", len(subs.Subscriptions), "triggerIDs", triggerCapIDs)
	e.metrics.IncrementWorkflowRegisteredCounter(ctx)
	e.cfg.Hooks.OnSubscribedToTriggers(triggerCapIDs)
	return nil
}

func (e *Engine) handleAllTriggerEvents(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case queueHead, isOpen := <-e.allTriggerEventsQueueCh:
			if !isOpen {
				return
			}
			// TODO(CAPPL-737): check if expired
			select {
			case e.executionsSemaphore <- struct{}{}: // block if too many concurrent workflow executions
				e.srvcEng.Go(func(srvcCtx context.Context) {
					e.startExecution(srvcCtx, queueHead)
					<-e.executionsSemaphore
				})
			case <-ctx.Done():
				return
			}
		}
	}
}

// startExecution initiates a new workflow execution, blocking until completed
func (e *Engine) startExecution(ctx context.Context, wrappedTriggerEvent enqueuedTriggerEvent) {
	triggerEvent := wrappedTriggerEvent.event.Event
	executionID, err := types.GenerateExecutionID(e.cfg.WorkflowID, triggerEvent.ID)
	if err != nil {
		e.lggr.Errorw("Failed to generate execution ID", "err", err, "triggerID", wrappedTriggerEvent.triggerCapID)
		return
	}

	// TODO(CAPPL-911): add rate-limiting

	subCtx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.WorkflowExecutionTimeoutMs))
	defer cancel()
	executionLogger := logger.With(e.lggr, "executionID", executionID, "triggerID", wrappedTriggerEvent.triggerCapID, "triggerIndex", wrappedTriggerEvent.triggerIndex)

	tid, err := safe.IntToUint64(wrappedTriggerEvent.triggerIndex)
	if err != nil {
		executionLogger.Errorw("Failed to convert trigger index to uint64", "err", err)
		return
	}

	e.meterReports.Add(executionID, metering.NewReport(e.cfg.Lggr))

	executionLogger.Infow("Workflow execution starting ...")
	_ = events.EmitExecutionStartedEvent(ctx, e.loggerLabels, triggerEvent.ID, executionID)

	result, err := e.cfg.Module.Execute(subCtx, &wasmpb.ExecuteRequest{
		Request: &wasmpb.ExecuteRequest_Trigger{
			Trigger: &sdkpb.Trigger{
				Id:      tid,
				Payload: triggerEvent.Payload,
			},
		},
		MaxResponseSize: uint64(e.cfg.LocalLimits.ModuleExecuteMaxResponseSizeBytes),
		// TODO(CAPPL-729): pass workflow config
	}, &CapabilityExecutor{Engine: e, WorkflowExecutionID: executionID})
	if err != nil {
		status := store.StatusErrored
		if errors.Is(err, context.DeadlineExceeded) {
			status = store.StatusTimeout
		}
		executionLogger.Errorw("Workflow execution failed", "err", err, "status", status)
		_ = events.EmitExecutionFinishedEvent(ctx, e.loggerLabels, status, executionID)
		e.meterReports.Delete(executionID)
		return
	}
	// TODO(CAPPL-737): measure and report execution time

	executionLogger.Infow("Workflow execution finished successfully")
	_ = events.EmitExecutionFinishedEvent(ctx, e.loggerLabels, store.StatusCompleted, executionID)
	e.meterReports.Delete(executionID)

	e.cfg.Hooks.OnResultReceived(result)
	e.cfg.Hooks.OnExecutionFinished(executionID)
}

func (e *Engine) close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(e.cfg.LocalLimits.ShutdownTimeoutMs))
	defer cancel()
	e.triggersRegMu.Lock()
	e.unregisterAllTriggers(ctx)
	e.triggersRegMu.Unlock()

	e.cfg.Module.Close()
	e.cfg.GlobalLimits.Decrement(e.cfg.WorkflowOwner)
	return nil
}

// NOTE: needs to be called under the triggersRegMu lock
func (e *Engine) unregisterAllTriggers(ctx context.Context) {
	for registrationID, trigger := range e.triggers {
		err := trigger.UnregisterTrigger(ctx, capabilities.TriggerRegistrationRequest{
			TriggerID: registrationID,
			Metadata: capabilities.RequestMetadata{
				WorkflowID:    e.cfg.WorkflowID,
				WorkflowDonID: e.localNode.WorkflowDON.ID,
			},
			Payload: trigger.payload,
		})
		if err != nil {
			e.cfg.Lggr.Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
		}
	}
	e.triggers = make(map[string]*triggerCapability)
	e.lggr.Infow("All triggers unregistered", "numTriggers", len(e.triggers))
	e.metrics.IncrementWorkflowUnregisteredCounter(ctx)
}

func (e *Engine) heartbeatLoop(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(e.cfg.LocalLimits.HeartbeatFrequencyMs) * time.Millisecond)
	defer ticker.Stop()
	e.lggr.Info("Starting heartbeat loop")
	e.metrics.EngineHeartbeatGauge(ctx, 1)

	for {
		select {
		case <-ctx.Done():
			e.metrics.EngineHeartbeatGauge(ctx, 0)
			e.lggr.Info("Shutting down heartbeat")
			return
		case <-ticker.C:
			e.lggr.Debugw("Engine heartbeat tick", "time", e.cfg.Clock.Now().Format(time.RFC3339))
			e.metrics.IncrementEngineHeartbeatCounter(ctx)
		}
	}
}
