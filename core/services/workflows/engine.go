package workflows

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/exec"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
)

type stepRequest struct {
	stepRef string
	state   store.WorkflowExecution
}

type stepUpdateChannel struct {
	executionID string
	ch          chan store.WorkflowExecutionStep
}

type stepUpdateManager struct {
	mu sync.RWMutex
	m  map[string]stepUpdateChannel
}

func (sucm *stepUpdateManager) add(executionID string, ch stepUpdateChannel) (added bool) {
	sucm.mu.RLock()
	_, ok := sucm.m[executionID]
	sucm.mu.RUnlock()
	if ok {
		return false
	}
	sucm.mu.Lock()
	defer sucm.mu.Unlock()
	if _, ok = sucm.m[executionID]; ok {
		return false
	}
	sucm.m[executionID] = ch
	return true
}

func (sucm *stepUpdateManager) remove(executionID string) {
	sucm.mu.Lock()
	defer sucm.mu.Unlock()
	if _, ok := sucm.m[executionID]; ok {
		close(sucm.m[executionID].ch)
		delete(sucm.m, executionID)
	}
}

func (sucm *stepUpdateManager) send(ctx context.Context, executionID string, stepUpdate store.WorkflowExecutionStep) error {
	sucm.mu.RLock()
	stepUpdateCh, ok := sucm.m[executionID]
	sucm.mu.RUnlock()
	if !ok {
		return fmt.Errorf("step update channel not found for execution %s, dropping step update", executionID)
	}

	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled before step update could be issued: %w", context.Cause(ctx))
	case stepUpdateCh.ch <- stepUpdate:
		return nil
	}
}

func (sucm *stepUpdateManager) len() int64 {
	sucm.mu.RLock()
	defer sucm.mu.RUnlock()
	return int64(len(sucm.m))
}

// Engine handles the lifecycle of a single workflow and its executions.
type Engine struct {
	services.StateMachine
	cma                  custmsg.Labeler
	metrics              workflowsMetricLabeler
	logger               logger.Logger
	registry             core.CapabilitiesRegistry
	workflow             *workflow
	env                  exec.Env
	localNode            capabilities.Node
	executionStates      store.Store
	pendingStepRequests  chan stepRequest
	triggerEvents        chan capabilities.TriggerResponse
	stepUpdatesChMap     stepUpdateManager
	wg                   sync.WaitGroup
	stopCh               services.StopChan
	newWorkerTimeout     time.Duration
	maxExecutionDuration time.Duration

	// testing lifecycle hook to signal when an execution is finished.
	onExecutionFinished func(string)
	// testing lifecycle hook to signal initialization status
	afterInit func(success bool)
	// Used for testing to control the number of retries
	// we'll do when initializing the engine.
	maxRetries int
	// Used for testing to control the retry interval
	// when initializing the engine.
	retryMs int

	maxWorkerLimit int

	clock clockwork.Clock
}

func (e *Engine) Start(_ context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, _ := e.stopCh.NewCtx()

		// spin up monitoring resources
		err := initMonitoringResources()
		if err != nil {
			return fmt.Errorf("could not initialize monitoring resources: %w", err)
		}

		e.wg.Add(e.maxWorkerLimit)
		for i := 0; i < e.maxWorkerLimit; i++ {
			go e.worker(ctx)
		}

		e.wg.Add(1)
		go e.init(ctx)

		return nil
	})
}

// resolveWorkflowCapabilities does the following:
//
// 1. Resolves the underlying capability for each trigger
// 2. Registers each step's capability to this workflow
func (e *Engine) resolveWorkflowCapabilities(ctx context.Context) error {
	//
	// Step 1. Resolve the underlying capability for each trigger
	//
	triggersInitialized := true
	for _, t := range e.workflow.triggers {
		tg, err := e.registry.GetTrigger(ctx, t.ID)
		if err != nil {
			e.logger.With(cIDKey, t.ID).Errorf("failed to get trigger capability: %s", err)
			cerr := e.cma.With(cIDKey, t.ID).SendLogAsCustomMessage(fmt.Sprintf("failed to resolve trigger: %s", err))
			if cerr != nil {
				return cerr
			}
			// we don't immediately return here, since we want to retry all triggers
			// to notify the user of all errors at once.
			triggersInitialized = false
		} else {
			t.trigger = tg
		}
	}
	if !triggersInitialized {
		return &workflowError{reason: "failed to resolve triggers", labels: map[string]string{
			wIDKey: e.workflow.id,
		}}
	}

	// Step 2. Walk the graph and register each step's capability to this workflow
	//
	// This means:
	// - fetching the capability
	// - register the capability to this workflow
	// - initializing the step's executionStrategy
	capabilityRegistrationErr := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
		// The graph contains a dummy step for triggers, but
		// we handle triggers separately since there might be more than one
		// trigger registered to a workflow.
		if s.Ref == workflows.KeywordTrigger {
			return nil
		}

		err := e.initializeCapability(ctx, s)
		if err != nil {
			cerr := e.cma.With(wIDKey, e.workflow.id, sIDKey, s.ID, sRKey, s.Ref).SendLogAsCustomMessage(fmt.Sprintf("failed to initialize capability for step: %s", err))
			if cerr != nil {
				return cerr
			}
			return &workflowError{err: err, reason: "failed to initialize capability for step",
				labels: map[string]string{
					wIDKey: e.workflow.id,
					sIDKey: s.ID,
					sRKey:  s.Ref,
				}}
		}

		return nil
	})

	return capabilityRegistrationErr
}

func (e *Engine) initializeCapability(ctx context.Context, step *step) error {
	// We use varadic err here so that err can be optional, but we assume that
	// its length is either 0 or 1
	newCPErr := func(reason string, errs ...error) *workflowError {
		var err error
		if len(errs) > 0 {
			err = errs[0]
		}

		return &workflowError{reason: reason, err: err, labels: map[string]string{
			wIDKey: e.workflow.id,
			sIDKey: step.ID,
		}}
	}

	// If the capability already exists, that means we've already registered it
	if step.capability != nil {
		return nil
	}

	cp, err := e.registry.Get(ctx, step.ID)
	if err != nil {
		return newCPErr("failed to get capability", err)
	}

	info, err := cp.Info(ctx)
	if err != nil {
		return newCPErr("failed to get capability info", err)
	}

	step.info = info

	// Special treatment for local targets - wrap into a transmission capability
	// If the DON is nil, this is a local target.
	if info.CapabilityType == capabilities.CapabilityTypeTarget && info.IsLocal {
		l := e.logger.With("capabilityID", step.ID)
		l.Debug("wrapping capability in local transmission protocol")
		cp = transmission.NewLocalTargetCapability(
			e.logger,
			step.ID,
			e.localNode,
			cp.(capabilities.TargetCapability),
		)
	}

	// We configure actions, consensus and targets here, and
	// they all satisfy the `CallbackCapability` interface
	cc, ok := cp.(capabilities.ExecutableCapability)
	if !ok {
		return newCPErr("capability does not satisfy CallbackCapability")
	}

	if step.config == nil {
		c, interpErr := exec.FindAndInterpolateEnvVars(step.Config, e.env)
		if interpErr != nil {
			return newCPErr("failed to convert interpolate env vars from config", interpErr)
		}

		config, ok := c.(map[string]any)
		if !ok {
			return newCPErr("failed to convert interpolate env vars from config into map")
		}

		configMap, newMapErr := values.NewMap(config)
		if newMapErr != nil {
			return newCPErr("failed to convert config to values.Map", newMapErr)
		}
		step.config = configMap
	}

	registrationRequest := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID: e.workflow.id,
		},
		Config: step.config,
	}

	err = cc.RegisterToWorkflow(ctx, registrationRequest)
	if err != nil {
		return newCPErr(fmt.Sprintf("failed to register capability to workflow (%+v)", registrationRequest), err)
	}

	step.capability = cc
	return nil
}

// init does the following:
//
//  1. Resolves the LocalDON information
//  2. Resolves the underlying capability for each trigger
//  3. Registers each step's capability to this workflow
//  4. Registers for trigger events now that all capabilities are resolved
//
// Steps 1-3 are retried every 5 seconds until successful.
func (e *Engine) init(ctx context.Context) {
	defer e.wg.Done()

	retryErr := retryable(ctx, e.logger, e.retryMs, e.maxRetries, func() error {
		// first wait for localDON to return a non-error response; this depends
		// on the underlying peerWrapper returning the PeerID.
		node, err := e.registry.LocalNode(ctx)
		if err != nil {
			return fmt.Errorf("failed to get donInfo: %w", err)
		}
		e.localNode = node

		err = e.resolveWorkflowCapabilities(ctx)
		if err != nil {
			return &workflowError{err: err, reason: "failed to resolve workflow capabilities",
				labels: map[string]string{
					wIDKey: e.workflow.id,
				}}
		}
		return nil
	})

	if retryErr != nil {
		e.logger.Errorf("initialization failed: %s", retryErr)
		e.afterInit(false)
		return
	}

	e.logger.Debug("capabilities resolved, resuming in-progress workflows")
	err := e.resumeInProgressExecutions(ctx)
	if err != nil {
		e.logger.Errorf("failed to resume in-progress workflows: %v", err)
	}

	e.logger.Debug("registering triggers")
	for idx, t := range e.workflow.triggers {
		terr := e.registerTrigger(ctx, t, idx)
		if terr != nil {
			e.logger.With(cIDKey, t.ID).Errorf("failed to register trigger: %s", terr)
			cerr := e.cma.With(cIDKey, t.ID).SendLogAsCustomMessage(fmt.Sprintf("failed to register trigger: %s", terr))
			if cerr != nil {
				e.logger.Errorf("failed to send custom message for trigger: %s", terr)
			}
		}
	}

	e.logger.Info("engine initialized")
	e.afterInit(true)
}

var (
	defaultOffset, defaultLimit = 0, 1_000
)

func (e *Engine) resumeInProgressExecutions(ctx context.Context) error {
	wipExecutions, err := e.executionStates.GetUnfinished(ctx, defaultOffset, defaultLimit)
	if err != nil {
		return err
	}

	// TODO: paginate properly
	if len(wipExecutions) >= defaultLimit {
		e.logger.Warnf("possible execution overflow during resumption, work in progress executions: %d >= %d", len(wipExecutions), defaultLimit)
	}

	// Cache the dependents associated with a step.
	// We may have to reprocess many executions, but should only
	// need to calculate the dependents of a step once since
	// they won't change.
	refToDeps := map[string][]*step{}
	for _, execution := range wipExecutions {
		for _, step := range execution.Steps {
			// NOTE: In order to determine what tasks need to be enqueued,
			// we look at any completed steps, and for each dependent,
			// check if they are ready to be enqueued.
			// This will also handle an execution that has stalled immediately on creation,
			// since we always create an execution with an initially completed trigger step.
			if step.Status != store.StatusCompleted {
				continue
			}

			sds, ok := refToDeps[step.Ref]
			if !ok {
				s, err := e.workflow.dependents(step.Ref)
				if err != nil {
					return err
				}

				sds = s
			}

			for _, sd := range sds {
				ch := make(chan store.WorkflowExecutionStep)
				added := e.stepUpdatesChMap.add(execution.ExecutionID, stepUpdateChannel{
					ch:          ch,
					executionID: execution.ExecutionID,
				})
				if added {
					// We trigger the `stepUpdateLoop` for this execution, since the loop is not running atm.
					e.wg.Add(1)
					go e.stepUpdateLoop(ctx, execution.ExecutionID, ch, execution.CreatedAt)
				}
				e.queueIfReady(execution, sd)
			}
		}
	}
	return nil
}

func generateTriggerId(workflowID string, triggerIdx int) string {
	return fmt.Sprintf("wf_%s_trigger_%d", workflowID, triggerIdx)
}

// registerTrigger is used during the initialization phase to bind a trigger to this workflow
func (e *Engine) registerTrigger(ctx context.Context, t *triggerCapability, triggerIdx int) error {
	triggerID := generateTriggerId(e.workflow.id, triggerIdx)

	tc, err := values.NewMap(t.Config)
	if err != nil {
		return err
	}

	t.config.Store(tc)

	triggerRegRequest := capabilities.TriggerRegistrationRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:               e.workflow.id,
			WorkflowOwner:            e.workflow.owner,
			WorkflowName:             e.workflow.name,
			WorkflowDonID:            e.localNode.WorkflowDON.ID,
			WorkflowDonConfigVersion: e.localNode.WorkflowDON.ConfigVersion,
			ReferenceID:              t.Ref,
		},
		Config:    t.config.Load(),
		TriggerID: triggerID,
	}
	eventsCh, err := t.trigger.RegisterTrigger(ctx, triggerRegRequest)
	if err != nil {
		e.metrics.with(cIDKey, t.ID, tIDKey, triggerID).incrementRegisterTriggerFailureCounter(ctx)
		// It's confusing that t.ID is different from triggerID, but
		// t.ID is the capability ID, and triggerID is the trigger ID.
		//
		// The capability ID is globally scoped, whereas the trigger ID
		// is scoped to this workflow.
		//
		// For example, t.ID might be "streams-trigger:network=mainnet@1.0.0"
		// and triggerID might be "wf_123_trigger_0"
		return &workflowError{err: err, reason: fmt.Sprintf("failed to register trigger: %+v", triggerRegRequest),
			labels: map[string]string{
				wIDKey: e.workflow.id,
				cIDKey: t.ID,
				tIDKey: triggerID,
			}}
	}

	e.wg.Add(1)
	go func() {
		defer e.wg.Done()

		for {
			select {
			case <-e.stopCh:
				return
			case event, isOpen := <-eventsCh:
				if !isOpen {
					return
				}

				select {
				case <-e.stopCh:
					return
				case e.triggerEvents <- event:
				}
			}
		}
	}()

	return nil
}

// stepUpdateLoop is a singleton goroutine per `Execution`, and it updates the `executionState` with the outcome of a `step`.
//
// Note: `executionState` is only mutated by this loop directly.
//
// This is important to avoid data races, and any accesses of `executionState` by any other
// goroutine should happen via a `stepRequest` message containing a copy of the latest
// `executionState`.
func (e *Engine) stepUpdateLoop(ctx context.Context, executionID string, stepUpdateCh chan store.WorkflowExecutionStep, workflowCreatedAt *time.Time) {
	defer e.wg.Done()
	lggr := e.logger.With(eIDKey, executionID)
	e.logger.Debugf("running stepUpdateLoop for execution %s", executionID)
	for {
		select {
		case <-ctx.Done():
			lggr.Debug("shutting down stepUpdateLoop")
			return
		case stepUpdate, open := <-stepUpdateCh:
			if !open {
				lggr.Debug("stepUpdate channel closed, shutting down stepUpdateLoop")
				return
			}
			// Executed synchronously to ensure we correctly schedule subsequent tasks.
			e.logger.Debugw(fmt.Sprintf("received step update for execution %s", stepUpdate.ExecutionID),
				eIDKey, stepUpdate.ExecutionID, sRKey, stepUpdate.Ref)
			err := e.handleStepUpdate(ctx, stepUpdate, workflowCreatedAt)
			if err != nil {
				e.logger.Errorf(fmt.Sprintf("failed to update step state: %+v, %s", stepUpdate, err),
					eIDKey, stepUpdate.ExecutionID, sRKey, stepUpdate.Ref)
			}
		}
	}
}

func generateExecutionID(workflowID, eventID string) (string, error) {
	s := sha256.New()
	_, err := s.Write([]byte(workflowID))
	if err != nil {
		return "", err
	}

	_, err = s.Write([]byte(eventID))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(s.Sum(nil)), nil
}

// startExecution kicks off a new workflow execution when a trigger event is received.
func (e *Engine) startExecution(ctx context.Context, executionID string, event *values.Map) error {
	lggr := e.logger.With("event", event, eIDKey, executionID)
	lggr.Debug("executing on a trigger event")
	ec := &store.WorkflowExecution{
		Steps: map[string]*store.WorkflowExecutionStep{
			workflows.KeywordTrigger: {
				Outputs: store.StepOutput{
					Value: event,
				},
				Status:      store.StatusCompleted,
				ExecutionID: executionID,
				Ref:         workflows.KeywordTrigger,
			},
		},
		WorkflowID:  e.workflow.id,
		ExecutionID: executionID,
		Status:      store.StatusStarted,
	}

	dbWex, err := e.executionStates.Add(ctx, ec)
	if err != nil {
		return err
	}

	// Find the tasks we need to fire when a trigger has fired and enqueue them.
	// This consists of a) nodes without a dependency and b) nodes which depend
	// on a trigger
	triggerDependents, err := e.workflow.dependents(workflows.KeywordTrigger)
	if err != nil {
		return err
	}

	ch := make(chan store.WorkflowExecutionStep)
	added := e.stepUpdatesChMap.add(executionID, stepUpdateChannel{
		ch:          ch,
		executionID: executionID,
	})
	if !added {
		// skip this execution since there's already a stepUpdateLoop running for the execution ID
		lggr.Debugf("won't start execution for execution %s, execution was already started", executionID)
		return nil
	}
	e.wg.Add(1)
	go e.stepUpdateLoop(ctx, executionID, ch, dbWex.CreatedAt)

	for _, td := range triggerDependents {
		e.queueIfReady(*ec, td)
	}

	return nil
}

func (e *Engine) handleStepUpdate(ctx context.Context, stepUpdate store.WorkflowExecutionStep, workflowCreatedAt *time.Time) error {
	l := e.logger.With(eIDKey, stepUpdate.ExecutionID, sRKey, stepUpdate.Ref)
	cma := e.cma.With(eIDKey, stepUpdate.ExecutionID, sRKey, stepUpdate.Ref)

	// If we've been executing for too long, let's time the workflow step out and continue.
	if workflowCreatedAt != nil && e.clock.Since(*workflowCreatedAt) > e.maxExecutionDuration {
		l.Info("execution timed out; setting step status to timeout")
		stepUpdate.Status = store.StatusTimeout
	}

	state, err := e.executionStates.UpsertStep(ctx, &stepUpdate)
	if err != nil {
		return err
	}

	workflowIsFullyProcessed, status, err := e.isWorkflowFullyProcessed(ctx, state)
	if err != nil {
		return err
	}

	if workflowIsFullyProcessed {
		switch status {
		case store.StatusTimeout:
			l.Info("execution timed out")
		case store.StatusCompleted:
			l.Info("workflow finished")
		case store.StatusErrored:
			l.Info("execution errored")
		case store.StatusCompletedEarlyExit:
			l.Info("execution terminated early")
			// NOTE: even though this marks the workflow as completed, any branches of the DAG
			// that don't depend on the step that signaled for an early exit will still complete.
			// This is to ensure that any side effects are executed consistently, since otherwise
			// the async nature of the workflow engine would provide no guarantees.
		}
		err = cma.SendLogAsCustomMessage(fmt.Sprintf("execution status: %s", status))
		if err != nil {
			return err
		}
		return e.finishExecution(ctx, state.ExecutionID, status)
	}

	// Finally, since the workflow hasn't timed out or completed, let's
	// check for any dependents that are ready to process.
	stepDependents, err := e.workflow.dependents(stepUpdate.Ref)
	if err != nil {
		return err
	}
	for _, sd := range stepDependents {
		e.queueIfReady(state, sd)
	}

	return nil
}

func (e *Engine) queueIfReady(state store.WorkflowExecution, step *step) {
	// Check if all dependencies are completed for the current step
	var waitingOnDependencies bool
	for _, dr := range step.Vertex.Dependencies {
		stepState, ok := state.Steps[dr]
		if !ok {
			waitingOnDependencies = true
			continue
		}

		// Unless the dependency is complete,
		// we'll mark waitingOnDependencies = true.
		// This includes cases where one of the dependent
		// steps has errored, since that means we shouldn't
		// schedule the step for execution.
		if stepState.Status != store.StatusCompleted {
			waitingOnDependencies = true
		}
	}

	// If all dependencies are completed, enqueue the step.
	if !waitingOnDependencies {
		e.logger.With(sRKey, step.Ref, eIDKey, state.ExecutionID, "state", copyState(state)).
			Debug("step request enqueued")
		e.pendingStepRequests <- stepRequest{
			state:   copyState(state),
			stepRef: step.Ref,
		}
	}
}

func (e *Engine) finishExecution(ctx context.Context, executionID string, status string) error {
	e.logger.With(eIDKey, executionID, "status", status).Info("finishing execution")
	metrics := e.metrics.with(eIDKey, executionID, "status", status)
	err := e.executionStates.UpdateStatus(ctx, executionID, status)
	if err != nil {
		return err
	}

	execState, err := e.executionStates.Get(ctx, executionID)
	if err != nil {
		return err
	}

	executionDuration := execState.FinishedAt.Sub(*execState.CreatedAt).Milliseconds()

	e.stepUpdatesChMap.remove(executionID)
	metrics.updateTotalWorkflowsGauge(ctx, e.stepUpdatesChMap.len())
	metrics.updateWorkflowExecutionLatencyGauge(ctx, executionDuration)
	e.onExecutionFinished(executionID)
	return nil
}

// worker is responsible for:
//   - handling a `pendingStepRequests`
//   - starting a new execution when a trigger emits a message on `triggerEvents`
func (e *Engine) worker(ctx context.Context) {
	defer e.wg.Done()

	for {
		select {
		case pendingStepRequest := <-e.pendingStepRequests:
			e.workerForStepRequest(ctx, pendingStepRequest)
		case resp, isOpen := <-e.triggerEvents:
			if !isOpen {
				e.logger.Error("trigger events channel is no longer open, skipping")
				continue
			}

			if resp.Err != nil {
				e.logger.Errorf("trigger event was an error %v; not executing", resp.Err)
				continue
			}

			te := resp.Event

			if te.ID == "" {
				e.logger.With(tIDKey, te.TriggerType).Error("trigger event ID is empty; not executing")
				continue
			}

			executionID, err := generateExecutionID(e.workflow.id, te.ID)
			if err != nil {
				e.logger.With(tIDKey, te.ID).Errorf("could not generate execution ID: %v", err)
				continue
			}

			err = e.startExecution(ctx, executionID, resp.Event.Outputs)
			if err != nil {
				e.logger.With(eIDKey, executionID).Errorf("failed to start execution: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) workerForStepRequest(ctx context.Context, msg stepRequest) {
	// Instantiate a child logger; in addition to the WorkflowID field the workflow
	// logger will already have, this adds the `stepRef` and `executionID`
	l := e.logger.With(sRKey, msg.stepRef, eIDKey, msg.state.ExecutionID)
	cma := e.cma.With(sRKey, msg.stepRef, eIDKey, msg.state.ExecutionID)

	l.Debug("executing on a step event")
	stepState := &store.WorkflowExecutionStep{
		Outputs:     store.StepOutput{},
		ExecutionID: msg.state.ExecutionID,
		Ref:         msg.stepRef,
	}

	// TODO ks-462 inputs
	err := cma.SendLogAsCustomMessage("executing step")
	if err != nil {
		return
	}
	inputs, outputs, err := e.executeStep(ctx, msg)
	var stepStatus string
	switch {
	case errors.Is(capabilities.ErrStopExecution, err):
		lmsg := "step executed successfully with a termination"
		l.Info(lmsg)
		cmErr := cma.SendLogAsCustomMessage(lmsg)
		if cmErr != nil {
			l.Errorf("failed to send custom message with msg: %s", lmsg)
		}
		stepStatus = store.StatusCompletedEarlyExit
	case err != nil:
		lmsg := "step executed successfully with a termination"
		l.Errorf("error executing step request: %s", err)
		cmErr := cma.SendLogAsCustomMessage(fmt.Sprintf("error executing step request: %s", err))
		if cmErr != nil {
			l.Errorf("failed to send custom message with msg: %s", lmsg)
		}
		stepStatus = store.StatusErrored
	default:
		lmsg := "step executed successfully with a termination"
		l.With("outputs", outputs).Info("step executed successfully")
		// TODO ks-462 emit custom message with outputs
		cmErr := cma.SendLogAsCustomMessage("step executed successfully")
		if cmErr != nil {
			l.Errorf("failed to send custom message with msg: %s", lmsg)
		}
		stepStatus = store.StatusCompleted
	}

	stepState.Status = stepStatus
	stepState.Outputs.Value = outputs
	stepState.Outputs.Err = err
	stepState.Inputs = inputs

	// Let's try and emit the stepUpdate.
	// If the context is canceled, we'll just drop the update.
	// This means the engine is shutting down and the
	// receiving loop may not pick up any messages we emit.
	// Note: When full persistence support is added, any hanging steps
	// like this one will get picked up again and will be reprocessed.
	l.Debugf("trying to send step state update for execution %s with status %s", stepState.ExecutionID, stepStatus)
	err = e.stepUpdatesChMap.send(ctx, stepState.ExecutionID, *stepState)
	if err != nil {
		l.Errorf("failed to issue step state update; error %v", err)
		return
	}
	l.Debugf("sent step state update for execution %s with status %s", stepState.ExecutionID, stepStatus)
}

func merge(baseConfig *values.Map, overrideConfig *values.Map) *values.Map {
	m := values.EmptyMap()

	for k, v := range baseConfig.Underlying {
		m.Underlying[k] = v
	}

	for k, v := range overrideConfig.Underlying {
		m.Underlying[k] = v
	}

	return m
}

func (e *Engine) configForStep(ctx context.Context, executionID string, step *step) (*values.Map, error) {
	ID := step.info.ID

	// If the capability info is missing a DON, then
	// the capability is local, and we should use the localNode's DON ID.
	var donID uint32
	if !step.info.IsLocal {
		donID = step.info.DON.ID
	} else {
		donID = e.localNode.WorkflowDON.ID
	}

	capConfig, err := e.registry.ConfigForCapability(ctx, ID, donID)
	if err != nil {
		e.logger.Warnw(fmt.Sprintf("could not retrieve config from remote registry: %s", err), "executionID", executionID, "capabilityID", ID)
		return step.config, nil
	}

	if capConfig.DefaultConfig == nil {
		return step.config, nil
	}

	// Merge the configs with registry config overriding the step config.  This is because
	// some config fields are sensitive and could affect the safe running of the capability,
	// so we avoid user provided values by overriding them with config from the capabilities registry.
	return merge(step.config, capConfig.DefaultConfig), nil
}

// executeStep executes the referenced capability within a step and returns the result.
func (e *Engine) executeStep(ctx context.Context, msg stepRequest) (*values.Map, values.Value, error) {
	step, err := e.workflow.Vertex(msg.stepRef)
	if err != nil {
		return nil, nil, err
	}

	var inputs any
	if step.Inputs.OutputRef != "" {
		inputs = step.Inputs.OutputRef
	} else {
		inputs = step.Inputs.Mapping
	}

	i, err := exec.FindAndInterpolateAllKeys(inputs, msg.state)
	if err != nil {
		return nil, nil, err
	}

	inputsMap, err := values.NewMap(i.(map[string]any))
	if err != nil {
		return nil, nil, err
	}

	config, err := e.configForStep(ctx, msg.state.ExecutionID, step)
	if err != nil {
		return nil, nil, err
	}

	tr := capabilities.CapabilityRequest{
		Inputs: inputsMap,
		Config: config,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:               msg.state.WorkflowID,
			WorkflowExecutionID:      msg.state.ExecutionID,
			WorkflowOwner:            e.workflow.owner,
			WorkflowName:             e.workflow.name,
			WorkflowDonID:            e.localNode.WorkflowDON.ID,
			WorkflowDonConfigVersion: e.localNode.WorkflowDON.ConfigVersion,
			ReferenceID:              msg.stepRef,
		},
	}

	e.metrics.incrementCapabilityInvocationCounter(ctx)
	output, err := step.capability.Execute(ctx, tr)
	if err != nil {
		return inputsMap, nil, err
	}

	return inputsMap, output.Value, err
}

func (e *Engine) deregisterTrigger(ctx context.Context, t *triggerCapability, triggerIdx int) error {
	deregRequest := capabilities.TriggerRegistrationRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:               e.workflow.id,
			WorkflowDonID:            e.localNode.WorkflowDON.ID,
			WorkflowDonConfigVersion: e.localNode.WorkflowDON.ConfigVersion,
			WorkflowName:             e.workflow.name,
			WorkflowOwner:            e.workflow.owner,
			ReferenceID:              t.Ref,
		},
		TriggerID: generateTriggerId(e.workflow.id, triggerIdx),
		Config:    t.config.Load(),
	}

	// if t.trigger == nil, then we haven't initialized the workflow
	// yet, and can safely consider the trigger deregistered with
	// no further action.
	if t.trigger != nil {
		return t.trigger.UnregisterTrigger(ctx, deregRequest)
	}

	return nil
}

func (e *Engine) isWorkflowFullyProcessed(ctx context.Context, state store.WorkflowExecution) (bool, string, error) {
	statuses := map[string]string{}
	// we need to first propagate the status of the errored status if it exists...
	err := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
		stateStep, ok := state.Steps[s.Ref]
		if !ok {
			// The step not existing on the state means that it has not been processed yet.
			// So ignore it.
			return nil
		}
		statuses[s.Ref] = stateStep.Status
		switch stateStep.Status {
		// For each step with any of the following statuses, propagate the statuses to its dependants
		// since they will not be executed.
		case store.StatusErrored, store.StatusCompletedEarlyExit, store.StatusTimeout:
			// Let's properly propagate the status to all dependents, not just direct dependents.
			queue := []string{s.Ref}
			for len(queue) > 0 {
				current := queue[0] // Grab the current step reference
				queue = queue[1:]   // Remove it from the queue

				// Find the dependents for the current step reference
				dependents, err := e.workflow.dependents(current)
				if err != nil {
					return err
				}

				// Propagate the status to all direct dependents
				// With no direct dependents, it will go to the next step reference in the queue.
				for _, sd := range dependents {
					if _, dependentProcessed := statuses[sd.Ref]; !dependentProcessed {
						statuses[sd.Ref] = stateStep.Status
						// Queue the dependent for to be processed later
						queue = append(queue, sd.Ref)
					}
				}
			}
			e.metrics.incrementTotalWorkflowStepErrorsCounter(ctx)
		}
		return nil
	})
	if err != nil {
		return false, "", err
	}

	workflowProcessed := true
	// Let's validate whether the workflow has been fully processed.
	err = e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
		// If the step is not part of the state, it is a pending step
		// so we should consider the workflow as not fully processed.
		if _, ok := statuses[s.Ref]; !ok {
			workflowProcessed = false
		}
		return nil
	})
	if err != nil {
		return false, "", err
	}

	if !workflowProcessed {
		return workflowProcessed, "", nil
	}

	var hasErrored, hasTimedOut, hasCompletedEarlyExit bool
	// Let's determine the status of the workflow.
	for _, status := range statuses {
		switch status {
		case store.StatusErrored:
			hasErrored = true
		case store.StatusTimeout:
			hasTimedOut = true
		case store.StatusCompletedEarlyExit:
			hasCompletedEarlyExit = true
		}
	}

	// The `errored` status has precedence over the other statuses to be returned, based on occurrence.
	// Status precedence: `errored` -> `timed_out` -> `completed_early_exit` -> `completed`.
	if hasErrored {
		return workflowProcessed, store.StatusErrored, nil
	}
	if hasTimedOut {
		return workflowProcessed, store.StatusTimeout, nil
	}
	if hasCompletedEarlyExit {
		return workflowProcessed, store.StatusCompletedEarlyExit, nil
	}
	return workflowProcessed, store.StatusCompleted, nil
}

func (e *Engine) Close() error {
	return e.StopOnce("Engine", func() error {
		e.logger.Info("shutting down engine")
		ctx := context.Background()
		// To shut down the engine, we'll start by deregistering
		// any triggers to ensure no new executions are triggered,
		// then we'll close down any background goroutines,
		// and finally, we'll deregister any workflow steps.
		for idx, t := range e.workflow.triggers {
			err := e.deregisterTrigger(ctx, t, idx)
			if err != nil {
				return err
			}
		}

		close(e.stopCh)
		e.wg.Wait()

		err := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
			if s.Ref == workflows.KeywordTrigger {
				return nil
			}

			reg := capabilities.UnregisterFromWorkflowRequest{
				Metadata: capabilities.RegistrationMetadata{
					WorkflowID: e.workflow.id,
				},
				Config: s.config,
			}

			// if capability is nil, then we haven't initialized
			// the workflow yet and can safely consider it deregistered
			// with no further action.
			if s.capability == nil {
				return nil
			}

			innerErr := s.capability.UnregisterFromWorkflow(ctx, reg)
			if innerErr != nil {
				return &workflowError{err: innerErr,
					reason: fmt.Sprintf("failed to unregister capability from workflow: %+v", reg),
					labels: map[string]string{
						wIDKey: e.workflow.id,
						sIDKey: s.ID,
						sRKey:  s.Ref,
					}}
			}

			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
}

type Config struct {
	Workflow             sdk.WorkflowSpec
	WorkflowID           string
	WorkflowOwner        string
	WorkflowName         string
	Lggr                 logger.Logger
	Registry             core.CapabilitiesRegistry
	MaxWorkerLimit       int
	QueueSize            int
	NewWorkerTimeout     time.Duration
	MaxExecutionDuration time.Duration
	Store                store.Store
	Config               []byte
	Binary               []byte

	// For testing purposes only
	maxRetries          int
	retryMs             int
	afterInit           func(success bool)
	onExecutionFinished func(weid string)
	clock               clockwork.Clock
}

const (
	defaultWorkerLimit          = 100
	defaultQueueSize            = 100000
	defaultNewWorkerTimeout     = 2 * time.Second
	defaultMaxExecutionDuration = 10 * time.Minute
)

func NewEngine(cfg Config) (engine *Engine, err error) {
	if cfg.Store == nil {
		return nil, &workflowError{reason: "store is nil",
			labels: map[string]string{
				wIDKey: cfg.WorkflowID,
			},
		}
	}

	if cfg.MaxWorkerLimit == 0 {
		cfg.MaxWorkerLimit = defaultWorkerLimit
	}

	if cfg.QueueSize == 0 {
		cfg.QueueSize = defaultQueueSize
	}

	if cfg.NewWorkerTimeout == 0 {
		cfg.NewWorkerTimeout = defaultNewWorkerTimeout
	}

	if cfg.MaxExecutionDuration == 0 {
		cfg.MaxExecutionDuration = defaultMaxExecutionDuration
	}

	if cfg.retryMs == 0 {
		cfg.retryMs = 5000
	}

	if cfg.afterInit == nil {
		cfg.afterInit = func(success bool) {}
	}

	if cfg.onExecutionFinished == nil {
		cfg.onExecutionFinished = func(weid string) {}
	}

	if cfg.clock == nil {
		cfg.clock = clockwork.NewRealClock()
	}

	// TODO: validation of the workflow spec
	// We'll need to check, among other things:
	// - that there are no step `ref` called `trigger` as this is reserved for any triggers
	// - that there are no duplicate `ref`s
	// - that the `ref` for any triggers is empty -- and filled in with `trigger`
	// - that the resulting graph is strongly connected (i.e. no disjointed subgraphs exist)
	// - etc.

	workflow, err := Parse(cfg.Workflow)
	if err != nil {
		return nil, err
	}

	workflow.id = cfg.WorkflowID
	workflow.owner = cfg.WorkflowOwner
	workflow.name = hex.EncodeToString([]byte(cfg.WorkflowName))

	engine = &Engine{
		logger:   cfg.Lggr.Named("WorkflowEngine").With("workflowID", cfg.WorkflowID),
		cma:      custmsg.NewLabeler().With(wIDKey, cfg.WorkflowID, woIDKey, cfg.WorkflowOwner, wnKey, workflow.name),
		metrics:  workflowsMetricLabeler{metrics.NewLabeler().With(wIDKey, cfg.WorkflowID, woIDKey, cfg.WorkflowOwner, wnKey, workflow.name)},
		registry: cfg.Registry,
		workflow: workflow,
		env: exec.Env{
			Config: cfg.Config,
			Binary: cfg.Binary,
		},
		executionStates:      cfg.Store,
		pendingStepRequests:  make(chan stepRequest, cfg.QueueSize),
		stepUpdatesChMap:     stepUpdateManager{m: map[string]stepUpdateChannel{}},
		triggerEvents:        make(chan capabilities.TriggerResponse),
		stopCh:               make(chan struct{}),
		newWorkerTimeout:     cfg.NewWorkerTimeout,
		maxExecutionDuration: cfg.MaxExecutionDuration,
		onExecutionFinished:  cfg.onExecutionFinished,
		afterInit:            cfg.afterInit,
		maxRetries:           cfg.maxRetries,
		retryMs:              cfg.retryMs,
		maxWorkerLimit:       cfg.MaxWorkerLimit,
		clock:                cfg.clock,
	}

	return engine, nil
}

type workflowError struct {
	labels map[string]string
	// err is the underlying error that caused this error
	err error
	// reason is a human-readable string that describes the error
	reason string
}

func (e *workflowError) Error() string {
	errStr := ""
	if e.err != nil {
		if e.reason != "" {
			errStr = fmt.Sprintf("%s: %v", e.reason, e.err)
		} else {
			errStr = e.err.Error()
		}
	} else {
		errStr = e.reason
	}

	// prefix the error with the labels
	for _, label := range orderedLabelKeys {
		// This will silently ignore any labels that are not present in the map
		// are we ok with this?
		if value, ok := e.labels[label]; ok {
			errStr = fmt.Sprintf("%s %s: %s", label, value, errStr)
		}
	}

	return errStr
}
