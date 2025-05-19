package v2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/services"

	sdkpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk/v2/pb"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/host"
	wasmpb "github.com/smartcontractkit/chainlink-common/pkg/workflows/wasm/v2/pb"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/safe"
)

type Engine struct {
	services.Service
	srvcEng *services.Engine

	cfg       *EngineConfig
	localNode capabilities.Node

	// registration ID -> trigger capability
	triggers map[string]capabilities.TriggerCapability
	// used to separate registration and unregistration phases
	triggersRegMu sync.Mutex

	allTriggerEventsQueueCh chan enqueuedTriggerEvent
	executionsSemaphore     chan struct{}
	capCallsSemaphore       chan struct{}
}

type enqueuedTriggerEvent struct {
	triggerCapID string
	triggerIndex int
	timestamp    time.Time
	event        capabilities.TriggerResponse
}

var _ host.CapabilityExecutor = (*Engine)(nil)

func NewEngine(ctx context.Context, cfg *EngineConfig) (*Engine, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	engine := &Engine{
		cfg:                     cfg,
		triggers:                make(map[string]capabilities.TriggerCapability),
		allTriggerEventsQueueCh: make(chan enqueuedTriggerEvent, cfg.LocalLimits.TriggerEventQueueSize),
		executionsSemaphore:     make(chan struct{}, cfg.LocalLimits.MaxConcurrentWorkflowExecutions),
		capCallsSemaphore:       make(chan struct{}, cfg.LocalLimits.MaxConcurrentCapabilityCallsPerWorkflow),
	}
	engine.Service, engine.srvcEng = services.Config{
		Name:  "WorkflowEngineV2",
		Start: engine.start,
		Close: engine.close,
	}.NewServiceEngine(cfg.Lggr.Named("WorkflowEngine").With("workflowID", cfg.WorkflowID))
	return engine, nil
}

func (e *Engine) start(_ context.Context) error {
	e.cfg.Module.Start()
	e.srvcEng.Go(e.init)
	e.srvcEng.Go(e.handleAllTriggerEvents)
	return nil
}

func (e *Engine) init(ctx context.Context) {
	// apply global engine instance limits
	// TODO(CAPPL-794): consider moving this outside of the engine, into the Syncer
	ownerAllow, globalAllow := e.cfg.GlobalLimits.Allow(e.cfg.WorkflowOwner)
	if !globalAllow {
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(types.ErrGlobalWorkflowCountLimitReached)
		return
	}
	if !ownerAllow {
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(types.ErrPerOwnerWorkflowCountLimitReached)
		return
	}

	// retrieve info about the current node we are running on
	retryErr := internal.RunWithRetries(
		ctx,
		e.cfg.Lggr,
		time.Millisecond*time.Duration(e.cfg.LocalLimits.CapRegistryAccessRetryIntervalMs),
		int(e.cfg.LocalLimits.MaxCapRegistryAccessRetries),
		func() error {
			// retry until the underlying peerWrapper service is ready
			node, err := e.cfg.CapRegistry.LocalNode(ctx)
			if err != nil {
				return fmt.Errorf("failed to get donInfo: %w", err)
			}
			e.localNode = node
			return nil
		})

	if retryErr != nil {
		e.cfg.Lggr.Errorw("Workflow Engine initialization failed", "err", retryErr)
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(retryErr)
		return
	}

	err := e.cfg.Module.SetCapabilityExecutor(e)
	if err != nil {
		e.cfg.Lggr.Errorw("Workflow Engine initialization failed", "err", err)
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	err = e.runTriggerSubscriptionPhase(ctx)
	if err != nil {
		e.cfg.Lggr.Errorw("Workflow Engine initialization failed", "err", err)
		// TODO(CAPPL-736): observability
		e.cfg.Hooks.OnInitialized(err)
		return
	}

	e.cfg.Lggr.Info("Workflow Engine initialized")
	e.cfg.Hooks.OnInitialized(nil)
}

func (e *Engine) runTriggerSubscriptionPhase(ctx context.Context) error {
	// call into the workflow to get trigger subscriptions
	subCtx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.TriggerSubscriptionRequestTimeoutMs))
	defer cancel()
	result, err := e.cfg.Module.Execute(subCtx, &wasmpb.ExecuteRequest{
		Id:              "subscribe_" + uuid.New().String(), // execution ID for the subscription phase (one-time, not very useful)
		Request:         &wasmpb.ExecuteRequest_Subscribe{},
		MaxResponseSize: uint64(e.cfg.LocalLimits.ModuleExecuteMaxResponseSizeBytes),
		// no Config needed
	})
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
		e.cfg.Lggr.Debugw("Registering trigger", "triggerID", sub.Id)
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
			e.cfg.Lggr.Errorw("One of trigger registrations failed - reverting all", "triggerID", sub.Id, "err", err)
			e.unregisterAllTriggers(ctx)
			return fmt.Errorf("failed to register trigger: %w", err)
		}
		e.triggers[registrationID] = triggerCap
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
						// TODO(CAPPL-736): observability
					}
				}
			}
		})
	}
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
		// TODO(CAPPL-736): observability
		return
	}

	subCtx, cancel := context.WithTimeout(ctx, time.Millisecond*time.Duration(e.cfg.LocalLimits.WorkflowExecutionTimeoutMs))
	defer cancel()

	tid, err := safe.IntToUint64(wrappedTriggerEvent.triggerIndex)
	if err != nil {
		// TODO(CAPPL-736): observability
		return
	}

	result, err := e.cfg.Module.Execute(subCtx, &wasmpb.ExecuteRequest{
		Id: executionID,
		Request: &wasmpb.ExecuteRequest_Trigger{
			Trigger: &sdkpb.Trigger{
				Id:      tid,
				Payload: triggerEvent.Payload,
			},
		},
		MaxResponseSize: uint64(e.cfg.LocalLimits.ModuleExecuteMaxResponseSizeBytes),
		// TODO(CAPPL-729): pass workflow config
	})
	if err != nil {
		e.cfg.Lggr.Errorw("Workflow execution failed", "err", err)
		// TODO(CAPPL-736): observability
		return
	}

	// TODO(CAPPL-736): handle execution result
	e.cfg.Lggr.Debugw("Workflow execution finished", "executionID", executionID, "result", result)
	e.cfg.Hooks.OnResultReceived(result)
	e.cfg.Hooks.OnExecutionFinished(executionID)
}

// CallCapability handles requests generated by the wasm guest
func (e *Engine) CallCapability(ctx context.Context, request *sdkpb.CapabilityRequest) (*sdkpb.CapabilityResponse, error) {
	select {
	case e.capCallsSemaphore <- struct{}{}: // block if too many concurrent capability calls
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	defer func() { <-e.capCallsSemaphore }()

	// TODO (CAPPL-735): use request.Metadata.WorkflowExecutionId to associate the call with a specific execution
	capability, err := e.cfg.CapRegistry.GetExecutable(ctx, request.Id)
	if err != nil {
		return nil, fmt.Errorf("trigger capability not found: %w", err)
	}

	capReq := capabilities.CapabilityRequest{
		Payload:      request.Payload,
		Method:       request.Method,
		CapabilityId: request.Id,
		Metadata: capabilities.RequestMetadata{
			WorkflowExecutionID: request.ExecutionId,
		},
	}

	// TODO(CAPPL-737): run with a timeout
	capResp, err := capability.Execute(ctx, capReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute capability: %w", err)
	}
	return &sdkpb.CapabilityResponse{
		Response: &sdkpb.CapabilityResponse_Payload{
			Payload: capResp.Payload,
		},
	}, nil
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
		})
		if err != nil {
			e.cfg.Lggr.Errorw("Failed to unregister trigger", "registrationId", registrationID, "err", err)
		}
	}
	e.triggers = make(map[string]capabilities.TriggerCapability)
}
