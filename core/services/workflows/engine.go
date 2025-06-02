package workflows

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/jonboulle/clockwork"

	"github.com/smartcontractkit/chainlink-common/pkg/aggregation"
	"github.com/smartcontractkit/chainlink-common/pkg/capabilities"
	"github.com/smartcontractkit/chainlink-common/pkg/custmsg"
	"github.com/smartcontractkit/chainlink-common/pkg/metrics"
	"github.com/smartcontractkit/chainlink-common/pkg/services"
	"github.com/smartcontractkit/chainlink-common/pkg/types/core"
	"github.com/smartcontractkit/chainlink-common/pkg/values"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/exec"
	"github.com/smartcontractkit/chainlink-common/pkg/workflows/sdk"
	billing "github.com/smartcontractkit/chainlink-protos/billing/go"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/transmission"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/platform"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/events"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/metering"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/ratelimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/store"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/syncerlimiter"
	"github.com/smartcontractkit/chainlink/v2/core/services/workflows/types"
)

const (
	fifteenMinutesSec            = 15 * 60
	reservedFieldNameStepTimeout = "cre_step_timeout"
	maxStepTimeoutOverrideSec    = 10 * 60 // 10 minutes
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
	sucm.mu.Lock()
	defer sucm.mu.Unlock()
	if _, ok := sucm.m[executionID]; ok {
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
	sucm.mu.Lock()
	defer sucm.mu.Unlock()
	stepUpdateCh, ok := sucm.m[executionID]

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

type SecretsFor func(ctx context.Context, workflowOwner, hexWorkflowName, decodedWorkflowName, workflowID string) (map[string]string, error)

type BillingClient interface {
	SubmitWorkflowReceipt(context.Context, *billing.SubmitWorkflowReceiptRequest) (*billing.SubmitWorkflowReceiptResponse, error)
}

// Engine handles the lifecycle of a single workflow and its executions.
type Engine struct {
	services.StateMachine
	cma                  custmsg.MessageEmitter
	metrics              workflowsMetricLabeler
	logger               logger.Logger
	registry             core.CapabilitiesRegistry
	workflow             *workflow
	secretsFetcher       SecretsFor
	env                  exec.Env
	localNode            atomic.Pointer[capabilities.Node]
	executionsStore      store.Store
	pendingStepRequests  chan stepRequest
	triggerEvents        chan capabilities.TriggerResponse
	stepUpdatesChMap     stepUpdateManager
	wg                   sync.WaitGroup
	stopCh               services.StopChan
	newWorkerTimeout     time.Duration
	maxExecutionDuration time.Duration
	heartbeatCadence     time.Duration
	stepTimeoutDuration  time.Duration
	billingClient        BillingClient

	// testing lifecycle hook to signal when an execution is finished.
	onExecutionFinished func(string)
	// testing lifecycle hook to signal initialization status
	afterInit func(success bool)

	// testing lifecycle hook to signal the execution was rate limited
	onRateLimit func(string)

	// Used for testing to control the number of retries
	// we'll do when initializing the engine.
	maxRetries int
	// Used for testing to control the retry interval
	// when initializing the engine.
	retryMs int

	maxWorkerLimit int

	clock          clockwork.Clock
	ratelimiter    *ratelimiter.RateLimiter
	workflowLimits *syncerlimiter.Limits
	meterReports   *metering.Reports
}

func (e *Engine) Start(_ context.Context) error {
	return e.StartOnce("Engine", func() error {
		// create a new context, since the one passed in via Start is short-lived.
		ctx, _ := e.stopCh.NewCtx()

		// validate if adding another workflow would exceed either the global or per owner engine count limit
		ownerAllow, globalAllow := e.workflowLimits.Allow(e.workflow.owner)
		if !globalAllow {
			e.metrics.incrementWorkflowLimitGlobalCounter(ctx)
			logCustMsg(ctx, e.cma.With(platform.KeyWorkflowID, e.workflow.id, platform.KeyWorkflowOwner, e.workflow.owner), types.ErrGlobalWorkflowCountLimitReached.Error(), e.logger)
			return types.ErrGlobalWorkflowCountLimitReached
		}

		if !ownerAllow {
			e.metrics.incrementWorkflowLimitPerOwnerCounter(ctx)
			logCustMsg(ctx, e.cma.With(platform.KeyWorkflowID, e.workflow.id, platform.KeyWorkflowOwner, e.workflow.owner), types.ErrPerOwnerWorkflowCountLimitReached.Error(), e.logger)
			return types.ErrPerOwnerWorkflowCountLimitReached
		}

		e.metrics.incrementWorkflowInitializationCounter(ctx)

		e.wg.Add(e.maxWorkerLimit)
		for i := 0; i < e.maxWorkerLimit; i++ {
			go e.worker(ctx)
		}

		e.wg.Add(1)
		go e.init(ctx)

		e.wg.Add(1)
		go e.heartbeat(ctx)

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
			log := e.logger.With(platform.KeyCapabilityID, t.ID)
			log.Errorf("failed to get trigger capability: %s", err)
			logCustMsg(ctx, e.cma.With(platform.KeyCapabilityID, t.ID), fmt.Sprintf("failed to resolve trigger: %s", err), log)
			// we don't immediately return here, since we want to retry all triggers
			// to notify the user of all errors at once.
			triggersInitialized = false
		} else {
			t.trigger = tg
		}
	}
	if !triggersInitialized {
		return &workflowError{reason: "failed to resolve triggers", labels: map[string]string{
			platform.KeyWorkflowID: e.workflow.id,
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
			logCustMsg(
				ctx,
				e.cma.With(platform.KeyWorkflowID, e.workflow.id, platform.KeyStepID, s.ID, platform.KeyStepRef, s.Ref),
				fmt.Sprintf("failed to initialize capability for step: %s", err),
				e.logger,
			)
			return &workflowError{err: err, reason: "failed to initialize capability for step",
				labels: map[string]string{
					platform.KeyWorkflowID: e.workflow.id,
					platform.KeyStepID:     s.ID,
					platform.KeyStepRef:    s.Ref,
				}}
		}

		return nil
	})

	return capabilityRegistrationErr
}

func (e *Engine) initializeCapability(ctx context.Context, step *step) error {
	l := e.logger.With("capabilityID", step.ID)

	// We use varadic err here so that err can be optional, but we assume that
	// its length is either 0 or 1
	newCPErr := func(reason string, errs ...error) *workflowError {
		var err error
		if len(errs) > 0 {
			err = errs[0]
		}

		return &workflowError{reason: reason, err: err, labels: map[string]string{
			platform.KeyWorkflowID: e.workflow.id,
			platform.KeyStepID:     step.ID,
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
		l.Debug("wrapping capability in local transmission protocol")
		cp = transmission.NewLocalTargetCapability(
			e.logger,
			step.ID,
			*e.localNode.Load(),
			cp.(capabilities.TargetCapability),
		)
	}

	// We configure actions, consensus and targets here, and
	// they all satisfy the `CallbackCapability` interface
	cc, ok := cp.(capabilities.ExecutableCapability)
	if !ok {
		return newCPErr("capability does not satisfy CallbackCapability")
	}

	stepConfig, err := e.configForStep(ctx, l, step)
	if err != nil {
		return newCPErr("failed to get config for step", err)
	}

	registrationRequest := capabilities.RegisterToWorkflowRequest{
		Metadata: capabilities.RegistrationMetadata{
			WorkflowID:    e.workflow.id,
			WorkflowOwner: e.workflow.owner,
			ReferenceID:   step.Vertex.Ref,
		},
		Config: stepConfig,
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

	retryErr := internal.RunWithRetries(ctx, e.logger, time.Millisecond*time.Duration(e.retryMs), e.maxRetries, func() error {
		// first wait for localDON to return a non-error response; this depends
		// on the underlying peerWrapper returning the PeerID.
		node, err := e.registry.LocalNode(ctx)
		if err != nil {
			return fmt.Errorf("failed to get donInfo: %w", err)
		}
		e.localNode.Store(&node)

		err = e.resolveWorkflowCapabilities(ctx)
		if err != nil {
			return &workflowError{err: err, reason: "failed to resolve workflow capabilities",
				labels: map[string]string{
					platform.KeyWorkflowID: e.workflow.id,
				}}
		}
		return nil
	})

	if retryErr != nil {
		e.logger.Errorf("initialization failed: %s", retryErr)
		logCustMsg(ctx, e.cma, fmt.Sprintf("workflow registration failed: %s", retryErr), e.logger)
		e.afterInit(false)
		return
	}

	e.logger.Debug("capabilities resolved")

	e.logger.Debug("registering triggers")
	for idx, t := range e.workflow.triggers {
		terr := e.registerTrigger(ctx, t, idx)
		if terr != nil {
			log := e.logger.With(platform.KeyCapabilityID, t.ID)
			log.Errorf("failed to register trigger: %s", terr)
			logCustMsg(ctx, e.cma.With(platform.KeyCapabilityID, t.ID), fmt.Sprintf("failed to register trigger: %s", terr), log)
		}
	}

	e.logger.Info("engine initialized")
	logCustMsg(ctx, e.cma, "workflow registered", e.logger)
	e.metrics.incrementWorkflowRegisteredCounter(ctx)
	e.afterInit(true)
}

func generateTriggerID(workflowID string, triggerIdx int) string {
	return fmt.Sprintf("wf_%s_trigger_%d", workflowID, triggerIdx)
}

// registerTrigger is used during the initialization phase to bind a trigger to this workflow
func (e *Engine) registerTrigger(ctx context.Context, t *triggerCapability, triggerIdx int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	triggerID := generateTriggerID(e.workflow.id, triggerIdx)

	tc, err := values.NewMap(t.Config)
	if err != nil {
		return err
	}

	t.config = tc

	ln := e.localNode.Load()
	triggerRegRequest := capabilities.TriggerRegistrationRequest{
		Metadata: capabilities.RequestMetadata{
			WorkflowID:               e.workflow.id,
			WorkflowOwner:            e.workflow.owner,
			WorkflowName:             e.workflow.name.Hex(),
			WorkflowDonID:            ln.WorkflowDON.ID,
			WorkflowDonConfigVersion: ln.WorkflowDON.ConfigVersion,
			ReferenceID:              t.Ref,
			DecodedWorkflowName:      e.workflow.name.String(),
		},
		Config:    t.config,
		TriggerID: triggerID,
	}
	eventsCh, err := t.trigger.RegisterTrigger(ctx, triggerRegRequest)
	if err != nil {
		e.metrics.with(platform.KeyTriggerID, triggerID).incrementRegisterTriggerFailureCounter(ctx)
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
				platform.KeyWorkflowID:   e.workflow.id,
				platform.KeyCapabilityID: t.ID,
				platform.KeyTriggerID:    triggerID,
			}}
	}

	// mark the trigger as successfully registered
	t.registered = true

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
	lggr := e.logger.With(platform.KeyWorkflowExecutionID, executionID)
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
			e.logger.Debugw("received step update for execution "+stepUpdate.ExecutionID,
				platform.KeyWorkflowExecutionID, stepUpdate.ExecutionID, platform.KeyStepRef, stepUpdate.Ref)
			err := e.handleStepUpdate(ctx, stepUpdate, workflowCreatedAt)
			if err != nil {
				e.logger.Errorf(fmt.Sprintf("failed to update step state: %+v, %s", stepUpdate, err),
					platform.KeyWorkflowExecutionID, stepUpdate.ExecutionID, platform.KeyStepRef, stepUpdate.Ref)
			}
		}
	}
}

// startExecution kicks off a new workflow execution when a trigger event is received.
func (e *Engine) startExecution(ctx context.Context, executionID string, triggerID string, event *values.Map) error {
	e.meterReports.Add(executionID, metering.NewReport(e.logger))

	err := events.EmitExecutionStartedEvent(ctx, e.cma, triggerID, executionID)
	if err != nil {
		e.logger.Errorf("failed to emit execution started event: %+v", err)
	}

	lggr := e.logger.With("event", event, platform.KeyWorkflowExecutionID, executionID)
	lggr.Debug("executing on a trigger event")
	workflowExecution, err := e.executionsStore.Add(ctx, map[string]*store.WorkflowExecutionStep{
		workflows.KeywordTrigger: {
			Outputs: store.StepOutput{
				Value: event,
			},
			Status:      store.StatusCompleted,
			ExecutionID: executionID,
			Ref:         workflows.KeywordTrigger,
		}}, executionID, e.workflow.id, store.StatusStarted)
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
	go e.stepUpdateLoop(ctx, executionID, ch, workflowExecution.CreatedAt)

	for _, td := range triggerDependents {
		e.queueIfReady(workflowExecution, td)
	}

	return nil
}

func (e *Engine) handleStepUpdate(ctx context.Context, stepUpdate store.WorkflowExecutionStep, workflowCreatedAt *time.Time) error {
	l := e.logger.With(platform.KeyWorkflowExecutionID, stepUpdate.ExecutionID, platform.KeyStepRef, stepUpdate.Ref)
	cma := e.cma.With(platform.KeyWorkflowExecutionID, stepUpdate.ExecutionID, platform.KeyStepRef, stepUpdate.Ref)

	// If we've been executing for too long, let's time the workflow step out and continue.
	if workflowCreatedAt != nil && e.clock.Since(*workflowCreatedAt) > e.maxExecutionDuration {
		l.Info("execution timed out; setting step status to timeout")
		stepUpdate.Status = store.StatusTimeout
	}

	state, err := e.executionsStore.UpsertStep(ctx, &stepUpdate)
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
			e.metrics.incrementTotalWorkflowStepErrorsCounter(ctx)
		case store.StatusCompletedEarlyExit:
			l.Info("execution terminated early")
			// NOTE: even though this marks the workflow as completed, any branches of the DAG
			// that don't depend on the step that signaled for an early exit will still complete.
			// This is to ensure that any side effects are executed consistently, since otherwise
			// the async nature of the workflow engine would provide no guarantees.
		}

		logCustMsg(ctx, cma, "execution status: "+status, l)

		// this case is only for resuming executions and should be updated when metering is added to save execution state
		if _, ok := e.meterReports.Get(stepUpdate.ExecutionID); !ok {
			e.meterReports.Add(stepUpdate.ExecutionID, metering.NewReport(e.logger))
		}

		return e.finishExecution(ctx, cma, state.ExecutionID, status)
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
		e.logger.With(platform.KeyStepRef, step.Ref, platform.KeyWorkflowExecutionID, state.ExecutionID, "state", state.DeepCopy()).
			Debug("step request enqueued")
		e.pendingStepRequests <- stepRequest{
			state:   state.DeepCopy(),
			stepRef: step.Ref,
		}
	}
}

func (e *Engine) finishExecution(ctx context.Context, cma custmsg.MessageEmitter, executionID string, status string) error {
	l := e.logger.With(platform.KeyWorkflowExecutionID, executionID, "status", status)

	l.Info("finishing execution")

	execState, err := e.executionsStore.FinishExecution(ctx, executionID, status)
	if err != nil {
		return fmt.Errorf("failed to mark execution as finished: %w", err)
	}

	report, exists := e.meterReports.Get(executionID)
	if exists {
		// send metering report to beholder
		if err = events.EmitMeteringReport(ctx, cma, report.Message()); err != nil {
			e.metrics.incrementWorkflowMissingMeteringReport(ctx)
			l.Warn(fmt.Sprintf("metering report send to beholder error %s", err))
		}

		// send metering report to billing if billing client is not nil
		if err = e.sendReportToBilling(ctx, report, e.workflow.id, executionID); err != nil {
			e.metrics.incrementWorkflowMissingMeteringReport(ctx)
			l.Warn(fmt.Sprintf("metering report send to billing error %s", err))
		}
	}

	// clean all per execution state trackers
	e.stepUpdatesChMap.remove(executionID)
	e.meterReports.Delete(executionID)

	executionDuration := int64(execState.FinishedAt.Sub(*execState.CreatedAt).Seconds())
	switch status {
	case store.StatusCompleted:
		e.metrics.updateWorkflowCompletedDurationHistogram(ctx, executionDuration)
	case store.StatusCompletedEarlyExit:
		e.metrics.updateWorkflowEarlyExitDurationHistogram(ctx, executionDuration)
	case store.StatusErrored:
		e.metrics.updateWorkflowErrorDurationHistogram(ctx, executionDuration)
	case store.StatusTimeout:
		// should expect the same values unless the timeout is adjusted.
		// using histogram as it gives count of executions for free
		e.metrics.updateWorkflowTimeoutDurationHistogram(ctx, executionDuration)
	}

	if executionDuration > fifteenMinutesSec {
		logCustMsg(ctx, cma, fmt.Sprintf("execution duration exceeded 15 minutes: %d (seconds)", executionDuration), l)
		l.Warnf("execution duration exceeded 15 minutes: %d (seconds)", executionDuration)
	}

	logCustMsg(ctx, cma, fmt.Sprintf("execution duration: %d (seconds)", executionDuration), l)
	l.Infof("execution duration: %d (seconds)", executionDuration)
	err = events.EmitExecutionFinishedEvent(ctx, cma, status, executionID)
	if err != nil {
		e.logger.Errorf("failed to emit execution finished event: %+v", err)
	}

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
				logCustMsg(ctx, e.cma, fmt.Sprintf("failed to resolve trigger: %s", resp.Err), e.logger)
				continue
			}

			te := resp.Event

			if te.ID == "" {
				e.logger.With(platform.KeyTriggerID, te.TriggerType).Error("trigger event ID is empty; not executing")
				continue
			}

			executionID, err := types.GenerateExecutionID(e.workflow.id, te.ID)
			if err != nil {
				e.logger.With(platform.KeyTriggerID, te.ID).Errorf("could not generate execution ID: %v", err)
				continue
			}

			senderAllowed, globalAllowed := e.ratelimiter.Allow(e.workflow.owner)
			if !senderAllowed {
				e.onRateLimit(executionID)
				e.logger.With(platform.KeyWorkflowID, e.workflow.id, platform.KeyWorkflowOwner, e.workflow.owner, platform.KeyWorkflowExecutionID, executionID).Errorf("failed to start execution: per sender rate limit exceeded")
				logCustMsg(ctx, e.cma.With(platform.KeyCapabilityID, te.ID), "failed to start execution: per sender rate limit exceeded", e.logger)
				e.metrics.with(platform.KeyTriggerID, te.ID).incrementWorkflowExecutionRateLimitPerUserCounter(ctx)
				continue
			}
			if !globalAllowed {
				e.onRateLimit(executionID)
				e.logger.With(platform.KeyWorkflowID, e.workflow.id, platform.KeyWorkflowOwner, e.workflow.owner, platform.KeyWorkflowExecutionID, executionID).Errorf("failed to start execution: global rate limit exceeded")
				logCustMsg(ctx, e.cma.With(platform.KeyCapabilityID, te.ID), "failed to start execution: global rate limit exceeded", e.logger)
				e.metrics.with(platform.KeyTriggerID, te.ID).incrementWorkflowExecutionRateLimitGlobalCounter(ctx)
				continue
			}

			cma := e.cma.With(platform.KeyWorkflowExecutionID, executionID)
			err = e.startExecution(ctx, executionID, te.ID, resp.Event.Outputs)
			if err != nil {
				e.logger.With(platform.KeyWorkflowExecutionID, executionID).Errorf("failed to start execution: %v", err)
				logCustMsg(ctx, cma, fmt.Sprintf("failed to start execution: %s", err), e.logger)
				e.metrics.with(platform.KeyTriggerID, te.ID).incrementTriggerWorkflowStarterErrorCounter(ctx)
			} else {
				e.logger.With(platform.KeyWorkflowExecutionID, executionID).Debug("execution started")
				logCustMsg(ctx, cma, "execution started", e.logger)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (e *Engine) workerForStepRequest(ctx context.Context, msg stepRequest) {
	// Instantiate a child logger; in addition to the WorkflowID field the workflow
	// logger will already have, this adds the `stepRef` and `executionID`
	l := e.logger.With(platform.KeyStepRef, msg.stepRef, platform.KeyWorkflowExecutionID, msg.state.ExecutionID)
	cma := e.cma.With(platform.KeyStepRef, msg.stepRef, platform.KeyWorkflowExecutionID, msg.state.ExecutionID)

	l.Debug("executing on a step event")
	stepState := &store.WorkflowExecutionStep{
		Outputs:     store.StepOutput{},
		ExecutionID: msg.state.ExecutionID,
		Ref:         msg.stepRef,
	}

	meteringReport, mrOK := e.meterReports.Get(msg.state.ExecutionID)
	if mrOK {
		curStep, err := e.workflow.Vertex(msg.stepRef)
		if err != nil {
			l.Error(fmt.Sprintf("failed to get current step %s for metering report: %s", stepState.Ref, err))
		}
		capInfo, err := curStep.capability.Info(ctx)
		if err != nil {
			l.Error(fmt.Sprintf("failed to get capability info for %s: %s", stepState.Ref, err))
		}
		err = meteringReport.ReserveStep(metering.ReportStepRef(stepState.Ref), capInfo)
		if err != nil {
			l.Error(fmt.Sprintf("failed to reserve on metering report for %s: %s", stepState.Ref, err))
		}
	} else {
		e.metrics.with(platform.KeyWorkflowID, e.workflow.id, platform.KeyWorkflowExecutionID, msg.state.ExecutionID).incrementWorkflowMissingMeteringReport(ctx)
		// TODO: to be bumped to error if all capabilities must implement metering
		l.Warnf("no metering report found for %v", msg.state.ExecutionID)
	}

	logCustMsg(ctx, cma, "executing step", l)

	stepExecutionStartTime := time.Now()
	inputs, response, sErr := e.executeStep(ctx, l, msg)
	stepExecutionDuration := time.Since(stepExecutionStartTime).Seconds()

	curStepID := "UNSET"
	curStep, verr := e.workflow.Vertex(msg.stepRef)
	if verr == nil {
		curStepID = curStep.ID
	} else {
		l.Errorf("failed to resolve step in workflow; error %v", verr)
	}
	e.metrics.with(platform.KeyCapabilityID, curStepID).updateWorkflowStepDurationHistogram(ctx, int64(stepExecutionDuration))

	var stepStatus string
	switch {
	case sErr != nil && capabilities.ErrStopExecution.Is(sErr):
		lmsg := "step executed successfully with a termination"
		l.Info(lmsg)
		logCustMsg(ctx, cma, lmsg, l)
		stepStatus = store.StatusCompletedEarlyExit
	case sErr != nil:
		lmsg := fmt.Sprintf("error executing step request: %s", sErr)
		l.Error(lmsg)
		logCustMsg(ctx, cma, lmsg, l)
		stepStatus = store.StatusErrored
	default:
		lmsg := "step executed successfully"
		l.With("outputs", response.Value).Info(lmsg)
		// TODO ks-462 emit custom message with outputs
		logCustMsg(ctx, cma, lmsg, l)
		stepStatus = store.StatusCompleted
	}

	if mrOK {
		err := meteringReport.SetStep(metering.ReportStepRef(stepState.Ref), response.Metadata.Metering)
		if err != nil {
			l.Error(fmt.Sprintf("failed to set metering report step for ref %s: %s", stepState.Ref, err))
		}
	}

	stepState.Status = stepStatus
	stepState.Outputs.Value = response.Value
	stepState.Outputs.Err = sErr
	stepState.Inputs = inputs

	// Let's try and emit the stepUpdate.
	// If the context is canceled, we'll just drop the update.
	// This means the engine is shutting down and the
	// receiving loop may not pick up any messages we emit.
	// Note: When full persistence support is added, any hanging steps
	// like this one will get picked up again and will be reprocessed.
	l.Debugf("trying to send step state update for execution %s with status %s", stepState.ExecutionID, stepStatus)
	if err := e.stepUpdatesChMap.send(ctx, stepState.ExecutionID, *stepState); err != nil {
		l.Errorf("failed to issue step state update; error %v", err)
		return
	}
	l.Debugf("sent step state update for execution %s with status %s", stepState.ExecutionID, stepStatus)
}

func merge(baseConfig *values.Map, capConfig capabilities.CapabilityConfiguration) *values.Map {
	restrictedKeys := map[string]bool{}
	for _, k := range capConfig.RestrictedKeys {
		restrictedKeys[k] = true
	}

	// Shallow copy the defaults set in the onchain capability config.
	m := values.EmptyMap()

	if capConfig.DefaultConfig != nil {
		for k, v := range capConfig.DefaultConfig.Underlying {
			m.Underlying[k] = v
		}
	}

	// Add in user-provided config, but skipping any restricted keys
	for k, v := range baseConfig.Underlying {
		if !restrictedKeys[k] {
			m.Underlying[k] = v
		}
	}

	if capConfig.RestrictedConfig == nil {
		return m
	}

	// Then overwrite the config with any restricted settings.
	for k, v := range capConfig.RestrictedConfig.Underlying {
		m.Underlying[k] = v
	}

	return m
}

func (e *Engine) interpolateEnvVars(config map[string]any, env exec.Env) (*values.Map, error) {
	conf, err := exec.FindAndInterpolateEnvVars(config, env)
	if err != nil {
		return nil, err
	}

	confm, ok := conf.(map[string]any)
	if !ok {
		return nil, err
	}

	configMap, err := values.NewMap(confm)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

// configForStep fetches the config for the step from the workflow registry (for secrets) and from the capabilities
// registry (for capability-level configuration). It doesn't perform any caching of the config values, since
// the two registries perform their own caching.
func (e *Engine) configForStep(ctx context.Context, lggr logger.Logger, step *step) (*values.Map, error) {
	secrets, err := e.secretsFetcher(ctx, e.workflow.owner, e.workflow.name.Hex(), e.workflow.name.String(), e.workflow.id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secrets: %w", err)
	}

	env := exec.Env{
		Config:  e.env.Config,
		Binary:  e.env.Binary,
		Secrets: secrets,
	}
	config, err := e.interpolateEnvVars(step.Config, env)
	if err != nil {
		return nil, fmt.Errorf("failed to interpolate env vars: %w", err)
	}

	ID := step.info.ID

	// If the capability info is missing a DON, then
	// the capability is local, and we should use the localNode's DON ID.
	var donID uint32
	if !step.info.IsLocal {
		donID = step.info.DON.ID
	} else {
		donID = e.localNode.Load().WorkflowDON.ID
	}

	capConfig, err := e.registry.ConfigForCapability(ctx, ID, donID)
	if err != nil {
		lggr.Warnw(fmt.Sprintf("could not retrieve config from remote registry: %s", err), "capabilityID", ID)
		return config, nil
	}

	// Merge the capability registry config with the config provided by the user.
	// We need to obey the following rules:
	// - Remove any restricted keys
	// - Overlay any restricted config
	// - Merge the other keys, with user keys taking precedence
	return merge(config, capConfig), nil
}

// executeStep executes the referenced capability within a step and returns the result.
func (e *Engine) executeStep(ctx context.Context, lggr logger.Logger, msg stepRequest) (*values.Map, capabilities.CapabilityResponse, error) {
	curStep, err := e.workflow.Vertex(msg.stepRef)
	if err != nil {
		return nil, capabilities.CapabilityResponse{}, err
	}

	var inputs any
	if curStep.Inputs.OutputRef != "" {
		inputs = curStep.Inputs.OutputRef
	} else {
		inputs = curStep.Inputs.Mapping
	}

	i, err := exec.FindAndInterpolateAllKeys(inputs, msg.state)
	if err != nil {
		return nil, capabilities.CapabilityResponse{}, err
	}

	inputsMap, err := values.NewMap(i.(map[string]any))
	if err != nil {
		return nil, capabilities.CapabilityResponse{}, err
	}

	config, err := e.configForStep(ctx, lggr, curStep)
	if err != nil {
		return nil, capabilities.CapabilityResponse{}, err
	}
	stepTimeoutDuration := e.stepTimeoutDuration
	if timeoutOverride, ok := config.Underlying[reservedFieldNameStepTimeout]; ok {
		var desiredTimeout int64
		err2 := timeoutOverride.UnwrapTo(&desiredTimeout)
		if err2 != nil {
			e.logger.Warnw("couldn't decode step timeout override, using default", "error", err2, "default", stepTimeoutDuration)
		} else if desiredTimeout == 0 {
			e.logger.Debugw("overridden timeout set to 0, using default", "default", stepTimeoutDuration)
		} else {
			if desiredTimeout > maxStepTimeoutOverrideSec {
				e.logger.Warnw("desired step timeout is too large, limiting to max value", "maxValue", maxStepTimeoutOverrideSec)
				desiredTimeout = maxStepTimeoutOverrideSec
			}
			stepTimeoutDuration = time.Duration(desiredTimeout) * time.Second
		}
	}

	ln := e.localNode.Load()
	tr := capabilities.CapabilityRequest{
		Inputs: inputsMap,
		Config: config,
		Metadata: capabilities.RequestMetadata{
			WorkflowID:               msg.state.WorkflowID,
			WorkflowExecutionID:      msg.state.ExecutionID,
			WorkflowOwner:            e.workflow.owner,
			WorkflowName:             e.workflow.name.Hex(),
			WorkflowDonID:            ln.WorkflowDON.ID,
			WorkflowDonConfigVersion: ln.WorkflowDON.ConfigVersion,
			ReferenceID:              msg.stepRef,
			DecodedWorkflowName:      e.workflow.name.String(),
		},
	}

	stepCtx, cancel := context.WithTimeout(ctx, stepTimeoutDuration)
	defer cancel()

	e.metrics.with(platform.KeyCapabilityID, curStep.ID).incrementCapabilityInvocationCounter(ctx)
	err = events.EmitCapabilityStartedEvent(ctx, e.cma, msg.state.ExecutionID, curStep.ID, msg.stepRef)
	if err != nil {
		e.logger.Errorf("failed to emit capability event: %v", err)
	}
	output, capErr := curStep.capability.Execute(stepCtx, tr)
	status := store.StatusCompleted

	if capErr != nil {
		status = store.StatusErrored
		if capabilities.ErrStopExecution.Is(capErr) {
			status = store.StatusCompletedEarlyExit
		}
	}

	defer func() {
		if err := events.EmitCapabilityFinishedEvent(ctx, e.cma, msg.state.ExecutionID, curStep.ID, msg.stepRef, status); err != nil {
			e.logger.Errorf("failed to emit capability event: %v", err)
		}
	}()

	if capErr != nil {
		e.metrics.with(platform.KeyStepRef, msg.stepRef, platform.KeyCapabilityID, curStep.ID).incrementCapabilityFailureCounter(ctx)
		return inputsMap, capabilities.CapabilityResponse{}, capErr
	}

	return inputsMap, output, nil
}

func (e *Engine) deregisterTrigger(ctx context.Context, t *triggerCapability, triggerIdx int) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.trigger != nil && t.registered {
		// if t.trigger == nil or !t.registered, then we haven't initialized the workflow
		// yet, and can safely consider the trigger deregistered with
		// no further action.
		ln := e.localNode.Load()
		deregRequest := capabilities.TriggerRegistrationRequest{
			Metadata: capabilities.RequestMetadata{
				WorkflowID:               e.workflow.id,
				WorkflowDonID:            ln.WorkflowDON.ID,
				WorkflowDonConfigVersion: ln.WorkflowDON.ConfigVersion,
				WorkflowName:             e.workflow.name.Hex(),
				WorkflowOwner:            e.workflow.owner,
				ReferenceID:              t.Ref,
				DecodedWorkflowName:      e.workflow.name.String(),
			},
			TriggerID: generateTriggerID(e.workflow.id, triggerIdx),
			Config:    t.config,
		}
		return t.trigger.UnregisterTrigger(ctx, deregRequest)
	}

	return nil
}

func (e *Engine) isWorkflowFullyProcessed(_ context.Context, state store.WorkflowExecution) (bool, string, error) {
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

// heartbeat runs by default every defaultHeartbeatCadence minutes
func (e *Engine) heartbeat(ctx context.Context) {
	defer e.wg.Done()

	ticker := time.NewTicker(e.heartbeatCadence)
	defer ticker.Stop()

	// Gauge values are "persisted" by the backend
	// and will be continually reported until the value changes.
	e.metrics.engineHeartbeatGauge(ctx, 1)

	for {
		select {
		case <-ctx.Done():
			e.metrics.engineHeartbeatGauge(ctx, 0)
			e.logger.Info("shutting down heartbeat")
			return
		case <-ticker.C:
			e.metrics.incrementEngineHeartbeatCounter(ctx)
			e.metrics.updateTotalWorkflowsGauge(ctx, e.stepUpdatesChMap.len())
			logCustMsg(ctx, e.cma, "engine heartbeat at: "+e.clock.Now().Format(time.RFC3339), e.logger)
		}
	}
}

func (e *Engine) sendReportToBilling(ctx context.Context, report *metering.Report, workflowID string, executionID string) error {
	if e.billingClient != nil {
		req := billing.SubmitWorkflowReceiptRequest{
			AccountId:           e.workflow.owner,
			WorkflowId:          workflowID,
			WorkflowExecutionId: executionID,
			Metering:            report.Message(),
		}

		resp, err := e.billingClient.SubmitWorkflowReceipt(ctx, &req)
		if err != nil {
			return err
		}

		if resp == nil || !resp.Success {
			return errors.New("failed to submit workflow receipt")
		}
	}

	return nil
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
				e.logger.Errorf("failed to deregister trigger: %v", err)
			}
		}

		close(e.stopCh)
		e.logger.Info("stopCh closed, waiting for workers to finish")
		e.wg.Wait()
		e.logger.Info("workers finished")

		err := e.workflow.walkDo(workflows.KeywordTrigger, func(s *step) error {
			if s.Ref == workflows.KeywordTrigger {
				return nil
			}

			// if capability is nil, then we haven't initialized
			// the workflow yet and can safely consider it deregistered
			// with no further action.
			if s.capability == nil {
				return nil
			}

			stepConfig, err := e.configForStep(ctx, e.logger, s)
			if err != nil {
				return fmt.Errorf("cannot fetch config for step: %w", err)
			}

			reg := capabilities.UnregisterFromWorkflowRequest{
				Metadata: capabilities.RegistrationMetadata{
					WorkflowID:    e.workflow.id,
					WorkflowOwner: e.workflow.owner,
					ReferenceID:   s.Vertex.Ref,
				},
				Config: stepConfig,
			}

			unregisterError := s.capability.UnregisterFromWorkflow(ctx, reg)
			if unregisterError != nil {
				// As a capability may be closing/closed down during shutdown it possible that unregistering the
				// workflow will fail but this is not necessarily indicative of an error so log it as a warning
				e.logger.Warn("failed to unregister capability from workflow",
					"workflowID", e.workflow.id,
					"stepID", s.ID,
					"stepRef", s.Ref,
					"error", unregisterError)
			}

			return nil
		})
		if err != nil {
			return err
		}
		// decrement the global and per owner engine counter
		e.workflowLimits.Decrement(e.workflow.owner)

		logCustMsg(ctx, e.cma, "workflow unregistered", e.logger)
		e.metrics.incrementWorkflowUnregisteredCounter(ctx)
		return nil
	})
}

func (e *Engine) HealthReport() map[string]error {
	return map[string]error{e.Name(): e.Healthy()}
}

func (e *Engine) Name() string {
	return e.logger.Name()
}

type Config struct {
	Workflow             sdk.WorkflowSpec
	WorkflowID           string
	WorkflowOwner        string
	WorkflowName         types.WorkflowName
	Lggr                 logger.Logger
	Registry             core.CapabilitiesRegistry
	MaxWorkerLimit       int
	QueueSize            int
	NewWorkerTimeout     time.Duration
	MaxExecutionDuration time.Duration
	Store                store.Store
	Config               []byte
	Binary               []byte
	SecretsFetcher       SecretsFor
	HeartbeatCadence     time.Duration
	StepTimeout          time.Duration
	BillingClient        BillingClient

	// RateLimiter limits the workflow execution steps globally and per
	// second that a workflow owner can make
	RateLimiter *ratelimiter.RateLimiter

	// WorkflowLimits specifies an upper limit on the count of workflows that can be
	// running globally and per workflow owner.
	WorkflowLimits *syncerlimiter.Limits

	// For testing purposes only
	maxRetries          int
	retryMs             int
	afterInit           func(success bool)
	onExecutionFinished func(weid string)
	onRateLimit         func(weid string)
	clock               clockwork.Clock
}

const (
	defaultWorkerLimit          = 100
	defaultQueueSize            = 1000
	defaultNewWorkerTimeout     = 2 * time.Second
	defaultMaxExecutionDuration = 10 * time.Minute
	defaultHeartbeatCadence     = 5 * time.Minute
	defaultStepTimeout          = 2 * time.Minute
)

func NewEngine(ctx context.Context, cfg Config) (engine *Engine, err error) {
	if cfg.Store == nil {
		return nil, &workflowError{reason: "store is nil",
			labels: map[string]string{
				platform.KeyWorkflowID: cfg.WorkflowID,
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

	if cfg.HeartbeatCadence == 0 {
		cfg.HeartbeatCadence = defaultHeartbeatCadence
	}

	if cfg.StepTimeout == 0 {
		cfg.StepTimeout = defaultStepTimeout
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

	if cfg.onRateLimit == nil {
		cfg.onRateLimit = func(weid string) {}
	}

	if cfg.clock == nil {
		cfg.clock = clockwork.NewRealClock()
	}

	if cfg.RateLimiter == nil {
		return nil, &workflowError{reason: "ratelimiter must be provided",
			labels: map[string]string{
				platform.KeyWorkflowID: cfg.WorkflowID,
			},
		}
	}

	if cfg.WorkflowLimits == nil {
		return nil, &workflowError{reason: "workflowLimits must be provided",
			labels: map[string]string{
				platform.KeyWorkflowID: cfg.WorkflowID,
			},
		}
	}

	if cfg.Registry == nil {
		return nil, &workflowError{reason: "capabilities registry must be provided",
			labels: map[string]string{
				platform.KeyWorkflowID: cfg.WorkflowID,
			},
		}
	}

	// TODO: validation of the workflow spec
	// We'll need to check, among other things:
	// - that there are no step `ref` called `trigger` as this is reserved for any triggers
	// - that there are no duplicate `ref`s
	// - that the `ref` for any triggers is empty -- and filled in with `trigger`
	// - that the resulting graph is strongly connected (i.e. no disjointed subgraphs exist)
	// - etc.

	// spin up monitoring resources
	em, err := initMonitoringResources()
	if err != nil {
		return nil, fmt.Errorf("could not initialize monitoring resources: %w", err)
	}

	nodeState, err := cfg.Registry.LocalNode(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get local node state: %w", err)
	}
	cma := custmsg.NewLabeler().With(platform.KeyWorkflowID, cfg.WorkflowID,
		platform.KeyWorkflowOwner, cfg.WorkflowOwner,
		platform.KeyWorkflowName, cfg.WorkflowName.String(),
		platform.KeyWorkflowVersion, platform.ValueWorkflowVersion,
		platform.KeyDonID, strconv.Itoa(int(nodeState.WorkflowDON.ID)),
		platform.KeyDonF, strconv.Itoa(int(nodeState.WorkflowDON.F)),
		platform.KeyDonN, strconv.Itoa(len(nodeState.WorkflowDON.Members)),
		platform.KeyDonQ, strconv.Itoa(aggregation.ByzantineQuorum(
			len(nodeState.WorkflowDON.Members),
			int(nodeState.WorkflowDON.F),
		)),
		platform.KeyP2PID, nodeState.PeerID.String(),
	)
	workflow, err := Parse(cfg.Workflow)
	if err != nil {
		logCustMsg(ctx, cma, fmt.Sprintf("failed to parse workflow: %s", err), cfg.Lggr)
		return nil, err
	}

	workflow.id = cfg.WorkflowID
	workflow.owner = cfg.WorkflowOwner
	workflow.name = cfg.WorkflowName

	lggr := cfg.Lggr.With("workflowID", cfg.WorkflowID)

	engine = &Engine{
		cma:            cma,
		logger:         lggr.Named("WorkflowEngine"),
		metrics:        workflowsMetricLabeler{metrics.NewLabeler().With(platform.KeyWorkflowID, cfg.WorkflowID, platform.KeyWorkflowOwner, cfg.WorkflowOwner, platform.KeyWorkflowName, cfg.WorkflowName.String()), *em},
		registry:       cfg.Registry,
		workflow:       workflow,
		secretsFetcher: cfg.SecretsFetcher,
		env: exec.Env{
			Config: cfg.Config,
			Binary: cfg.Binary,
		},
		executionsStore:      cfg.Store,
		pendingStepRequests:  make(chan stepRequest, cfg.QueueSize),
		stepUpdatesChMap:     stepUpdateManager{m: map[string]stepUpdateChannel{}},
		triggerEvents:        make(chan capabilities.TriggerResponse),
		stopCh:               make(chan struct{}),
		newWorkerTimeout:     cfg.NewWorkerTimeout,
		stepTimeoutDuration:  cfg.StepTimeout,
		maxExecutionDuration: cfg.MaxExecutionDuration,
		heartbeatCadence:     cfg.HeartbeatCadence,
		onExecutionFinished:  cfg.onExecutionFinished,
		onRateLimit:          cfg.onRateLimit,
		afterInit:            cfg.afterInit,
		maxRetries:           cfg.maxRetries,
		retryMs:              cfg.retryMs,
		maxWorkerLimit:       cfg.MaxWorkerLimit,
		clock:                cfg.clock,
		ratelimiter:          cfg.RateLimiter,
		workflowLimits:       cfg.WorkflowLimits,
		meterReports:         metering.NewReports(),
		billingClient:        cfg.BillingClient,
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
	for label := range platform.LabelKeysSorted() {
		// This will silently ignore any labels that are not present in the map
		// are we ok with this?
		if value, ok := e.labels[label]; ok {
			errStr = fmt.Sprintf("%s %s: %s", label, value, errStr)
		}
	}

	return errStr
}

func logCustMsg(ctx context.Context, cma custmsg.MessageEmitter, msg string, log logger.Logger) {
	err := cma.Emit(ctx, msg)
	if err != nil {
		log.Errorf("failed to send custom message with msg: %s, err: %v", msg, err)
	}
}
